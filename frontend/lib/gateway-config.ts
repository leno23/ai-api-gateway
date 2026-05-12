/**
 * 网关基址（build-time public env）。未配置时开发环境应给出明确提示（见 console-shell-auth spec）。
 */
export function getGatewayBaseUrl(): string {
  const v = process.env.NEXT_PUBLIC_GATEWAY_API_URL?.trim();
  return v ?? "";
}

export function isGatewayConfigured(): boolean {
  return Boolean(getGatewayBaseUrl());
}
