"use client";

import { useCallback } from "react";

import { gatewayFetch } from "@/lib/api/client";
import { useAuth } from "@/contexts/auth-context";

export function useGatewayRequest() {
  const { token } = useAuth();

  return useCallback(
    <T,>(path: string, init: RequestInit = {}) =>
      gatewayFetch<T>(path, { ...init, token }),
    [token],
  );
}
