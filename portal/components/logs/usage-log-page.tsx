"use client";

import { Button, Form, Table, Tag, Toast, Typography } from "@douyinfe/semi-ui";
import { useCallback, useState } from "react";

import type { UsageLogRow } from "@/lib/api/dashboard-types";
import { ApiError } from "@/lib/api/client";
import { useDeferredEffect } from "@/lib/hooks/use-deferred-effect";
import { useGatewayRequest } from "@/lib/hooks/use-gateway-request";

const { Title } = Typography;

type FilterValues = {
  token_name?: string;
  model?: string;
  request_id?: string;
  token_group?: string;
  range?: Date[];
};

function toISO(d?: Date) {
  return d ? d.toISOString() : undefined;
}

export function UsageLogPage() {
  const request = useGatewayRequest();
  const [items, setItems] = useState<UsageLogRow[]>([]);
  const [total, setTotal] = useState(0);
  const [page, setPage] = useState(1);
  const [loading, setLoading] = useState(false);
  const [filters, setFilters] = useState<FilterValues>({});

  const load = useCallback(async () => {
    setLoading(true);
    const params = new URLSearchParams({
      page: String(page),
      page_size: "10",
    });
    if (filters.token_name) params.set("token_name", filters.token_name);
    if (filters.model) params.set("model", filters.model);
    if (filters.request_id) params.set("request_id", filters.request_id);
    if (filters.token_group) params.set("token_group", filters.token_group);
    const [start, end] = filters.range ?? [];
    if (start) params.set("start", toISO(start)!);
    if (end) params.set("end", toISO(end)!);
    try {
      const res = await request<{
        success: boolean;
        data: { items: UsageLogRow[]; total: number };
      }>(`/api/logs/usage?${params}`);
      setItems(res.data.items);
      setTotal(res.data.total);
    } catch (e) {
      Toast.error(e instanceof ApiError ? e.message : "加载失败");
    } finally {
      setLoading(false);
    }
  }, [request, page, filters]);

  useDeferredEffect(() => load(), [load]);

  return (
    <div className="space-y-4">
      <Title heading={4}>使用日志</Title>
      <Form
        layout="horizontal"
        onSubmit={(v) => {
          setFilters(v as FilterValues);
          setPage(1);
        }}
      >
        <Form.DatePicker
          field="range"
          label="时间"
          type="dateTimeRange"
          style={{ width: 360 }}
        />
        <Form.Input field="token_name" label="令牌" placeholder="令牌名称" />
        <Form.Input field="model" label="模型" placeholder="模型名" />
        <Form.Input field="request_id" label="Request ID" />
        <Form.Input field="token_group" label="分组" />
        <Button htmlType="submit" type="primary" theme="solid">
          查询
        </Button>
      </Form>

      <Table
        rowKey="id"
        loading={loading}
        dataSource={items}
        pagination={{
          currentPage: page,
          pageSize: 10,
          total,
          onPageChange: setPage,
        }}
        columns={[
          {
            title: "时间",
            dataIndex: "created_at",
            width: 180,
            render: (v: string) => new Date(v).toLocaleString(),
          },
          { title: "令牌", dataIndex: "token_name", width: 120 },
          {
            title: "分组",
            dataIndex: "token_group",
            render: (v: string) => <Tag>{v || "—"}</Tag>,
          },
          { title: "模型", dataIndex: "model" },
          { title: "Request ID", dataIndex: "request_id", width: 200 },
          {
            title: "首字(ms)",
            dataIndex: "time_to_first_ms",
            render: (v: number | null) => v ?? "—",
          },
          { title: "输入 Tok", dataIndex: "prompt_tokens" },
          {
            title: "消耗额度",
            dataIndex: "cost_quota",
          },
        ]}
      />
    </div>
  );
}
