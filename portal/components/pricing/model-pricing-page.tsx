"use client";

import {
  Button,
  Card,
  Input,
  Radio,
  RadioGroup,
  Switch,
  Table,
  Tag,
  Toast,
  Typography,
} from "@douyinfe/semi-ui";
import { IconCopy, IconSearch } from "@douyinfe/semi-icons";
import { useCallback, useEffect, useMemo, useState } from "react";

import { PricingFilterSidebar } from "@/components/pricing/pricing-filter-sidebar";
import type { CatalogListResponse, CatalogModelItem } from "@/lib/api/catalog-types";
import { gatewayFetch } from "@/lib/api/client";
import { formatQuotaPerMillion } from "@/lib/format-quota";
import { useDeferredEffect } from "@/lib/hooks/use-deferred-effect";

const { Title, Text, Paragraph } = Typography;

type ViewMode = "grid" | "table";
type SizeMode = "M" | "L";

function buildQuery(params: Record<string, string | number | undefined>) {
  const q = new URLSearchParams();
  Object.entries(params).forEach(([k, v]) => {
    if (v !== undefined && v !== "") q.set(k, String(v));
  });
  const s = q.toString();
  return s ? `?${s}` : "";
}

export function ModelPricingPage() {
  const [loading, setLoading] = useState(true);
  const [items, setItems] = useState<CatalogModelItem[]>([]);
  const [total, setTotal] = useState(0);
  const [providers, setProviders] = useState<string[]>([]);
  const [tokenGroups, setTokenGroups] = useState<
    { slug: string; name: string; multiplier: number }[]
  >([]);

  const [provider, setProvider] = useState<string>("");
  const [endpointType, setEndpointType] = useState<string>("");
  const [billingType, setBillingType] = useState<string>("");
  const [tag, setTag] = useState<string>("");
  const [tokenGroup, setTokenGroup] = useState("default");
  const [search, setSearch] = useState("");
  const [debouncedQ, setDebouncedQ] = useState("");
  const [page, setPage] = useState(1);
  const [showPrice, setShowPrice] = useState(true);
  const [showMultiplier, setShowMultiplier] = useState(true);
  const [viewMode, setViewMode] = useState<ViewMode>("grid");
  const [sizeMode, setSizeMode] = useState<SizeMode>("M");

  useEffect(() => {
    const t = setTimeout(() => setDebouncedQ(search.trim()), 300);
    return () => clearTimeout(t);
  }, [search]);

  const load = useCallback(async () => {
    setLoading(true);
    try {
      const res = await gatewayFetch<CatalogListResponse>(
        `/api/models${buildQuery({
          page,
          page_size: 24,
          provider: provider || undefined,
          endpoint_type: endpointType || undefined,
          billing_type: billingType || undefined,
          tag: tag || undefined,
          token_group: tokenGroup,
          q: debouncedQ || undefined,
        })}`,
      );
      setItems(res.data.items);
      setTotal(res.data.total);
      setProviders(res.meta.providers ?? []);
      setTokenGroups(res.meta.token_groups ?? []);
    } catch {
      Toast.error("加载模型列表失败");
    } finally {
      setLoading(false);
    }
  }, [
    page,
    provider,
    endpointType,
    billingType,
    tag,
    tokenGroup,
    debouncedQ,
  ]);

  useDeferredEffect(() => load(), [load]);

  const resetFilters = () => {
    setProvider("");
    setEndpointType("");
    setBillingType("");
    setTag("");
    setTokenGroup("default");
    setSearch("");
    setDebouncedQ("");
    setPage(1);
  };

  const onProviderChange = (v: string) => {
    setProvider(v);
    setPage(1);
  };
  const onTokenGroupChange = (v: string) => {
    setTokenGroup(v);
    setPage(1);
  };
  const onEndpointTypeChange = (v: string) => {
    setEndpointType(v);
    setPage(1);
  };
  const onBillingTypeChange = (v: string) => {
    setBillingType(v);
    setPage(1);
  };
  const onTagChange = (v: string) => {
    setTag(v);
    setPage(1);
  };

  const priceFor = (row: CatalogModelItem) =>
    showMultiplier ? row.prices_applied : row.prices;

  const copyList = () => {
    const text = items.map((i) => i.model).join("\n");
    void navigator.clipboard.writeText(text);
    Toast.success("已复制模型列表");
  };

  const columns = useMemo(
    () => [
      {
        title: "模型",
        dataIndex: "model",
        render: (_: unknown, row: CatalogModelItem) => (
          <div>
            <Text strong>{row.display_name}</Text>
            <br />
            <Text type="tertiary" size="small">
              {row.model}
            </Text>
          </div>
        ),
      },
      {
        title: "供应商",
        dataIndex: "provider",
        render: (v: string) => <Tag color="blue">{v}</Tag>,
      },
      {
        title: "端点",
        dataIndex: "endpoint_type",
      },
      {
        title: "计费",
        dataIndex: "billing_label",
      },
      ...(showPrice
        ? [
            {
              title: "输入 / 1M",
              render: (_: unknown, row: CatalogModelItem) => {
                const p = priceFor(row);
                return row.billing_type === 2
                  ? formatQuotaPerMillion(p.unit_price)
                  : formatQuotaPerMillion(p.input_per_million);
              },
            },
            {
              title: "输出 / 1M",
              render: (_: unknown, row: CatalogModelItem) =>
                row.billing_type === 2
                  ? "—"
                  : formatQuotaPerMillion(priceFor(row).output_per_million),
            },
          ]
        : []),
    ],
    [showPrice, showMultiplier, priceFor],
  );

  const gridClass =
    sizeMode === "L"
      ? "grid gap-4 sm:grid-cols-2 xl:grid-cols-3"
      : "grid gap-3 sm:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4";

  return (
    <div className="mx-auto flex max-w-7xl gap-6 p-6">
      <PricingFilterSidebar
        total={total}
        providers={providers}
        tokenGroups={tokenGroups}
        provider={provider}
        tokenGroup={tokenGroup}
        endpointType={endpointType}
        billingType={billingType}
        tag={tag}
        onProviderChange={onProviderChange}
        onTokenGroupChange={onTokenGroupChange}
        onEndpointTypeChange={onEndpointTypeChange}
        onBillingTypeChange={onBillingTypeChange}
        onTagChange={onTagChange}
        onReset={resetFilters}
      />

      <div className="min-w-0 flex-1">
        <div className="portal-pricing-banner portal-card mb-6 overflow-hidden px-6 py-8 text-white shadow-md">
          <div className="flex flex-wrap items-center justify-between gap-4">
            <div>
              <Title heading={3} className="!text-white !mb-1">
                {provider ? provider : "全部供应商"}
              </Title>
              <Text className="!text-white/85">
                共 {total} 个模型
                {tokenGroup ? ` · 分组倍率 ${tokenGroup}` : ""}
              </Text>
              <Paragraph className="!mt-2 !mb-0 !text-sm !text-white/75">
                查看可用的 AI 模型与分组定价，支持按供应商、端点类型筛选。
              </Paragraph>
            </div>
            <Tag color="white" size="large" className="!rounded-full !font-medium">
              模型广场
            </Tag>
          </div>
        </div>
        <div className="mb-4 flex flex-wrap items-center gap-3">
          <Input
            prefix={<IconSearch />}
            placeholder="搜索模型名称"
            value={search}
            onChange={setSearch}
            showClear
            className="!w-64"
          />
          <Button icon={<IconCopy />} onClick={copyList}>
            复制列表
          </Button>
          <Switch checked={showPrice} onChange={setShowPrice} aria-label="显示价格" />
          <Text size="small">价格</Text>
          <Switch
            checked={showMultiplier}
            onChange={setShowMultiplier}
            aria-label="倍率"
          />
          <Text size="small">倍率</Text>
          <RadioGroup
            type="button"
            value={viewMode}
            onChange={(e) => setViewMode(e.target.value as ViewMode)}
          >
            <Radio value="grid">卡片</Radio>
            <Radio value="table">表格</Radio>
          </RadioGroup>
          <RadioGroup
            type="button"
            value={sizeMode}
            onChange={(e) => setSizeMode(e.target.value as SizeMode)}
          >
            <Radio value="M">M</Radio>
            <Radio value="L">L</Radio>
          </RadioGroup>
          <Text type="tertiary" className="ml-auto">
            共 {total} 个模型
            {showMultiplier ? ` · 分组 ${tokenGroup}` : ""}
          </Text>
        </div>

        {viewMode === "table" ? (
          <Table
            columns={columns}
            dataSource={items}
            loading={loading}
            pagination={{
              currentPage: page,
              pageSize: 24,
              total,
              onPageChange: setPage,
            }}
            rowKey="model"
          />
        ) : (
          <>
            <div className={gridClass}>
              {items.map((row) => {
                const p = priceFor(row);
                return (
                  <Card
                    key={row.model}
                    shadows="hover"
                    className="!rounded-xl"
                    bodyStyle={{ padding: sizeMode === "L" ? 20 : 14 }}
                    loading={loading}
                  >
                    <div className="mb-2 flex items-start justify-between gap-2">
                      <div>
                        <Title heading={sizeMode === "L" ? 5 : 6}>
                          {row.display_name}
                        </Title>
                        <Text type="tertiary" size="small">
                          {row.model}
                        </Text>
                      </div>
                      <Tag color="cyan">{row.billing_label}</Tag>
                    </div>
                    <div className="mb-2 flex flex-wrap gap-1">
                      <Tag size="small">{row.provider}</Tag>
                      <Tag size="small" color="grey">
                        {row.endpoint_type}
                      </Tag>
                    </div>
                    {showPrice && (
                      <div className="space-y-1 text-sm">
                        {row.billing_type === 2 ? (
                          <Text>
                            单价：{formatQuotaPerMillion(p.unit_price)} 额度
                          </Text>
                        ) : (
                          <>
                            <Text className="block">
                              输入 / 1M：{formatQuotaPerMillion(p.input_per_million)}
                            </Text>
                            <Text className="block">
                              补全 / 1M：
                              {formatQuotaPerMillion(p.output_per_million)}
                            </Text>
                            {p.cache_read_per_million > 0 && (
                              <Text type="tertiary" className="block">
                                缓存读 / 1M：
                                {formatQuotaPerMillion(p.cache_read_per_million)}
                              </Text>
                            )}
                          </>
                        )}
                      </div>
                    )}
                  </Card>
                );
              })}
            </div>
            {total > 24 && (
              <div className="mt-4 flex justify-center gap-2">
                <Button disabled={page <= 1} onClick={() => setPage((p) => p - 1)}>
                  上一页
                </Button>
                <Text className="self-center">
                  {page} / {Math.ceil(total / 24)}
                </Text>
                <Button
                  disabled={page >= Math.ceil(total / 24)}
                  onClick={() => setPage((p) => p + 1)}
                >
                  下一页
                </Button>
              </div>
            )}
          </>
        )}
      </div>
    </div>
  );
}
