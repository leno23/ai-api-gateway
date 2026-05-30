"use client";

import { Button, Empty, Form, Input, Table, Tag, Toast, Typography } from "@douyinfe/semi-ui";
import { useCallback, useState } from "react";

import type { TaskLogRow } from "@/lib/api/dashboard-types";
import { ApiError } from "@/lib/api/client";
import { useDeferredEffect } from "@/lib/hooks/use-deferred-effect";
import { useGatewayRequest } from "@/lib/hooks/use-gateway-request";

const { Title } = Typography;

type FilterValues = {
  task_id?: string;
  range?: Date[];
};

export function TaskLogPage() {
  const request = useGatewayRequest();
  const [items, setItems] = useState<TaskLogRow[]>([]);
  const [total, setTotal] = useState(0);
  const [page, setPage] = useState(1);
  const [loading, setLoading] = useState(false);
  const [filters, setFilters] = useState<FilterValues>({});

  const load = useCallback(async () => {
    setLoading(true);
    const params = new URLSearchParams({ page: String(page), page_size: "10" });
    if (filters.task_id) params.set("task_id", filters.task_id);
    const [start, end] = filters.range ?? [];
    if (start) params.set("start", start.toISOString());
    if (end) params.set("end", end.toISOString());
    try {
      const res = await request<{
        success: boolean;
        data: { items: TaskLogRow[]; total: number };
      }>(`/api/logs/tasks?${params}`);
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
      <Title heading={4}>任务日志</Title>
      <Form
        layout="horizontal"
        onSubmit={(v) => {
          setFilters(v as FilterValues);
          setPage(1);
        }}
      >
        <Form.DatePicker field="range" label="时间" type="dateTimeRange" style={{ width: 360 }} />
        <Form.Input field="task_id" label="任务 ID" />
        <Button htmlType="submit" type="primary" theme="solid">
          查询
        </Button>
      </Form>

      {!loading && items.length === 0 ? (
        <Empty title="暂无任务记录" description="异步任务提交后将显示在此" />
      ) : (
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
              title: "提交时间",
              dataIndex: "submitted_at",
              render: (v: string) => new Date(v).toLocaleString(),
            },
            {
              title: "结束时间",
              dataIndex: "finished_at",
              render: (v: string | null) => (v ? new Date(v).toLocaleString() : "—"),
            },
            {
              title: "耗时(ms)",
              dataIndex: "duration_ms",
              render: (v: number | null) => v ?? "—",
            },
            { title: "平台", dataIndex: "platform" },
            { title: "类型", dataIndex: "type" },
            { title: "任务 ID", dataIndex: "task_id" },
            {
              title: "状态",
              dataIndex: "status",
              render: (v: string) => <Tag>{v}</Tag>,
            },
            {
              title: "进度",
              dataIndex: "progress",
              render: (v: number) => `${v}%`,
            },
          ]}
        />
      )}
    </div>
  );
}
