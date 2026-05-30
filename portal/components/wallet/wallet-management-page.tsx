"use client";

import {
  Banner,
  Button,
  Card,
  Input,
  InputNumber,
  Toast,
  Typography,
} from "@douyinfe/semi-ui";
import { IconCopy } from "@douyinfe/semi-icons";
import { useCallback, useState } from "react";

import type { WalletSummaryResponse } from "@/lib/api/wallet-types";
import { ApiError } from "@/lib/api/client";
import { useDeferredEffect } from "@/lib/hooks/use-deferred-effect";
import { useGatewayRequest } from "@/lib/hooks/use-gateway-request";

const { Title, Text } = Typography;

export function WalletManagementPage() {
  const request = useGatewayRequest();
  const [data, setData] = useState<WalletSummaryResponse["data"] | null>(null);
  const [loading, setLoading] = useState(true);
  const [redeemCode, setRedeemCode] = useState("");
  const [rechargeAmount, setRechargeAmount] = useState(10000);

  const load = useCallback(async () => {
    setLoading(true);
    try {
      const res = await request<WalletSummaryResponse>("/api/wallet/summary");
      setData(res.data);
    } catch (e) {
      Toast.error(e instanceof ApiError ? e.message : "加载失败");
    } finally {
      setLoading(false);
    }
  }, [request]);

  useDeferredEffect(() => load(), [load]);

  const copyInvite = async () => {
    if (!data?.invite_url) return;
    await navigator.clipboard.writeText(data.invite_url);
    Toast.success("已复制邀请链接");
  };

  return (
    <div className="space-y-6">
      <Title heading={4}>钱包管理</Title>

      <div className="grid gap-4 md:grid-cols-3">
        <Card loading={loading} className="!rounded-xl" shadows="hover">
          <Text type="tertiary">当前余额（额度）</Text>
          <Title heading={3} className="!mt-2">
            {data?.quota ?? "—"}
          </Title>
        </Card>
        <Card loading={loading} className="!rounded-xl" shadows="hover">
          <Text type="tertiary">历史消耗</Text>
          <Title heading={3} className="!mt-2">
            {data?.used_quota ?? "—"}
          </Title>
        </Card>
        <Card loading={loading} className="!rounded-xl" shadows="hover">
          <Text type="tertiary">累计请求数</Text>
          <Title heading={3} className="!mt-2">
            {data?.request_count ?? "—"}
          </Title>
        </Card>
      </div>

      <Card title="兑换码充值" className="!rounded-xl">
        <div className="flex flex-wrap gap-2">
          <Input
            placeholder="输入兑换码"
            value={redeemCode}
            onChange={setRedeemCode}
            className="!w-64"
          />
          <Button
            type="primary"
            theme="solid"
            onClick={async () => {
              if (!redeemCode.trim()) {
                Toast.warning("请输入兑换码");
                return;
              }
              try {
                await request("/api/wallet/redeem", {
                  method: "POST",
                  body: JSON.stringify({ code: redeemCode.trim() }),
                });
                Toast.success("兑换成功");
                setRedeemCode("");
                void load();
              } catch (e) {
                Toast.error(e instanceof ApiError ? e.message : "兑换失败");
              }
            }}
          >
            兑换
          </Button>
        </div>
      </Card>

      <Card title="邀请返利" className="!rounded-xl">
        <Text type="secondary" className="mb-3 block">
          好友通过邀请链接注册并充值后，你将获得待结算收益，可划转到余额。
        </Text>
        <div className="mb-3 flex flex-wrap items-center gap-2">
          <Input readonly value={data?.invite_url ?? ""} className="!min-w-[280px] !flex-1" />
          <Button icon={<IconCopy />} onClick={() => void copyInvite()}>
            复制链接
          </Button>
        </div>
        <Text>
          待使用收益：<strong>{data?.affiliate_pending ?? 0}</strong> 额度
        </Text>
        <Button
          className="!mt-3"
          type="primary"
          theme="solid"
          disabled={!data?.affiliate_pending}
          onClick={async () => {
            try {
              await request("/api/wallet/affiliate/transfer", { method: "POST" });
              Toast.success("已划转到余额");
              void load();
            } catch (e) {
              Toast.error(e instanceof ApiError ? e.message : "划转失败");
            }
          }}
        >
          划转到余额
        </Button>
      </Card>

      <Card title="在线充值" className="!rounded-xl">
        {data?.recharge_enabled ? (
          <div className="flex flex-wrap items-end gap-3">
            <div>
              <Text type="tertiary" size="small">
                充值额度（mock）
              </Text>
              <InputNumber
                value={rechargeAmount}
                onChange={(v) => setRechargeAmount(Number(v) || 0)}
                min={1}
                className="!mt-1 !w-40"
              />
            </div>
            <Button
              type="primary"
              theme="solid"
              onClick={async () => {
                try {
                  await request("/api/wallet/recharge/mock", {
                    method: "POST",
                    body: JSON.stringify({ amount: rechargeAmount }),
                  });
                  Toast.success("充值成功（mock）");
                  void load();
                } catch (e) {
                  Toast.error(e instanceof ApiError ? e.message : "充值失败");
                }
              }}
            >
              立即充值
            </Button>
          </div>
        ) : (
          <Banner
            type="info"
            description="在线支付未开启，请联系管理员或使用兑换码充值。"
            closeIcon={null}
          />
        )}
      </Card>
    </div>
  );
}
