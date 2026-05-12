"use client";

import { useRouter } from "next/navigation";
import { useCallback } from "react";

import { ApiError, gatewayFetch } from "@/lib/api/client";
import { useAuth } from "@/contexts/auth-context";

/**
 * 带会话处理的网关请求：401 清会话并回登录；403 跳转无权限页。
 */
export function useGatewayRequest() {
  const { token, logout } = useAuth();
  const router = useRouter();

  return useCallback(
    async <T,>(
      path: string,
      init: RequestInit & { token?: string | null } = {},
    ): Promise<T> => {
      try {
        return await gatewayFetch<T>(path, { ...init, token: init.token ?? token });
      } catch (e) {
        if (e instanceof ApiError) {
          if (e.status === 401) {
            logout();
            router.replace("/login");
          } else if (e.status === 403) {
            router.replace("/forbidden");
          }
        }
        throw e;
      }
    },
    [token, logout, router],
  );
}
