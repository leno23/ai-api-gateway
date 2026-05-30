"use client";

import { AuthProvider } from "@/contexts/auth-context";
import { PortalLocaleProvider } from "@/contexts/portal-locale-context";

export function AppProviders({ children }: { children: React.ReactNode }) {
  return (
    <PortalLocaleProvider>
      <AuthProvider>{children}</AuthProvider>
    </PortalLocaleProvider>
  );
}
