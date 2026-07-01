"use client";

import { Banner } from "@douyinfe/semi-ui";

export function ConfigError() {
  return (
    <div className="flex min-h-screen items-center justify-center p-6">
      <Banner
        type="danger"
        fullMode={false}
        bordered
        closeIcon={null}
        title="未配置网关地址"
        description={
          <>
            请在 <code className="rounded bg-black/5 px-1">portal/.env.local</code>{" "}
            中设置 <code className="rounded bg-black/5 px-1">NEXT_PUBLIC_GATEWAY_API_URL</code>
            （无尾斜杠），例如 <code className="rounded bg-black/5 px-1">http://localhost:8080</code>
          </>
        }
      />
    </div>
  );
}
