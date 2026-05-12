import { getGatewayBaseUrl } from "@/lib/gateway-config";

export class ApiError extends Error {
  status: number;
  body: unknown;

  constructor(message: string, status: number, body?: unknown) {
    super(message);
    this.name = "ApiError";
    this.status = status;
    this.body = body;
  }
}

function parseErrorMessage(data: unknown, fallback: string): string {
  if (typeof data === "object" && data !== null && "error" in data) {
    return String((data as { error: unknown }).error);
  }
  return fallback;
}

export async function gatewayFetch<T>(
  path: string,
  init: RequestInit & { token?: string | null } = {},
): Promise<T> {
  const base = getGatewayBaseUrl().replace(/\/$/, "");
  const { token, ...rest } = init;
  const headers = new Headers(rest.headers);
  if (!headers.has("Content-Type")) {
    headers.set("Content-Type", "application/json");
  }
  if (token) {
    headers.set("Authorization", `Bearer ${token}`);
  }

  const res = await fetch(`${base}${path}`, { ...rest, headers });
  const text = await res.text();
  let data: unknown = null;
  if (text) {
    try {
      data = JSON.parse(text) as unknown;
    } catch {
      data = text;
    }
  }

  if (!res.ok) {
    throw new ApiError(
      parseErrorMessage(data, res.statusText),
      res.status,
      data,
    );
  }
  return data as T;
}
