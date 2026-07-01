"use client";

import { isGatewayConfigured } from "@/lib/gateway-config";

import { ConfigError } from "./config-error";

export function ConfigGate({ children }: { children: React.ReactNode }) {
  if (!isGatewayConfigured()) {
    return <ConfigError />;
  }
  return <>{children}</>;
}
