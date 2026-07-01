"use client";

import { Card, Typography } from "@douyinfe/semi-ui";
import Link from "next/link";

import { usePortalLocale } from "@/contexts/portal-locale-context";

const { Title, Paragraph, Text } = Typography;

export function AboutPageContent() {
  const { messages: m } = usePortalLocale();

  return (
    <div className="mx-auto max-w-3xl px-6 py-16">
      <Title heading={2}>{m.about.title}</Title>
      <Paragraph className="!mt-4 !text-base">{m.about.intro}</Paragraph>

      <Card className="!mt-8 !rounded-xl" title={m.about.missionTitle}>
        <Paragraph>{m.about.mission}</Paragraph>
      </Card>

      <Card className="!mt-6 !rounded-xl" title={m.about.featureTitle}>
        <ul className="list-disc space-y-2 pl-5">
          {m.about.features.map((item) => (
            <li key={item}>
              <Text>{item}</Text>
            </li>
          ))}
        </ul>
      </Card>

      <Card className="!mt-6 !rounded-xl" title={m.about.contactTitle}>
        <Paragraph>{m.about.contact}</Paragraph>
        <Paragraph className="!mt-4 !mb-0">
          <Link href="/docs" className="text-[var(--portal-primary)]">
            {m.nav.docs}
          </Link>
          {" · "}
          <Link href="/pricing" className="text-[var(--portal-primary)]">
            {m.nav.pricing}
          </Link>
        </Paragraph>
      </Card>
    </div>
  );
}
