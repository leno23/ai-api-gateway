"use client";

import {
  Button,
  Dropdown,
  Form,
  Input,
  Modal,
  Select,
  Switch,
  Table,
  Tag,
  Toast,
  Typography,
} from "@douyinfe/semi-ui";
import { IconCopy, IconPlus } from "@douyinfe/semi-icons";
import { useRouter } from "next/navigation";
import { useCallback, useMemo, useState } from "react";
import { QRCodeSVG } from "qrcode.react";

import type { ApiTokenRow, TokenCreateResponse, TokenListResponse } from "@/lib/api/token-types";
import { ApiError } from "@/lib/api/client";
import type { TokenGroupItem } from "@/lib/api/catalog-types";
import { useDeferredEffect } from "@/lib/hooks/use-deferred-effect";
import { useGatewayRequest } from "@/lib/hooks/use-gateway-request";

const { Title, Text } = Typography;

type TokenFormValues = {
  name: string;
  token_group: string;
  quota_unlimited: boolean;
  quota_limit?: number;
  models_text?: string;
  ip_whitelist_text?: string;
};

export function TokenManagementPage() {
  const request = useGatewayRequest();
  const router = useRouter();
  const [loading, setLoading] = useState(false);
  const [items, setItems] = useState<ApiTokenRow[]>([]);
  const [total, setTotal] = useState(0);
  const [page, setPage] = useState(1);
  const [selectedKeys, setSelectedKeys] = useState<number[]>([]);
  const [groups, setGroups] = useState<TokenGroupItem[]>([]);

  const [modalOpen, setModalOpen] = useState(false);
  const [editing, setEditing] = useState<ApiTokenRow | null>(null);
  const [secretModal, setSecretModal] = useState<{ key: string; name: string } | null>(null);

  const load = useCallback(async () => {
    setLoading(true);
    try {
      const [listRes, groupRes] = await Promise.all([
        request<TokenListResponse>(`/api/tokens?page=${page}&page_size=10`),
        request<{ success: boolean; data: TokenGroupItem[] }>("/api/token-groups"),
      ]);
      setItems(listRes.data.items);
      setTotal(listRes.data.total);
      setGroups(groupRes.data);
    } catch (e) {
      Toast.error(e instanceof ApiError ? e.message : "加载失败");
    } finally {
      setLoading(false);
    }
  }, [request, page]);

  useDeferredEffect(() => load(), [load]);

  const openCreate = () => {
    setEditing(null);
    setModalOpen(true);
  };

  const openEdit = (row: ApiTokenRow) => {
    setEditing(row);
    setModalOpen(true);
  };

  const copyText = async (text: string, label: string) => {
    await navigator.clipboard.writeText(text);
    Toast.success(`已复制${label}`);
  };

  const toggleEnabled = async (row: ApiTokenRow, enabled: boolean) => {
    try {
      await request(`/api/tokens/${row.id}`, {
        method: "PUT",
        body: JSON.stringify({
          status: enabled ? 1 : 0,
        }),
      });
      Toast.success(enabled ? "已启用" : "已禁用");
      void load();
    } catch (e) {
      Toast.error(e instanceof ApiError ? e.message : "操作失败");
    }
  };

  const deleteOne = async (id: number) => {
    try {
      await request(`/api/tokens/${id}`, { method: "DELETE" });
      Toast.success("已删除");
      void load();
    } catch (e) {
      Toast.error(e instanceof ApiError ? e.message : "删除失败");
    }
  };

  const batchDelete = async () => {
    if (selectedKeys.length === 0) return;
    try {
      await request("/api/tokens/batch-delete", {
        method: "POST",
        body: JSON.stringify({ ids: selectedKeys }),
      });
      Toast.success("批量删除成功");
      setSelectedKeys([]);
      void load();
    } catch (e) {
      Toast.error(e instanceof ApiError ? e.message : "批量删除失败");
    }
  };

  const openPlayground = (row: ApiTokenRow) => {
    sessionStorage.setItem("playground_token_id", String(row.id));
    sessionStorage.setItem("playground_token_name", row.name);
    router.push("/console/playground");
  };

  const columns = useMemo(
    () => [
      {
        title: "名称",
        dataIndex: "name",
      },
      {
        title: "状态",
        dataIndex: "status",
        render: (_: unknown, row: ApiTokenRow) => (
          <Tag color={row.enabled ? "green" : "grey"}>
            {row.enabled ? "已启用" : "已禁用"}
          </Tag>
        ),
      },
      {
        title: "密钥",
        dataIndex: "key_masked",
        render: (_: unknown, row: ApiTokenRow) => (
          <Button
            theme="borderless"
            icon={<IconCopy />}
            onClick={() => void copyText(row.key_masked, "掩码密钥")}
          >
            {row.key_masked}
          </Button>
        ),
      },
      {
        title: "分组",
        dataIndex: "token_group",
        render: (v: string) => <Tag>{v}</Tag>,
      },
      {
        title: "额度",
        render: (_: unknown, row: ApiTokenRow) =>
          row.quota_limit == null
            ? "无限"
            : `${row.used_quota} / ${row.quota_limit}`,
      },
      {
        title: "模型",
        render: (_: unknown, row: ApiTokenRow) =>
          row.models?.length ? row.models.join(", ") : "全部",
      },
      {
        title: "操作",
        render: (_: unknown, row: ApiTokenRow) => (
          <div className="flex flex-wrap gap-1">
            <Dropdown
              trigger="click"
              render={
                <Dropdown.Menu>
                  <Dropdown.Item onClick={() => openPlayground(row)}>聊天</Dropdown.Item>
                  <Dropdown.Item onClick={() => openEdit(row)}>编辑</Dropdown.Item>
                  <Dropdown.Item
                    onClick={() => void toggleEnabled(row, !row.enabled)}
                  >
                    {row.enabled ? "禁用" : "启用"}
                  </Dropdown.Item>
                  <Dropdown.Divider />
                  <Dropdown.Item type="danger" onClick={() => void deleteOne(row.id)}>
                    删除
                  </Dropdown.Item>
                </Dropdown.Menu>
              }
            >
              <Button size="small">操作</Button>
            </Dropdown>
          </div>
        ),
      },
    ],
    [],
  );

  const initialValues: TokenFormValues = editing
    ? {
        name: editing.name,
        token_group: editing.token_group,
        quota_unlimited: editing.quota_limit == null,
        quota_limit: editing.quota_limit ?? undefined,
        models_text: editing.models?.join("\n") ?? "",
        ip_whitelist_text: editing.ip_whitelist?.join("\n") ?? "",
      }
    : {
        name: "",
        token_group: "default",
        quota_unlimited: true,
        models_text: "",
        ip_whitelist_text: "",
      };

  return (
    <div className="space-y-4">
      <div className="flex flex-wrap items-center justify-between gap-3">
        <Title heading={4}>令牌管理</Title>
        <div className="flex gap-2">
          <Button
            type="danger"
            disabled={selectedKeys.length === 0}
            onClick={() => void batchDelete()}
          >
            删除所选
          </Button>
          <Button theme="solid" type="primary" icon={<IconPlus />} onClick={openCreate}>
            添加令牌
          </Button>
        </div>
      </div>

      <Table
        rowKey="id"
        columns={columns}
        dataSource={items}
        loading={loading}
        pagination={{
          currentPage: page,
          pageSize: 10,
          total,
          onPageChange: setPage,
        }}
        rowSelection={{
          selectedRowKeys: selectedKeys,
          onChange: (keys) => setSelectedKeys(keys as number[]),
        }}
      />

      <Modal
        title={editing ? "编辑令牌" : "添加令牌"}
        visible={modalOpen}
        onCancel={() => setModalOpen(false)}
        footer={null}
        width={520}
      >
        <Form
          key={editing?.id ?? "new"}
          initValues={initialValues}
          onSubmit={async (values) => {
            const v = values as TokenFormValues;
            const models = v.models_text
              ? v.models_text.split(/[\n,]/).map((s) => s.trim()).filter(Boolean)
              : [];
            const ipWhitelist = v.ip_whitelist_text
              ? v.ip_whitelist_text.split(/[\n,]/).map((s) => s.trim()).filter(Boolean)
              : [];
            const body: Record<string, unknown> = {
              name: v.name,
              token_group: v.token_group,
              models,
              ip_whitelist: ipWhitelist,
            };
            if (!v.quota_unlimited && v.quota_limit != null) {
              body.quota_limit = v.quota_limit;
            }
            try {
              if (editing) {
                const updateBody: Record<string, unknown> = {
                  name: v.name,
                  token_group: v.token_group,
                  models,
                  ip_whitelist: ipWhitelist,
                };
                if (v.quota_unlimited) {
                  updateBody.clear_quota_limit = true;
                } else if (v.quota_limit != null) {
                  updateBody.quota_limit = v.quota_limit;
                }
                await request(`/api/tokens/${editing.id}`, {
                  method: "PUT",
                  body: JSON.stringify(updateBody),
                });
                Toast.success("已保存");
                setModalOpen(false);
                void load();
              } else {
                const res = await request<TokenCreateResponse>("/api/tokens", {
                  method: "POST",
                  body: JSON.stringify(body),
                });
                setModalOpen(false);
                setSecretModal({
                  key: res.data.api_key,
                  name: res.data.token.name,
                });
                void load();
              }
            } catch (e) {
              Toast.error(e instanceof ApiError ? e.message : "保存失败");
            }
          }}
        >
          <Form.Input field="name" label="名称" rules={[{ required: true, message: "必填" }]} />
          <Form.Select
            field="token_group"
            label="分组"
            optionList={groups.map((g) => ({
              value: g.slug,
              label: `${g.name}（${g.multiplier}x）`,
            }))}
          />
          <Form.Switch field="quota_unlimited" label="无限额度" />
          <Form.InputNumber field="quota_limit" label="额度上限" min={0} />
          <Form.TextArea
            field="models_text"
            label="可用模型（每行一个，空=全部）"
            autosize={{ minRows: 2, maxRows: 6 }}
          />
          <Form.TextArea
            field="ip_whitelist_text"
            label="IP 白名单（每行一个，空=不限制）"
            autosize={{ minRows: 2, maxRows: 6 }}
          />
          <Button htmlType="submit" type="primary" theme="solid" block>
            {editing ? "保存" : "创建"}
          </Button>
        </Form>
      </Modal>

      <Modal
        title="请妥善保存密钥"
        visible={secretModal != null}
        onCancel={() => setSecretModal(null)}
        footer={null}
        width={480}
      >
        {secretModal && (
          <div className="space-y-4 text-center">
            <Text type="secondary">令牌「{secretModal.name}」仅显示一次</Text>
            <div className="break-all rounded-lg bg-[var(--portal-bg)] p-3 font-mono text-sm">
              {secretModal.key}
            </div>
            <div className="flex justify-center gap-2">
              <Button
                icon={<IconCopy />}
                onClick={() => void copyText(secretModal.key, "密钥")}
              >
                复制
              </Button>
            </div>
            <div className="flex justify-center">
              <QRCodeSVG value={secretModal.key} size={160} />
            </div>
          </div>
        )}
      </Modal>
    </div>
  );
}
