"use client";

import { LocaleProvider as SemiLocaleProvider } from "@douyinfe/semi-ui";
import en_US from "@douyinfe/semi-ui/lib/es/locale/source/en_US";
import zh_CN from "@douyinfe/semi-ui/lib/es/locale/source/zh_CN";
import {
  createContext,
  useCallback,
  useContext,
  useEffect,
  useMemo,
  useState,
} from "react";

import { useDeferredEffect } from "@/lib/hooks/use-deferred-effect";
import { t, type MessageTree, type PortalLocale } from "@/lib/i18n/messages";

const STORAGE_KEY = "portal_locale";

type PortalLocaleContextValue = {
  locale: PortalLocale;
  setLocale: (locale: PortalLocale) => void;
  messages: MessageTree;
};

const PortalLocaleContext = createContext<PortalLocaleContextValue | null>(null);

function readStoredLocale(): PortalLocale {
  if (typeof window === "undefined") return "zh-CN";
  const raw = localStorage.getItem(STORAGE_KEY);
  return raw === "en" ? "en" : "zh-CN";
}

export function PortalLocaleProvider({ children }: { children: React.ReactNode }) {
  const [locale, setLocaleState] = useState<PortalLocale>("zh-CN");
  const [ready, setReady] = useState(false);

  useDeferredEffect(() => {
    setLocaleState(readStoredLocale());
    setReady(true);
  }, []);

  const setLocale = useCallback((next: PortalLocale) => {
    setLocaleState(next);
    localStorage.setItem(STORAGE_KEY, next);
  }, []);

  useEffect(() => {
    if (!ready) return;
    document.documentElement.lang = locale === "en" ? "en" : "zh-CN";
  }, [locale, ready]);

  const value = useMemo(
    () => ({
      locale,
      setLocale,
      messages: t(locale),
    }),
    [locale, setLocale],
  );

  const semiLocale = locale === "en" ? en_US : zh_CN;

  if (!ready) {
    return (
      <SemiLocaleProvider locale={zh_CN}>
        <PortalLocaleContext.Provider value={value}>
          {children}
        </PortalLocaleContext.Provider>
      </SemiLocaleProvider>
    );
  }

  return (
    <SemiLocaleProvider locale={semiLocale}>
      <PortalLocaleContext.Provider value={value}>{children}</PortalLocaleContext.Provider>
    </SemiLocaleProvider>
  );
}

export function usePortalLocale(): PortalLocaleContextValue {
  const ctx = useContext(PortalLocaleContext);
  if (!ctx) {
    throw new Error("usePortalLocale must be used within PortalLocaleProvider");
  }
  return ctx;
}
