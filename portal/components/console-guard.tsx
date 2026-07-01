"use client";

import { Toast } from "@douyinfe/semi-ui";
import { useRouter } from "next/navigation";
import { useEffect } from "react";

import { useAuth } from "@/contexts/auth-context";

export function ConsoleGuard({ children }: { children: React.ReactNode }) {
  const { token, ready } = useAuth();
  const router = useRouter();

  useEffect(() => {
    if (!ready) return;
    if (!token) {
      Toast.warning("未登录或登录已过期");
      router.replace("/login?expired=true");
    }
  }, [ready, token, router]);

  if (!ready) {
    return (
      <div className="flex flex-1 items-center justify-center p-12 text-[var(--semi-color-text-2)]">
        加载中…
      </div>
    );
  }

  if (!token) {
    return null;
  }

  return <>{children}</>;
}
