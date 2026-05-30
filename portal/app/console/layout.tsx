"use client";

import { Layout } from "@douyinfe/semi-ui";

import { ConsoleGuard } from "@/components/console-guard";
import { ConsoleSidebar } from "@/components/console-sidebar";

const { Content } = Layout;

export default function ConsoleLayout({
  children,
}: {
  children: React.ReactNode;
}) {
  return (
    <ConsoleGuard>
      <Layout className="min-h-[calc(100vh-3.5rem)] bg-[var(--portal-bg)]">
        <ConsoleSidebar />
        <Content className="!p-6">{children}</Content>
      </Layout>
    </ConsoleGuard>
  );
}
