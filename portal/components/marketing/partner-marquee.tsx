"use client";

import { Typography } from "@douyinfe/semi-ui";

import { usePortalLocale } from "@/contexts/portal-locale-context";

const { Text } = Typography;

const PARTNERS = [
  "OpenAI",
  "Anthropic",
  "Google",
  "DeepSeek",
  "Meta",
  "Mistral",
  "Qwen",
  "Moonshot",
  "Zhipu",
  "MiniMax",
];

export function PartnerMarquee() {
  const { messages } = usePortalLocale();
  const row = [...PARTNERS, ...PARTNERS];

  return (
    <section className="overflow-hidden border-y border-[var(--semi-color-border)] bg-white py-10">
      <Text type="tertiary" className="mb-6 block text-center">
        {messages.home.partnersTitle}
      </Text>
      <div className="relative flex">
        <div className="animate-partner-marquee flex shrink-0 gap-12 px-6">
          {row.map((name, i) => (
            <span
              key={`${name}-${i}`}
              className="shrink-0 text-lg font-semibold tracking-wide text-[var(--semi-color-text-2)]"
            >
              {name}
            </span>
          ))}
        </div>
      </div>
    </section>
  );
}
