"use client";

import { AntdRegistry } from "@ant-design/nextjs-registry";
import { Alert, App, ConfigProvider, theme } from "antd";
import zhCN from "antd/locale/zh_CN";

import { AuthProvider } from "@/contexts/auth-context";
import { isGatewayConfigured } from "@/lib/gateway-config";

function GatewayConfigGate({ children }: { children: React.ReactNode }) {
  if (!isGatewayConfigured()) {
    return (
      <div className="mx-auto mt-20 max-w-xl p-8">
        <Alert
          type="error"
          showIcon
          message="缺少网关地址配置"
          description={
            <span>
              请在 <code className="rounded bg-neutral-100 px-1">frontend/.env.local</code>{" "}
              中设置{" "}
              <code className="rounded bg-neutral-100 px-1">
                NEXT_PUBLIC_GATEWAY_API_URL
              </code>
              （无尾斜杠），例如{" "}
              <code className="rounded bg-neutral-100 px-1">
                http://127.0.0.1:8080
              </code>
              。可参考 <code className="rounded bg-neutral-100 px-1">.env.example</code>
              。
            </span>
          }
        />
      </div>
    );
  }
  return <>{children}</>;
}

export function AppProviders({ children }: { children: React.ReactNode }) {
  return (
    <AntdRegistry>
      <ConfigProvider
        locale={zhCN}
        theme={{
          algorithm: theme.defaultAlgorithm,
          token: { borderRadius: 6 },
        }}
      >
        <App>
          <AuthProvider>
            <GatewayConfigGate>{children}</GatewayConfigGate>
          </AuthProvider>
        </App>
      </ConfigProvider>
    </AntdRegistry>
  );
}
