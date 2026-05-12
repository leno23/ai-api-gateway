"use client";

import { useRouter } from "next/navigation";
import { useEffect } from "react";

import { AdminShell } from "@/components/admin-shell";
import { useAuth } from "@/contexts/auth-context";

export default function AdminLayout({
  children,
}: {
  children: React.ReactNode;
}) {
  const { token, ready, logout } = useAuth();
  const router = useRouter();

  useEffect(() => {
    if (!ready) return;
    if (!token) {
      router.replace("/login");
    }
  }, [token, ready, router]);

  if (!ready || !token) {
    return (
      <div className="flex min-h-screen items-center justify-center text-neutral-400">
        加载中…
      </div>
    );
  }

  return (
    <AdminShell
      onLogout={() => {
        logout();
        router.replace("/login");
      }}
    >
      {children}
    </AdminShell>
  );
}
