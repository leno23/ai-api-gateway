import type { Metadata } from "next";

import { AppProviders } from "@/app/providers";
import { ConfigGate } from "@/components/config-gate";
import { SiteHeader } from "@/components/site-header";

import "./globals.css";

export const metadata: Metadata = {
  title: "AI API Gateway",
  description: "企业级多模型接入网关",
};

export default function RootLayout({
  children,
}: Readonly<{
  children: React.ReactNode;
}>) {
  return (
    <html lang="zh-CN" className="h-full antialiased">
      <body className="flex min-h-screen flex-col">
        <AppProviders>
          <ConfigGate>
            <SiteHeader />
            <main className="flex flex-1 flex-col">{children}</main>
          </ConfigGate>
        </AppProviders>
      </body>
    </html>
  );
}
