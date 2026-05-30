"use client";

import {
  Button,
  Input,
  Select,
  Slider,
  Switch,
  TextArea,
  Toast,
  Typography,
} from "@douyinfe/semi-ui";
import {
  IconCopy,
  IconDelete,
  IconDownload,
  IconEdit,
  IconRefresh,
  IconUpload,
} from "@douyinfe/semi-icons";
import { useSearchParams } from "next/navigation";
import { useCallback, useEffect, useMemo, useRef, useState } from "react";

import type { CatalogListResponse, TokenGroupItem } from "@/lib/api/catalog-types";
import { ApiError } from "@/lib/api/client";
import type {
  PlaygroundConfig,
  PlaygroundMessage,
} from "@/lib/api/playground-types";
import { streamPlaygroundChat } from "@/lib/api/playground-stream";
import { useAuth } from "@/contexts/auth-context";
import { useDeferredEffect } from "@/lib/hooks/use-deferred-effect";
import { useGatewayRequest } from "@/lib/hooks/use-gateway-request";

const { Title, Text } = Typography;

const CONFIG_STORAGE_KEY = "portal_playground_config_v1";

function newId(): string {
  return `${Date.now()}-${Math.random().toString(36).slice(2, 9)}`;
}

const defaultConfig = (): PlaygroundConfig => ({
  version: 1,
  token_group: "default",
  model: "",
  custom_body: false,
  custom_body_raw: "",
  image_url: "",
  max_tokens: 2048,
  temperature: 1,
  temperature_enabled: true,
  top_p: 1,
  top_p_enabled: false,
  frequency_penalty: 0,
  frequency_penalty_enabled: false,
  presence_penalty: 0,
  presence_penalty_enabled: false,
});

export function PlaygroundPage() {
  const { token } = useAuth();
  const request = useGatewayRequest();
  const searchParams = useSearchParams();
  const abortRef = useRef<AbortController | null>(null);
  const bottomRef = useRef<HTMLDivElement>(null);

  const [config, setConfig] = useState<PlaygroundConfig>(defaultConfig);
  const [messages, setMessages] = useState<PlaygroundMessage[]>([]);
  const [input, setInput] = useState("");
  const [streaming, setStreaming] = useState(false);
  const [models, setModels] = useState<string[]>([]);
  const [groups, setGroups] = useState<TokenGroupItem[]>([]);
  const [apiKeyId, setApiKeyId] = useState<number | undefined>();
  const [tokenHint, setTokenHint] = useState<string | null>(null);
  const [debugOpen, setDebugOpen] = useState(false);
  const [debugPayload, setDebugPayload] = useState<unknown>(null);
  const [editingId, setEditingId] = useState<string | null>(null);
  const [editDraft, setEditDraft] = useState("");

  useDeferredEffect(() => {
    const raw = sessionStorage.getItem(CONFIG_STORAGE_KEY);
    if (raw) {
      try {
        const parsed = JSON.parse(raw) as PlaygroundConfig;
        setConfig((c) => ({ ...c, ...parsed, version: 1 }));
      } catch {
        /* ignore */
      }
    }
    const fromQuery = searchParams.get("token_id");
    if (fromQuery) {
      sessionStorage.setItem("playground_token_id", fromQuery);
    }
    const idStr = sessionStorage.getItem("playground_token_id");
    const name = sessionStorage.getItem("playground_token_name");
    if (idStr) {
      const id = Number(idStr);
      if (!Number.isNaN(id)) setApiKeyId(id);
      setTokenHint(name ? `${name} (#${idStr})` : `#${idStr}`);
    }
  }, [searchParams]);

  useEffect(() => {
    sessionStorage.setItem(CONFIG_STORAGE_KEY, JSON.stringify(config));
  }, [config]);

  const loadCatalog = useCallback(async () => {
    try {
      const group = config.token_group || "default";
      const [catalogRes, groupRes] = await Promise.all([
        request<CatalogListResponse>(
          `/api/models?page=1&page_size=200&token_group=${encodeURIComponent(group)}`,
        ),
        request<{ success: boolean; data: TokenGroupItem[] }>("/api/token-groups"),
      ]);
      const names = catalogRes.data.items.map((m) => m.model);
      setModels(names);
      setGroups(groupRes.data);
      setConfig((c) => ({
        ...c,
        model: c.model && names.includes(c.model) ? c.model : (names[0] ?? ""),
      }));
    } catch (e) {
      Toast.error(e instanceof ApiError ? e.message : "加载模型失败");
    }
  }, [request, config.token_group]);

  useDeferredEffect(() => loadCatalog(), [loadCatalog]);

  const patchConfig = (partial: Partial<PlaygroundConfig>) => {
    setConfig((c) => ({ ...c, ...partial }));
  };

  const buildApiMessages = useCallback(
    (list: PlaygroundMessage[]) => {
      const out: { role: PlaygroundMessage["role"]; content: string }[] = [];
      for (const m of list) {
        if (m.role === "assistant" && !m.content.trim()) continue;
        if (m.role === "user" && config.image_url?.trim()) {
          out.push({
            role: "user",
            content: m.content,
          });
          continue;
        }
        out.push({ role: m.role, content: m.content });
      }
      return out;
    },
    [config.image_url],
  );

  const runCompletion = useCallback(
    async (baseMessages: PlaygroundMessage[], regenerate = false) => {
      if (!token) {
        Toast.error("请先登录");
        return;
      }
      if (!config.model && !config.custom_body) {
        Toast.error("请选择模型");
        return;
      }
      if (streaming) return;

      let contextMessages = [...baseMessages];
      if (regenerate) {
        while (
          contextMessages.length > 0 &&
          contextMessages[contextMessages.length - 1].role === "assistant"
        ) {
          contextMessages = contextMessages.slice(0, -1);
        }
      }

      const last = contextMessages[contextMessages.length - 1];
      if (!last || last.role !== "user" || !last.content.trim()) {
        Toast.warning("请输入消息");
        return;
      }

      const assistantId = newId();
      setMessages([
        ...contextMessages,
        { id: assistantId, role: "assistant", content: "" },
      ]);
      setStreaming(true);
      abortRef.current?.abort();
      abortRef.current = new AbortController();

      let customRaw: unknown;
      if (config.custom_body && config.custom_body_raw.trim()) {
        try {
          customRaw = JSON.parse(config.custom_body_raw) as unknown;
        } catch {
          setStreaming(false);
          Toast.error("自定义 Body JSON 无效");
          return;
        }
      }

      const reqBody = {
        model: config.model,
        messages: buildApiMessages(contextMessages),
        max_tokens: config.max_tokens,
        api_key_id: apiKeyId,
        token_group: config.token_group,
        custom_body: config.custom_body,
        ...(customRaw !== undefined ? { custom_body_raw: customRaw } : {}),
        ...(config.temperature_enabled ? { temperature: config.temperature } : {}),
        ...(config.top_p_enabled ? { top_p: config.top_p } : {}),
        ...(config.frequency_penalty_enabled
          ? { frequency_penalty: config.frequency_penalty }
          : {}),
        ...(config.presence_penalty_enabled
          ? { presence_penalty: config.presence_penalty }
          : {}),
      };

      setDebugPayload(reqBody);

      await streamPlaygroundChat(
        token,
        reqBody,
        {
          onDelta: (text) => {
            setMessages((prev) =>
              prev.map((m) =>
                m.id === assistantId ? { ...m, content: m.content + text } : m,
              ),
            );
          },
          onDone: () => {
            setStreaming(false);
            setMessages((prev) =>
              prev.map((m) =>
                m.id === assistantId
                  ? {
                      ...m,
                      debug: {
                        request: reqBody,
                        finished_at: new Date().toISOString(),
                      },
                    }
                  : m,
              ),
            );
          },
          onError: (message) => {
            setStreaming(false);
            Toast.error(message);
          },
        },
        abortRef.current.signal,
      );
    },
    [token, config, apiKeyId, streaming, buildApiMessages],
  );

  const sendMessage = async () => {
    const text = input.trim();
    if (!text) return;
    const userMsg: PlaygroundMessage = { id: newId(), role: "user", content: text };
    setInput("");
    const next = [...messages, userMsg];
    setMessages(next);
    await runCompletion(next);
  };

  const regenerate = async () => {
    if (!messages.some((m) => m.role === "user")) {
      Toast.warning("没有可重新生成的对话");
      return;
    }
    await runCompletion(messages, true);
  };

  const copyMessage = async (content: string) => {
    try {
      await navigator.clipboard.writeText(content);
      Toast.success("已复制");
    } catch {
      Toast.error("复制失败");
    }
  };

  const deleteMessage = (id: string) => {
    setMessages((prev) => prev.filter((m) => m.id !== id));
  };

  const startEdit = (m: PlaygroundMessage) => {
    setEditingId(m.id);
    setEditDraft(m.content);
  };

  const saveEdit = () => {
    if (!editingId) return;
    setMessages((prev) =>
      prev.map((m) => (m.id === editingId ? { ...m, content: editDraft } : m)),
    );
    setEditingId(null);
    setEditDraft("");
  };

  const exportConfig = () => {
    const blob = new Blob([JSON.stringify(config, null, 2)], {
      type: "application/json",
    });
    const url = URL.createObjectURL(blob);
    const a = document.createElement("a");
    a.href = url;
    a.download = "playground-config.json";
    a.click();
    URL.revokeObjectURL(url);
    Toast.success("配置已导出");
  };

  const importConfig = () => {
    const inputEl = document.createElement("input");
    inputEl.type = "file";
    inputEl.accept = "application/json";
    inputEl.onchange = async () => {
      const file = inputEl.files?.[0];
      if (!file) return;
      try {
        const text = await file.text();
        const parsed = JSON.parse(text) as PlaygroundConfig;
        setConfig({ ...defaultConfig(), ...parsed, version: 1 });
        Toast.success("配置已导入");
      } catch {
        Toast.error("无效的配置文件");
      }
    };
    inputEl.click();
  };

  const groupOptions = useMemo(
    () => groups.map((g) => ({ value: g.slug, label: `${g.name} (${g.multiplier}x)` })),
    [groups],
  );

  const modelOptions = useMemo(
    () => models.map((m) => ({ value: m, label: m })),
    [models],
  );

  useEffect(() => {
    bottomRef.current?.scrollIntoView({ behavior: "smooth" });
  }, [messages, streaming]);

  return (
    <div className="flex h-[calc(100vh-120px)] min-h-[560px] flex-col gap-3">
      {tokenHint && (
        <Text type="secondary" size="small">
          当前令牌：{tokenHint}
        </Text>
      )}

      <div className="flex min-h-0 flex-1 gap-3">
        <aside className="flex w-[300px] shrink-0 flex-col overflow-hidden rounded-xl border border-[var(--semi-color-border)] bg-white">
          <div className="border-b border-[var(--semi-color-border)] px-4 py-3">
            <Text strong>模型配置</Text>
          </div>
          <div className="flex-1 overflow-y-auto p-4">
            <div className="mb-4 flex items-center justify-between gap-2">
              <Text size="small" type="secondary">
                自定义请求体模式
              </Text>
              <Switch
                checked={config.custom_body}
                onChange={(v) => patchConfig({ custom_body: v })}
                size="small"
              />
            </div>

            {config.custom_body ? (
            <TextArea
              value={config.custom_body_raw ?? ""}
              onChange={(v) => patchConfig({ custom_body_raw: v })}
              rows={16}
              placeholder='{"model":"gpt-4o","messages":[...]}'
            />
          ) : (
            <div className="flex flex-col gap-4">
              <div>
                <Text type="secondary" size="small" className="mb-1 block">
                  分组
                </Text>
                <Select
                  value={config.token_group}
                  optionList={groupOptions}
                  onChange={(v) => patchConfig({ token_group: String(v) })}
                  className="w-full"
                />
              </div>
              <div>
                <Text type="secondary" size="small" className="mb-1 block">
                  模型
                </Text>
                <Select
                  filter
                  value={config.model}
                  optionList={modelOptions}
                  onChange={(v) => patchConfig({ model: String(v) })}
                  className="w-full"
                />
              </div>
              <div>
                <Text type="secondary" size="small" className="mb-1 block">
                  图片 URL（多模态占位）
                </Text>
                <Input
                  value={config.image_url ?? ""}
                  onChange={(v) => patchConfig({ image_url: v })}
                  placeholder="https://..."
                />
              </div>
              <div>
                <Text type="secondary" size="small" className="mb-1 block">
                  Max tokens: {config.max_tokens}
                </Text>
                <Slider
                  min={256}
                  max={8192}
                  step={256}
                  value={config.max_tokens}
                  onChange={(v) => patchConfig({ max_tokens: Number(v) })}
                />
              </div>
              <ParamSlider
                label="Temperature"
                enabled={config.temperature_enabled}
                onEnabledChange={(v) => patchConfig({ temperature_enabled: v })}
                value={config.temperature}
                min={0}
                max={2}
                step={0.1}
                onChange={(v) => patchConfig({ temperature: v })}
              />
              <ParamSlider
                label="Top P"
                enabled={config.top_p_enabled}
                onEnabledChange={(v) => patchConfig({ top_p_enabled: v })}
                value={config.top_p}
                min={0}
                max={1}
                step={0.05}
                onChange={(v) => patchConfig({ top_p: v })}
              />
              <ParamSlider
                label="Frequency penalty"
                enabled={config.frequency_penalty_enabled}
                onEnabledChange={(v) => patchConfig({ frequency_penalty_enabled: v })}
                value={config.frequency_penalty}
                min={-2}
                max={2}
                step={0.1}
                onChange={(v) => patchConfig({ frequency_penalty: v })}
              />
              <ParamSlider
                label="Presence penalty"
                enabled={config.presence_penalty_enabled}
                onEnabledChange={(v) => patchConfig({ presence_penalty_enabled: v })}
                value={config.presence_penalty}
                min={-2}
                max={2}
                step={0.1}
                onChange={(v) => patchConfig({ presence_penalty: v })}
              />
            </div>
            )}
          </div>
          <div className="flex gap-2 border-t border-[var(--semi-color-border)] p-3">
            <Button icon={<IconDownload />} type="primary" theme="solid" block onClick={exportConfig}>
              导出
            </Button>
            <Button icon={<IconUpload />} block onClick={importConfig}>
              导入
            </Button>
          </div>
        </aside>

        <section className="flex min-w-0 flex-1 flex-col rounded-xl border border-[var(--semi-color-border)] bg-white">
          <div className="flex items-center justify-between border-b border-[var(--semi-color-border)] px-4 py-3">
            <div>
              <Text strong>AI 对话</Text>
              {config.model && (
                <Text type="tertiary" size="small" className="block">
                  {config.model}
                </Text>
              )}
            </div>
            <Button
              size="small"
              theme="borderless"
              onClick={() => {
                setDebugOpen((o) => !o);
                if (!debugOpen) setDebugPayload(null);
              }}
            >
              {debugOpen ? "隐藏调试" : "显示调试"}
            </Button>
          </div>
          <div className="flex-1 overflow-y-auto p-4">
            {messages.length === 0 && (
              <Text type="tertiary" className="block text-center py-12">
                发送消息开始对话
              </Text>
            )}
            {messages.map((m) => (
              <div
                key={m.id}
                className={`mb-4 flex ${m.role === "user" ? "justify-end" : "justify-start"}`}
              >
                <div
                  className={`max-w-[85%] rounded-xl px-4 py-3 ${
                    m.role === "user"
                      ? "bg-[#007AFF] text-white"
                      : "bg-[#f6f7f9] text-[var(--semi-color-text-0)]"
                  }`}
                >
                  {editingId === m.id ? (
                    <div className="flex flex-col gap-2">
                      <TextArea value={editDraft} onChange={setEditDraft} rows={4} />
                      <div className="flex gap-2">
                        <Button size="small" onClick={saveEdit}>
                          保存
                        </Button>
                        <Button
                          size="small"
                          type="tertiary"
                          onClick={() => setEditingId(null)}
                        >
                          取消
                        </Button>
                      </div>
                    </div>
                  ) : (
                    <pre className="whitespace-pre-wrap font-sans text-sm">{m.content}</pre>
                  )}
                  {m.role === "assistant" && editingId !== m.id && (
                    <div className="mt-2 flex flex-wrap gap-1">
                      <Button
                        size="small"
                        type="tertiary"
                        icon={<IconCopy />}
                        onClick={() => void copyMessage(m.content)}
                      />
                      <Button
                        size="small"
                        type="tertiary"
                        icon={<IconEdit />}
                        onClick={() => startEdit(m)}
                      />
                      <Button
                        size="small"
                        type="tertiary"
                        icon={<IconDelete />}
                        onClick={() => deleteMessage(m.id)}
                      />
                      {m.id === messages[messages.length - 1]?.id && (
                        <Button
                          size="small"
                          type="tertiary"
                          icon={<IconRefresh />}
                          onClick={() => void regenerate()}
                          disabled={streaming}
                        />
                      )}
                    </div>
                  )}
                </div>
              </div>
            ))}
            <div ref={bottomRef} />
          </div>
          <div className="border-t border-[var(--semi-color-border)] p-3">
            <div className="flex gap-2">
              <TextArea
                value={input}
                onChange={setInput}
                placeholder="请输入您的问题…"
                autosize={{ minRows: 2, maxRows: 6 }}
                onEnterPress={(e) => {
                  if (!e.shiftKey) {
                    e.preventDefault();
                    void sendMessage();
                  }
                }}
                disabled={streaming || config.custom_body}
              />
              <Button
                type="primary"
                theme="solid"
                loading={streaming}
                disabled={config.custom_body}
                onClick={() => void sendMessage()}
                className="self-end"
              >
                发送
              </Button>
            </div>
          </div>
        </section>

        {debugOpen && (
          <aside className="w-[280px] shrink-0 overflow-y-auto rounded-xl border border-[var(--semi-color-border)] bg-white p-3">
            <div className="mb-2 flex items-center justify-between">
              <Text strong>调试信息</Text>
              <Button size="small" type="tertiary" onClick={() => setDebugOpen(false)}>
                关闭
              </Button>
            </div>
            <pre className="whitespace-pre-wrap break-all text-xs">
              {JSON.stringify(debugPayload, null, 2)}
            </pre>
          </aside>
        )}
      </div>
    </div>
  );
}

function ParamSlider({
  label,
  enabled,
  onEnabledChange,
  value,
  min,
  max,
  step,
  onChange,
}: {
  label: string;
  enabled: boolean;
  onEnabledChange: (v: boolean) => void;
  value: number;
  min: number;
  max: number;
  step: number;
  onChange: (v: number) => void;
}) {
  return (
    <div>
      <div className="mb-1 flex items-center justify-between">
        <Text type="secondary" size="small">
          {label}
        </Text>
        <Switch size="small" checked={enabled} onChange={onEnabledChange} />
      </div>
      {enabled && (
        <Slider
          min={min}
          max={max}
          step={step}
          value={value}
          onChange={(v) => onChange(Number(v))}
        />
      )}
    </div>
  );
}
