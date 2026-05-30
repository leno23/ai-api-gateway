"use client";

import { Button, Modal, Typography } from "@douyinfe/semi-ui";
import { useCallback, useState } from "react";

import type { AnnouncementItem } from "@/lib/api/wallet-types";
import { gatewayFetch } from "@/lib/api/client";
import { useDeferredEffect } from "@/lib/hooks/use-deferred-effect";

const { Paragraph, Title } = Typography;

const DISMISS_PREFIX = "portal_announcement_dismiss_";

function todayKey() {
  return new Date().toISOString().slice(0, 10);
}

type AnnouncementModalProps = {
  placement?: string;
};

export function AnnouncementModal({ placement = "home" }: AnnouncementModalProps) {
  const [visible, setVisible] = useState(false);
  const [item, setItem] = useState<AnnouncementItem | null>(null);

  const load = useCallback(async () => {
    try {
      const res = await gatewayFetch<{ success: boolean; data: AnnouncementItem[] }>(
        `/api/announcements?placement=${placement}`,
      );
      const first = res.data[0];
      if (!first) return;
      const key = `${DISMISS_PREFIX}${first.id}_${todayKey()}`;
      if (localStorage.getItem(key) === "1") return;
      setItem(first);
      setVisible(true);
    } catch {
      // ignore when gateway offline
    }
  }, [placement]);

  useDeferredEffect(() => load(), [load]);

  if (!item) return null;

  return (
    <Modal
      title={item.title}
      visible={visible}
      onCancel={() => setVisible(false)}
      footer={
        <div className="flex w-full justify-end gap-2">
          <Button
            onClick={() => {
              localStorage.setItem(`${DISMISS_PREFIX}${item.id}_${todayKey()}`, "1");
              setVisible(false);
            }}
          >
            今日不再提示
          </Button>
          <Button type="primary" theme="solid" onClick={() => setVisible(false)}>
            知道了
          </Button>
        </div>
      }
    >
      <Paragraph>{item.content}</Paragraph>
      {item.level === "warning" && (
        <Title heading={6} type="warning">
          请关注系统通知
        </Title>
      )}
    </Modal>
  );
}
