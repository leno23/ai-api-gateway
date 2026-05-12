"use client";

import { useRouter } from "next/navigation";
import { useEffect } from "react";

import { useAuth } from "@/contexts/auth-context";

export default function HomePage() {
  const { token, ready } = useAuth();
  const router = useRouter();

  useEffect(() => {
    if (!ready) return;
    router.replace(token ? "/admin" : "/login");
  }, [token, ready, router]);

  return (
    <div className="flex flex-1 items-center justify-center p-8 text-neutral-500">
      正在跳转…
    </div>
  );
}
