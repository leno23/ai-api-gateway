"use client";

import { Button, Tag, Typography } from "@douyinfe/semi-ui";

const { Title, Text } = Typography;

type FilterTagProps = {
  active: boolean;
  label: string;
  onClick: () => void;
};

function FilterTag({ active, label, onClick }: FilterTagProps) {
  return (
    <Tag
      size="large"
      color={active ? "purple" : "white"}
      className={`!cursor-pointer !rounded-full !px-3 !py-1 ${
        active ? "" : "!border !border-[var(--semi-color-border)]"
      }`}
      onClick={onClick}
    >
      {label}
    </Tag>
  );
}

type FilterSectionProps = {
  title: string;
  children: React.ReactNode;
};

function FilterSection({ title, children }: FilterSectionProps) {
  return (
    <div>
      <Text strong className="mb-2 block text-sm">
        {title}
      </Text>
      <div className="flex flex-wrap gap-2">{children}</div>
    </div>
  );
}

export type PricingFilterSidebarProps = {
  total: number;
  providers: string[];
  tokenGroups: { slug: string; name: string; multiplier: number }[];
  provider: string;
  tokenGroup: string;
  endpointType: string;
  billingType: string;
  tag: string;
  onProviderChange: (v: string) => void;
  onTokenGroupChange: (v: string) => void;
  onEndpointTypeChange: (v: string) => void;
  onBillingTypeChange: (v: string) => void;
  onTagChange: (v: string) => void;
  onReset: () => void;
};

export function PricingFilterSidebar({
  total,
  providers,
  tokenGroups,
  provider,
  tokenGroup,
  endpointType,
  billingType,
  tag,
  onProviderChange,
  onTokenGroupChange,
  onEndpointTypeChange,
  onBillingTypeChange,
  onTagChange,
  onReset,
}: PricingFilterSidebarProps) {
  return (
    <aside className="hidden w-64 shrink-0 md:block">
      <div className="sticky top-20 rounded-xl border border-[var(--semi-color-border)] bg-white p-4">
        <div className="mb-4 flex items-center justify-between">
          <Title heading={6}>筛选</Title>
          <Button theme="borderless" size="small" onClick={onReset}>
            重置
          </Button>
        </div>
        <div className="space-y-5">
          <FilterSection title="供应商">
            <FilterTag
              active={!provider}
              label={`全部供应商 ${total}`}
              onClick={() => onProviderChange("")}
            />
            {providers.map((p) => (
              <FilterTag
                key={p}
                active={provider === p}
                label={p}
                onClick={() => onProviderChange(p)}
              />
            ))}
          </FilterSection>

          <FilterSection title="可用令牌分组">
            {tokenGroups.map((g) => (
              <FilterTag
                key={g.slug}
                active={tokenGroup === g.slug}
                label={`${g.name} ${g.multiplier}x`}
                onClick={() => onTokenGroupChange(g.slug)}
              />
            ))}
          </FilterSection>

          <FilterSection title="计费类型">
            <FilterTag
              active={!billingType}
              label="全部类型"
              onClick={() => onBillingTypeChange("")}
            />
            <FilterTag
              active={billingType === "1"}
              label="按量计费"
              onClick={() => onBillingTypeChange("1")}
            />
            <FilterTag
              active={billingType === "2"}
              label="按次计费"
              onClick={() => onBillingTypeChange("2")}
            />
          </FilterSection>

          <FilterSection title="端点类型">
            <FilterTag
              active={!endpointType}
              label="全部端点"
              onClick={() => onEndpointTypeChange("")}
            />
            <FilterTag
              active={endpointType === "openai"}
              label="openai"
              onClick={() => onEndpointTypeChange("openai")}
            />
            <FilterTag
              active={endpointType === "anthropic"}
              label="anthropic"
              onClick={() => onEndpointTypeChange("anthropic")}
            />
          </FilterSection>

          <FilterSection title="标签">
            <FilterTag active={!tag} label="全部标签" onClick={() => onTagChange("")} />
            <FilterTag active={tag === "chat"} label="chat" onClick={() => onTagChange("chat")} />
            <FilterTag
              active={tag === "image"}
              label="image"
              onClick={() => onTagChange("image")}
            />
          </FilterSection>
        </div>
      </div>
    </aside>
  );
}
