"use client";

import {
  createContext,
  useCallback,
  useContext,
  useEffect,
  useMemo,
  useState,
} from "react";

import type { UserProfile } from "@/lib/api/types";
import { gatewayFetch } from "@/lib/api/client";
import { authStorage } from "@/lib/auth-storage";
import { isGatewayConfigured } from "@/lib/gateway-config";

type AuthContextValue = {
  token: string | null;
  user: UserProfile | null;
  ready: boolean;
  setToken: (t: string | null) => void;
  setUser: (u: UserProfile | null) => void;
  refreshUser: () => Promise<void>;
  logout: () => void;
};

const AuthContext = createContext<AuthContextValue | null>(null);

export function AuthProvider({ children }: { children: React.ReactNode }) {
  const [token, setTokenState] = useState<string | null>(null);
  const [user, setUser] = useState<UserProfile | null>(null);
  const [ready, setReady] = useState(false);

  const refreshUser = useCallback(async () => {
    const t = authStorage.getToken();
    if (!t || !isGatewayConfigured()) {
      setUser(null);
      return;
    }
    try {
      const res = await gatewayFetch<{ success: boolean; data: UserProfile }>(
        "/api/user/self",
        { token: t },
      );
      setUser(res.data);
    } catch {
      setUser(null);
    }
  }, []);

  useEffect(() => {
    let cancelled = false;
    const t = authStorage.getToken();
    queueMicrotask(async () => {
      if (cancelled) return;
      setTokenState(t);
      if (t) {
        await refreshUser();
      }
      if (!cancelled) setReady(true);
    });
    return () => {
      cancelled = true;
    };
  }, [refreshUser]);

  const setToken = useCallback(
    (next: string | null) => {
      if (next) {
        authStorage.setToken(next);
      } else {
        authStorage.clearToken();
        setUser(null);
      }
      setTokenState(next);
    },
    [],
  );

  const logout = useCallback(() => {
    authStorage.clearToken();
    setTokenState(null);
    setUser(null);
  }, []);

  const value = useMemo(
    () => ({
      token,
      user,
      ready,
      setToken,
      setUser,
      refreshUser,
      logout,
    }),
    [token, user, ready, setToken, refreshUser, logout],
  );

  return <AuthContext.Provider value={value}>{children}</AuthContext.Provider>;
}

export function useAuth(): AuthContextValue {
  const ctx = useContext(AuthContext);
  if (!ctx) {
    throw new Error("useAuth must be used within AuthProvider");
  }
  return ctx;
}
