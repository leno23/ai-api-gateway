"use client";

import { Button, Card, Typography } from "@douyinfe/semi-ui";
import { IconExternalOpen } from "@douyinfe/semi-icons";
import Link from "next/link";

import { usePortalLocale } from "@/contexts/portal-locale-context";
import { getGatewayBaseUrl } from "@/lib/gateway-config";
import { getDocsUrl } from "@/lib/portal-config";

const { Title, Paragraph, Text } = Typography;

export function DocsPageContent() {
  const { messages: m } = usePortalLocale();
  const externalDocs = getDocsUrl();
  const gatewayBase = getGatewayBaseUrl().replace(/\/$/, "");

  return (
    <div className="mx-auto max-w-3xl px-6 py-16">
      <div className="flex flex-wrap items-start justify-between gap-4">
        <Title heading={2}>{m.docs.title}</Title>
        {externalDocs ? (
          <a href={externalDocs} target="_blank" rel="noopener noreferrer">
            <Button icon={<IconExternalOpen />} theme="solid" type="primary">
              {m.docs.openExternal}
            </Button>
          </a>
        ) : null}
      </div>
      <Paragraph className="!mt-4">{m.docs.intro}</Paragraph>

      <Card className="!mt-8 !rounded-xl" title={m.docs.quickStart}>
        <div className="space-y-6">
          <div>
            <Text strong>{m.docs.authTitle}</Text>
            <Paragraph className="!mt-1 !mb-0">{m.docs.authBody}</Paragraph>
          </div>
          <div>
            <Text strong>{m.docs.baseUrlTitle}</Text>
            <Paragraph className="!mt-1 !mb-0">
              {m.docs.baseUrlBody}
              {gatewayBase ? (
                <>
                  <br />
                  <code className="mt-2 inline-block rounded bg-[var(--portal-bg)] px-2 py-1 text-sm">
                    {gatewayBase}/v1
                  </code>
                </>
              ) : null}
            </Paragraph>
          </div>
          <div>
            <Text strong>{m.docs.errorsTitle}</Text>
            <Paragraph className="!mt-1 !mb-0">{m.docs.errorsBody}</Paragraph>
          </div>
        </div>
      </Card>

      <Paragraph className="!mt-8">
        <Link href="/openapi.yaml" className="text-[var(--portal-primary)]">
          OpenAPI
        </Link>
        {gatewayBase ? (
          <>
            {" "}
            (
            <a
              href={`${gatewayBase}/openapi.yaml`}
              target="_blank"
              rel="noopener noreferrer"
              className="text-[var(--portal-primary)]"
            >
              {gatewayBase}/openapi.yaml
            </a>
            )
          </>
        ) : (
          "（配置 NEXT_PUBLIC_GATEWAY_API_URL 后显示网关地址）"
        )}
      </Paragraph>
    </div>
  );
}
