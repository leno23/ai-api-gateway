export function getGatewayBaseUrl(): string {
  return process.env.NEXT_PUBLIC_GATEWAY_API_URL ?? "";
}

export function isGatewayConfigured(): boolean {
  return getGatewayBaseUrl().trim().length > 0;
}
