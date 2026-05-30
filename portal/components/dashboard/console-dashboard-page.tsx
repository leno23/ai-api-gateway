"use client";

import Link from "next/link";
import {
  Button,
  Card,
  Collapse,
  TabPane,
  Tabs,
  Tag,
  Toast,
  Typography,
} from "@douyinfe/semi-ui";
import { IconCopy, IconExternalOpen } from "@douyinfe/semi-icons";
import { useCallback, useMemo, useState } from "react";
import {
  Bar,
  BarChart,
  CartesianGrid,
  Line,
  LineChart,
  ResponsiveContainer,
  Tooltip,
  XAxis,
  YAxis,
} from "recharts";

import type {
  DashboardChartsResponse,
  DashboardStatsResponse,
  NodePingResult,
  NodesPingResponse,
  PortalNode,
} from "@/lib/api/dashboard-types";
import { DashboardAnnouncementsCard } from "@/components/dashboard/dashboard-announcements-card";
import { StatSparkline } from "@/components/dashboard/stat-sparkline";
import { AnnouncementModal } from "@/components/announcement/announcement-modal";
import { usePortalLocale } from "@/contexts/portal-locale-context";
import { ApiError } from "@/lib/api/client";
import { useAuth } from "@/contexts/auth-context";
import { useDeferredEffect } from "@/lib/hooks/use-deferred-effect";
import { useGatewayRequest } from "@/lib/hooks/use-gateway-request";

const { Title, Text } = Typography;

const SERVICES = [
  { name: "API 网关", status: "ok" as const },
  { name: "OpenAI 兼容接口", status: "ok" as const },
  { name: "操练场", status: "ok" as const },
];

function greetingLabel(): string {
  const h = new Date().getHours();
  if (h < 6) return "凌晨好";
  if (h < 12) return "早上好";
  if (h < 18) return "下午好";
  return "晚上好";
}

export function ConsoleDashboardPage() {
  const request = useGatewayRequest();
  const { user } = useAuth();
  const { messages: dm } = usePortalLocale();
  const [stats, setStats] = useState<DashboardStatsResponse["data"] | null>(null);
  const [charts, setCharts] = useState<DashboardChartsResponse["data"] | null>(null);
  const [nodes, setNodes] = useState<PortalNode[]>([]);
  const [pingMap, setPingMap] = useState<Record<string, NodePingResult>>({});
  const [pinging, setPinging] = useState(false);
  const [range, setRange] = useState("7d");
  const [loading, setLoading] = useState(true);

  const load = useCallback(async () => {
    setLoading(true);
    try {
      const [s, c, n] = await Promise.all([
        request<DashboardStatsResponse>("/api/dashboard/stats"),
        request<DashboardChartsResponse>(`/api/dashboard/charts?range=${range}`),
        request<{ success: boolean; data: PortalNode[] }>("/api/dashboard/nodes"),
      ]);
      setStats(s.data);
      setCharts(c.data);
      setNodes(n.data);
    } catch (e) {
      Toast.error(e instanceof ApiError ? e.message : "加载看板失败");
    } finally {
      setLoading(false);
    }
  }, [request, range]);

  useDeferredEffect(() => load(), [load]);

  const runPing = useCallback(async () => {
    setPinging(true);
    try {
      const res = await request<NodesPingResponse>("/api/nodes/ping");
      const next: Record<string, NodePingResult> = {};
      for (const row of res.data) {
        next[row.url] = row;
      }
      setPingMap(next);
    } catch (e) {
      Toast.error(e instanceof ApiError ? e.message : "测速失败");
    } finally {
      setPinging(false);
    }
  }, [request]);

  const metricCards = useMemo(() => {
    if (!stats) return [];
    return [
      {
        key: "account",
        title: "账户余额",
        value: stats.account.quota,
        sub: `已消耗 ${stats.account.used_quota} · 分组 ${stats.account.group}`,
        action: (
          <Link href="/console/wallet">
            <Button size="small" theme="solid" type="primary">
              充值
            </Button>
          </Link>
        ),
      },
      {
        key: "requests",
        title: "24h 请求",
        value: stats.usage.request_count_24h,
        spark: stats.sparklines.requests,
        color: "#007AFF",
      },
      {
        key: "tokens",
        title: "24h Tokens",
        value: stats.usage.total_tokens_24h,
      },
      {
        key: "cost",
        title: "24h 消耗额度",
        value: stats.consumption.cost_quota_24h,
        spark: stats.sparklines.cost,
        color: "#7C3AED",
      },
    ];
  }, [stats]);

  const trendData = useMemo(
    () =>
      (charts?.call_trend ?? []).map((p) => ({
        hour: p.hour.slice(11, 16),
        count: p.count,
        cost: p.cost,
      })),
    [charts],
  );

  const modelData = useMemo(
    () => charts?.consumption_by_model?.slice(0, 10) ?? [],
    [charts],
  );

  return (
    <div className="flex flex-col gap-6 lg:flex-row">
      <AnnouncementModal placement="console" />
      <div className="min-w-0 flex-1 space-y-6">
        <div className="flex flex-wrap items-center justify-between gap-3">
          <div>
            <Title heading={4} className="!mb-0">
              👋 {greetingLabel()}
              {user?.display_name || user?.username
                ? `，${user.display_name || user.username}`
                : ""}
            </Title>
            <Text type="tertiary" size="small" className="mt-1 block">
              数据看板 · RPM {stats?.performance.rpm ?? "—"} · TPM{" "}
              {stats?.performance.tpm ?? "—"} · 平均延迟{" "}
              {stats?.performance.avg_latency_ms ?? "—"} ms
            </Text>
          </div>
          <Tabs type="button" activeKey={range} onChange={setRange}>
            <TabPane tab="24h" itemKey="24h" />
            <TabPane tab="7d" itemKey="7d" />
            <TabPane tab="30d" itemKey="30d" />
          </Tabs>
        </div>

        <div className="grid gap-4 md:grid-cols-2 xl:grid-cols-4">
          {metricCards.map((card) => (
            <Card key={card.key} loading={loading} className="!rounded-xl" shadows="hover">
              <div className="flex items-start justify-between gap-2">
                <Text type="tertiary" size="small">
                  {card.title}
                </Text>
                {"action" in card ? card.action : null}
              </div>
              <Title heading={3} className="!mt-2 !mb-0">
                {card.value}
              </Title>
              {"sub" in card && card.sub ? (
                <Text type="tertiary" size="small" className="mt-1 block">
                  {card.sub}
                </Text>
              ) : null}
              {"spark" in card && card.spark ? (
                <StatSparkline data={card.spark} color={card.color} />
              ) : null}
            </Card>
          ))}
        </div>

        <Card className="!rounded-xl" loading={loading}>
          <Tabs type="line">
            <TabPane tab="消耗分布" itemKey="dist">
              <div className="h-64 pt-4">
                <ResponsiveContainer width="100%" height="100%">
                  <BarChart data={modelData}>
                    <CartesianGrid strokeDasharray="3 3" />
                    <XAxis dataKey="model" tick={{ fontSize: 11 }} />
                    <YAxis />
                    <Tooltip />
                    <Bar dataKey="cost" fill="#007AFF" name="消耗额度" />
                  </BarChart>
                </ResponsiveContainer>
              </div>
            </TabPane>
            <TabPane tab="调用趋势" itemKey="trend">
              <div className="h-64 pt-4">
                <ResponsiveContainer width="100%" height="100%">
                  <LineChart data={trendData}>
                    <CartesianGrid strokeDasharray="3 3" />
                    <XAxis dataKey="hour" />
                    <YAxis />
                    <Tooltip />
                    <Line type="monotone" dataKey="count" stroke="#007AFF" name="次数" />
                  </LineChart>
                </ResponsiveContainer>
              </div>
            </TabPane>
            <TabPane tab="排行" itemKey="rank">
              <div className="space-y-2 pt-4">
                {(charts?.ranking ?? []).slice(0, 8).map((row) => (
                  <div key={row.model} className="flex justify-between text-sm">
                    <Text>{row.model}</Text>
                    <Tag>{row.cost}</Tag>
                  </div>
                ))}
              </div>
            </TabPane>
          </Tabs>
        </Card>
      </div>

      <aside className="w-full shrink-0 space-y-4 lg:w-80">
        <Card
          title="API 信息"
          className="!rounded-xl"
          headerExtraContent={
            <Button size="small" loading={pinging} onClick={() => void runPing()}>
              {pinging ? dm.dashboard.pinging : dm.dashboard.ping}
            </Button>
          }
        >
          {nodes.length === 0 && !loading ? (
            <Text type="tertiary" size="small">
              暂无节点配置
            </Text>
          ) : null}
          {nodes.map((node) => {
            const ping = pingMap[node.url];
            return (
              <div
                key={node.url}
                className="mb-3 flex items-start justify-between gap-2 border-b border-[var(--semi-color-border)] pb-3 last:mb-0 last:border-0 last:pb-0"
              >
                <div className="min-w-0 flex-1">
                  <div className="flex flex-wrap items-center gap-2">
                    <Text strong>{node.name}</Text>
                    {node.region ? (
                      <Tag size="small" color="grey">
                        {node.region}
                      </Tag>
                    ) : null}
                    {ping ? (
                      <Tag color={ping.ok ? "green" : "red"} size="small">
                        {ping.ok
                          ? `${dm.dashboard.pingOk} ${ping.latency_ms}ms`
                          : dm.dashboard.pingFail}
                      </Tag>
                    ) : null}
                  </div>
                  <Text type="tertiary" size="small" className="break-all">
                    {node.url}
                  </Text>
                </div>
                <div className="flex shrink-0 gap-1">
                  <Button
                    size="small"
                    icon={<IconCopy />}
                    aria-label="复制"
                    onClick={() => {
                      void navigator.clipboard.writeText(node.url);
                      Toast.success("已复制");
                    }}
                  />
                  <Button
                    size="small"
                    icon={<IconExternalOpen />}
                    aria-label="跳转"
                    onClick={() => window.open(node.url, "_blank", "noopener,noreferrer")}
                  />
                </div>
              </div>
            );
          })}
        </Card>

        <DashboardAnnouncementsCard placement="console" />

        <Card title="FAQ" className="!rounded-xl">
          <Collapse accordion>
            <Collapse.Panel header="如何创建 API 密钥？" itemKey="1">
              进入「令牌管理」→「添加令牌」，创建后请立即保存完整密钥，关闭弹窗后将无法再次查看。
            </Collapse.Panel>
            <Collapse.Panel header="402 额度不足怎么办？" itemKey="2">
              前往「钱包管理」使用兑换码充值，或联系管理员调整额度与分组倍率。
            </Collapse.Panel>
            <Collapse.Panel header="如何切换模型与分组定价？" itemKey="3">
              在「模型广场」左侧筛选供应商与令牌分组；操练场与 API 调用需使用对应分组下的令牌。
            </Collapse.Panel>
            <Collapse.Panel header="节点测速显示失败？" itemKey="4">
              可能是网络或节点 URL 不可达；可点击「跳转」在浏览器中直接访问节点地址排查。
            </Collapse.Panel>
          </Collapse>
        </Card>

        <Card title="服务可用性" className="!rounded-xl">
          <ul className="space-y-2">
            {SERVICES.map((svc) => (
              <li key={svc.name} className="flex items-center justify-between text-sm">
                <Text>{svc.name}</Text>
                <Tag color="green" size="small">
                  正常
                </Tag>
              </li>
            ))}
          </ul>
          <Text type="tertiary" size="small" className="mt-3 block">
            上游模型供应商状态将随版本迭代接入自动探测。
          </Text>
        </Card>
      </aside>
    </div>
  );
}
