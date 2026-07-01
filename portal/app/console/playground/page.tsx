"use client";

import { Suspense } from "react";

import { PlaygroundPage } from "@/components/playground/playground-page";

export default function ConsolePlaygroundRoute() {
  return (
    <Suspense fallback={<div className="p-6">加载中…</div>}>
      <PlaygroundPage />
    </Suspense>
  );
}
