"use client";

import {
  App,
  Button,
  Drawer,
  Form,
  Input,
  InputNumber,
  Popconfirm,
  Select,
  Space,
  Table,
} from "antd";
import type { ColumnsType } from "antd/es/table";
import { useCallback, useEffect, useState } from "react";

import { ApiError } from "@/lib/api/client";
import type {
  ChannelListItem,
  ChannelPayload,
  ChannelsListResponse,
} from "@/lib/api/types";
import { useGatewayRequest } from "@/lib/hooks/use-gateway-request";

type ChannelFormValues = {
  name: string;
  provider: string;
  base_url: string;
  api_key?: string;
  models?: string[];
  model_mapping_text?: string;
  priority?: number;
  weight?: number;
  status?: number;
  rate_limit?: number;
};

export default function AdminChannelsPage() {
  const request = useGatewayRequest();
  const { message } = App.useApp();
  const [loading, setLoading] = useState(false);
  const [items, setItems] = useState<ChannelListItem[]>([]);
  const [drawerOpen, setDrawerOpen] = useState(false);
  const [editing, setEditing] = useState<ChannelListItem | null>(null);
  const [form] = Form.useForm<ChannelFormValues>();

  const load = useCallback(async () => {
    setLoading(true);
    try {
      const data = await request<ChannelsListResponse>("/admin/channels");
      setItems(data.items);
    } catch (e) {
      if (e instanceof ApiError) {
        message.error(e.message || "加载失败");
      }
    } finally {
      setLoading(false);
    }
  }, [request, message]);

  useEffect(() => {
    let cancelled = false;
    queueMicrotask(() => {
      if (!cancelled) {
        void load();
      }
    });
    return () => {
      cancelled = true;
    };
  }, [load]);

  const openCreate = () => {
    setEditing(null);
    form.resetFields();
    form.setFieldsValue({
      models: [],
      priority: 0,
      weight: 100,
      status: 1,
      rate_limit: 0,
      model_mapping_text: "",
    });
    setDrawerOpen(true);
  };

  const openEdit = (row: ChannelListItem) => {
    setEditing(row);
    form.setFieldsValue({
      name: row.name,
      provider: row.provider,
      base_url: row.base_url,
      api_key: "",
      models: row.models ?? [],
      model_mapping_text: row.model_mapping
        ? JSON.stringify(row.model_mapping, null, 2)
        : "",
      priority: row.priority,
      weight: row.weight,
      status: row.status,
      rate_limit: row.rate_limit,
    });
    setDrawerOpen(true);
  };

  const submit = async () => {
    try {
      const v = await form.validateFields();
      let model_mapping: Record<string, string> | undefined;
      const raw = v.model_mapping_text?.trim();
      if (raw) {
        try {
          model_mapping = JSON.parse(raw) as Record<string, string>;
          if (typeof model_mapping !== "object" || model_mapping === null) {
            throw new Error("not object");
          }
        } catch {
          message.error("model_mapping 不是合法 JSON 对象");
          return;
        }
      }

      const payload: ChannelPayload = {
        name: v.name,
        provider: v.provider,
        base_url: v.base_url,
        api_key: v.api_key ?? "",
        models: v.models ?? [],
        model_mapping,
        priority: v.priority,
        weight: v.weight,
        status: v.status,
        rate_limit: v.rate_limit,
      };

      if (editing) {
        if (!payload.api_key.trim()) {
          message.error("更新渠道时必须重新填写完整 API Key");
          return;
        }
        await request(`/admin/channels/${editing.id}`, {
          method: "PUT",
          body: JSON.stringify(payload),
        });
        message.success("已保存");
      } else {
        if (!payload.api_key.trim()) {
          message.error("请填写 API Key");
          return;
        }
        await request("/admin/channels", {
          method: "POST",
          body: JSON.stringify(payload),
        });
        message.success("已创建");
      }
      setDrawerOpen(false);
      await load();
    } catch (e) {
      if (e instanceof ApiError) {
        message.error(e.message || "保存失败");
      }
    }
  };

  const columns: ColumnsType<ChannelListItem> = [
    { title: "ID", dataIndex: "id", width: 72 },
    { title: "名称", dataIndex: "name" },
    { title: "提供商", dataIndex: "provider", width: 100 },
    { title: "Base URL", dataIndex: "base_url", ellipsis: true },
    { title: "Key 预览", dataIndex: "api_key_preview", width: 120 },
    { title: "优先级", dataIndex: "priority", width: 80 },
    { title: "权重", dataIndex: "weight", width: 80 },
    { title: "状态", dataIndex: "status", width: 72 },
    { title: "限流", dataIndex: "rate_limit", width: 80 },
    {
      title: "操作",
      key: "actions",
      width: 160,
      render: (_, row) => (
        <Space>
          <Button type="link" size="small" onClick={() => openEdit(row)}>
            编辑
          </Button>
          <Popconfirm
            title="确定删除该渠道？"
            okText="删除"
            cancelText="取消"
            okButtonProps={{ danger: true }}
            onConfirm={async () => {
              try {
                await request(`/admin/channels/${row.id}`, { method: "DELETE" });
                message.success("已删除");
                await load();
              } catch (e) {
                if (e instanceof ApiError) {
                  message.error(e.message || "删除失败");
                }
              }
            }}
          >
            <Button type="link" size="small" danger>
              删除
            </Button>
          </Popconfirm>
        </Space>
      ),
    },
  ];

  return (
    <div className="flex flex-col gap-4">
      <div className="flex justify-between gap-4">
        <h1 className="text-lg font-semibold text-neutral-900">渠道管理</h1>
        <Button type="primary" onClick={openCreate}>
          新建渠道
        </Button>
      </div>
      <Table<ChannelListItem>
        rowKey="id"
        loading={loading}
        columns={columns}
        dataSource={items}
        scroll={{ x: 1100 }}
        pagination={false}
      />

      <Drawer
        title={editing ? `编辑渠道 #${editing.id}` : "新建渠道"}
        width={480}
        open={drawerOpen}
        onClose={() => setDrawerOpen(false)}
        extra={
          <Space>
            <Button onClick={() => setDrawerOpen(false)}>取消</Button>
            <Button type="primary" onClick={() => void submit()}>
              保存
            </Button>
          </Space>
        }
        styles={{ body: { paddingBottom: 48 } }}
      >
        <Form form={form} layout="vertical" className="max-w-full">
          <Form.Item name="name" label="名称" rules={[{ required: true }]}>
            <Input />
          </Form.Item>
          <Form.Item name="provider" label="提供商" rules={[{ required: true }]}>
            <Input placeholder="openai / custom" />
          </Form.Item>
          <Form.Item name="base_url" label="Base URL" rules={[{ required: true }]}>
            <Input placeholder="https://api.openai.com/v1" />
          </Form.Item>
          <Form.Item
            name="api_key"
            label={editing ? "API Key（更新时必填完整密钥）" : "API Key"}
            rules={editing ? [] : [{ required: true, message: "请输入 API Key" }]}
          >
            <Input.Password placeholder={editing ? "重新填写完整 Key" : ""} />
          </Form.Item>
          <Form.Item name="models" label="模型列表（标签输入）">
            <Select
              mode="tags"
              className="w-full"
              placeholder="输入模型名后回车"
              tokenSeparators={[",", " "]}
            />
          </Form.Item>
          <Form.Item
            name="model_mapping_text"
            label="model_mapping（JSON，可选）"
            tooltip='例如 {"gpt-4o":"upstream-name"}，留空表示不传映射'
          >
            <Input.TextArea rows={4} placeholder="{}" />
          </Form.Item>
          <Form.Item name="priority" label="优先级">
            <InputNumber className="w-full" />
          </Form.Item>
          <Form.Item name="weight" label="权重">
            <InputNumber className="w-full" />
          </Form.Item>
          <Form.Item name="status" label="状态">
            <InputNumber className="w-full" placeholder="1 启用 / 0 停用" />
          </Form.Item>
          <Form.Item name="rate_limit" label="限流（RPM）">
            <InputNumber className="w-full" min={0} />
          </Form.Item>
        </Form>
      </Drawer>
    </div>
  );
}
