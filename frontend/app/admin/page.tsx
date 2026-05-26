"use client";

import { Typography } from "antd";
import Link from "next/link";

import { getGatewayBaseUrl } from "@/lib/gateway-config";

export default function AdminHomePage() {
  const base = getGatewayBaseUrl().replace(/\/$/, "");

  return (
    <div className="flex max-w-2xl flex-col gap-6">
      <div>
        <Typography.Title level={4} className="!mb-2">
          概览
        </Typography.Title>
        <Typography.Paragraph type="secondary" className="!mb-0">
          使用左侧菜单管理渠道、兑换码与用户状态。契约与网关{" "}
          <Link href={`${base}/openapi.yaml`} target="_blank" rel="noreferrer">
            OpenAPI
          </Link>{" "}
          对齐。
        </Typography.Paragraph>
      </div>
      <ul className="list-inside list-disc text-neutral-600">
        <li>
          健康检查：<code className="rounded bg-neutral-100 px-1">{base}/health</code>
        </li>
        <li>
          指标（只读）：{" "}
          <Link href={`${base}/metrics`} target="_blank" rel="noreferrer">
            /metrics
          </Link>
        </li>
      </ul>
    </div>
  );
}
