"use client";

import {
  createContext,
  useCallback,
  useContext,
  useEffect,
  useMemo,
  useState,
} from "react";

import { authStorage } from "@/lib/auth-storage";

type AuthContextValue = {
  token: string | null;
  ready: boolean;
  setToken: (t: string | null) => void;
  logout: () => void;
};

const AuthContext = createContext<AuthContextValue | null>(null);

export function AuthProvider({ children }: { children: React.ReactNode }) {
  const [token, setTokenState] = useState<string | null>(null);
  const [ready, setReady] = useState(false);

  useEffect(() => {
    let cancelled = false;
    queueMicrotask(() => {
      if (!cancelled) {
        setTokenState(authStorage.getToken());
        setReady(true);
      }
    });
    return () => {
      cancelled = true;
    };
  }, []);

  const setToken = useCallback((t: string | null) => {
    if (t) {
      authStorage.setToken(t);
    } else {
      authStorage.clearToken();
    }
    setTokenState(t);
  }, []);

  const logout = useCallback(() => {
    authStorage.clearToken();
    setTokenState(null);
  }, []);

  const value = useMemo(
    () => ({ token, ready, setToken, logout }),
    [token, ready, setToken, logout],
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
