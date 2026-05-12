"use client";

import {
  App,
  Button,
  Card,
  Col,
  Form,
  Input,
  InputNumber,
  Row,
  Space,
  Table,
  Tabs,
} from "antd";
import type { ColumnsType } from "antd/es/table";
import { useCallback, useEffect, useState } from "react";

import { ApiError } from "@/lib/api/client";
import type {
  BatchRedeemResponse,
  RedeemCodeRow,
  RedeemListResponse,
  RedeemStatsResponse,
} from "@/lib/api/types";
import { useGatewayRequest } from "@/lib/hooks/use-gateway-request";

export default function AdminRedeemPage() {
  const request = useGatewayRequest();
  const { message } = App.useApp();
  const [stats, setStats] = useState<Record<string, number>>({});
  const [listLoading, setListLoading] = useState(false);
  const [rows, setRows] = useState<RedeemCodeRow[]>([]);
  const [total, setTotal] = useState(0);
  const [page, setPage] = useState(1);
  const [pageSize, setPageSize] = useState(50);
  const [statusFilter, setStatusFilter] = useState<number | undefined>(undefined);
  const [batchLoading, setBatchLoading] = useState(false);
  const [codesResult, setCodesResult] = useState<string[] | null>(null);
  const [batchForm] = Form.useForm<{
    count: number;
    quota: number;
    expires_days?: number;
  }>();

  const loadStats = useCallback(async () => {
    try {
      const data = await request<RedeemStatsResponse>("/admin/redeem/stats");
      setStats(data.by_status ?? {});
    } catch (e) {
      if (e instanceof ApiError) {
        message.error(e.message || "统计加载失败");
      }
    }
  }, [request, message]);

  const loadList = useCallback(async () => {
    setListLoading(true);
    try {
      const qs = new URLSearchParams({
        page: String(page),
        page_size: String(pageSize),
      });
      if (statusFilter !== undefined) {
        qs.set("status", String(statusFilter));
      }
      const data = await request<RedeemListResponse>(`/admin/redeem/codes?${qs}`);
      setRows(data.items);
      setTotal(data.total);
    } catch (e) {
      if (e instanceof ApiError) {
        message.error(e.message || "列表加载失败");
      }
    } finally {
      setListLoading(false);
    }
  }, [request, message, page, pageSize, statusFilter]);

  useEffect(() => {
    let cancelled = false;
    queueMicrotask(() => {
      if (!cancelled) {
        void loadStats();
      }
    });
    return () => {
      cancelled = true;
    };
  }, [loadStats]);

  useEffect(() => {
    let cancelled = false;
    queueMicrotask(() => {
      if (!cancelled) {
        void loadList();
      }
    });
    return () => {
      cancelled = true;
    };
  }, [loadList]);

  const columns: ColumnsType<RedeemCodeRow> = [
    { title: "ID", dataIndex: "ID", width: 72 },
    {
      title: "兑换码",
      dataIndex: "Code",
      render: (code: string) => maskCode(code),
    },
    { title: "额度", dataIndex: "Quota" },
    { title: "状态", dataIndex: "Status" },
    { title: "使用者", dataIndex: "UsedBy", render: (v) => v ?? "—" },
    {
      title: "过期",
      dataIndex: "ExpiresAt",
      render: (v: string | null | undefined) => v ?? "—",
    },
    {
      title: "创建时间",
      dataIndex: "CreatedAt",
      render: (v: string) => formatTime(v),
    },
  ];

  return (
    <div className="flex flex-col gap-4">
      <h1 className="text-lg font-semibold text-neutral-900">兑换运营</h1>

      <Tabs
        items={[
          {
            key: "stats",
            label: "统计",
            children: (
              <Row gutter={[16, 16]}>
                {Object.entries(stats).map(([k, v]) => (
                  <Col xs={12} md={8} lg={6} key={k}>
                    <Card size="small" title={`状态 ${k}`}>
                      <span className="text-2xl font-semibold">{v}</span>
                    </Card>
                  </Col>
                ))}
                {Object.keys(stats).length === 0 && (
                  <Col span={24}>
                    <Card size="small">暂无统计数据</Card>
                  </Col>
                )}
              </Row>
            ),
          },
          {
            key: "batch",
            label: "批量生成",
            children: (
              <Card className="max-w-xl" size="small">
                <Form
                  form={batchForm}
                  layout="vertical"
                  initialValues={{ count: 10, quota: 1000 }}
                  onFinish={async (v) => {
                    setBatchLoading(true);
                    try {
                      const body: {
                        count: number;
                        quota: number;
                        expires_days?: number;
                      } = { count: v.count, quota: v.quota };
                      if (
                        v.expires_days != null &&
                        !Number.isNaN(v.expires_days) &&
                        v.expires_days > 0
                      ) {
                        body.expires_days = v.expires_days;
                      }
                      const data = await request<BatchRedeemResponse>(
                        "/admin/redeem/batch",
                        {
                          method: "POST",
                          body: JSON.stringify(body),
                        },
                      );
                      setCodesResult(data.codes);
                      message.success(`已生成 ${data.count} 条`);
                      await loadStats();
                      await loadList();
                    } catch (e) {
                      if (e instanceof ApiError) {
                        message.error(e.message || "生成失败");
                      }
                    } finally {
                      setBatchLoading(false);
                    }
                  }}
                >
                  <Form.Item
                    name="count"
                    label="数量"
                    rules={[{ required: true, type: "number", min: 1, max: 5000 }]}
                  >
                    <InputNumber className="w-full" min={1} max={5000} />
                  </Form.Item>
                  <Form.Item
                    name="quota"
                    label="单码额度"
                    rules={[{ required: true, type: "number", min: 1 }]}
                  >
                    <InputNumber className="w-full" min={1} />
                  </Form.Item>
                  <Form.Item name="expires_days" label="有效天数（可选）">
                    <InputNumber className="w-full" min={1} placeholder="不填则不过期" />
                  </Form.Item>
                  <Form.Item>
                    <Button type="primary" htmlType="submit" loading={batchLoading}>
                      生成
                    </Button>
                  </Form.Item>
                </Form>
                {codesResult && codesResult.length > 0 && (
                  <div className="mt-4">
                    <div className="mb-2 flex items-center justify-between gap-2">
                      <span className="text-sm font-medium text-neutral-700">
                        本次生成的兑换码（请立即复制保存）
                      </span>
                      <Button
                        size="small"
                        onClick={() => {
                          void navigator.clipboard.writeText(codesResult.join("\n"));
                          message.success("已复制到剪贴板");
                        }}
                      >
                        复制全部
                      </Button>
                    </div>
                    <Input.TextArea
                      readOnly
                      rows={Math.min(12, codesResult.length)}
                      value={codesResult.join("\n")}
                      className="font-mono text-sm"
                    />
                  </div>
                )}
              </Card>
            ),
          },
          {
            key: "list",
            label: "兑换码列表",
            children: (
              <div className="flex flex-col gap-3">
                <Space wrap>
                  <span className="text-neutral-600">状态筛选：</span>
                  <Button
                    size="small"
                    type={statusFilter === undefined ? "primary" : "default"}
                    onClick={() => {
                      setPage(1);
                      setStatusFilter(undefined);
                    }}
                  >
                    全部
                  </Button>
                  {[1, 2, 3].map((s) => (
                    <Button
                      key={s}
                      size="small"
                      type={statusFilter === s ? "primary" : "default"}
                      onClick={() => {
                        setPage(1);
                        setStatusFilter(s);
                      }}
                    >
                      {s}
                    </Button>
                  ))}
                </Space>
                <Table<RedeemCodeRow>
                  rowKey="ID"
                  loading={listLoading}
                  columns={columns}
                  dataSource={rows}
                  pagination={{
                    current: page,
                    pageSize,
                    total,
                    showSizeChanger: true,
                    pageSizeOptions: [20, 50, 100, 200],
                    onChange: (p, ps) => {
                      setPage(p);
                      setPageSize(ps);
                    },
                  }}
                />
              </div>
            ),
          },
        ]}
      />
    </div>
  );
}

function maskCode(code: string) {
  if (code.length <= 8) {
    return "****";
  }
  return `${code.slice(0, 4)}…${code.slice(-4)}`;
}

function formatTime(iso: string) {
  try {
    return new Date(iso).toLocaleString();
  } catch {
    return iso;
  }
}
