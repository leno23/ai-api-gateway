"use client";

import { Card, Tag, Typography } from "@douyinfe/semi-ui";
import { useCallback, useState } from "react";

import type { AnnouncementItem } from "@/lib/api/wallet-types";
import { ApiError } from "@/lib/api/client";
import { useDeferredEffect } from "@/lib/hooks/use-deferred-effect";
import { useGatewayRequest } from "@/lib/hooks/use-gateway-request";

const { Text } = Typography;

const LEVEL_META: Record<
  string,
  { label: string; color: "grey" | "blue" | "green" | "orange" | "red" }
> = {
  default: { label: "默认", color: "grey" },
  info: { label: "进行中", color: "blue" },
  success: { label: "成功", color: "green" },
  warning: { label: "警告", color: "orange" },
  error: { label: "错误", color: "red" },
};

function levelMeta(level: string) {
  return LEVEL_META[level] ?? LEVEL_META.default;
}

type DashboardAnnouncementsCardProps = {
  placement?: string;
};

export function DashboardAnnouncementsCard({
  placement = "console",
}: DashboardAnnouncementsCardProps) {
  const request = useGatewayRequest();
  const [items, setItems] = useState<AnnouncementItem[]>([]);
  const [loading, setLoading] = useState(true);

  const load = useCallback(async () => {
    setLoading(true);
    try {
      const res = await request<{ success: boolean; data: AnnouncementItem[] }>(
        `/api/announcements?placement=${placement}`,
      );
      let list = res.data ?? [];
      if (list.length === 0 && placement !== "home") {
        const fallback = await request<{ success: boolean; data: AnnouncementItem[] }>(
          "/api/announcements?placement=home",
        );
        list = fallback.data ?? [];
      }
      setItems(list.slice(0, 6));
    } catch (e) {
      if (!(e instanceof ApiError)) {
        setItems([]);
      }
    } finally {
      setLoading(false);
    }
  }, [request, placement]);

  useDeferredEffect(() => load(), [load]);

  return (
    <Card title="系统公告" className="!rounded-xl" loading={loading}>
      <div className="mb-3 flex flex-wrap gap-2">
        {Object.entries(LEVEL_META).map(([key, meta]) => (
          <Tag key={key} size="small" color={meta.color}>
            {meta.label}
          </Tag>
        ))}
      </div>
      {items.length === 0 ? (
        <Text type="tertiary" size="small">
          暂无公告
        </Text>
      ) : (
        <ul className="space-y-3">
          {items.map((item) => {
            const meta = levelMeta(item.level);
            return (
              <li
                key={item.id}
                className="border-b border-[var(--semi-color-border)] pb-3 last:border-0 last:pb-0"
              >
                <div className="mb-1 flex items-center gap-2">
                  <Tag size="small" color={meta.color}>
                    {meta.label}
                  </Tag>
                  <Text strong size="small">
                    {item.title}
                  </Text>
                </div>
                <Text type="tertiary" size="small" className="line-clamp-2">
                  {item.content}
                </Text>
              </li>
            );
          })}
        </ul>
      )}
    </Card>
  );
}
