"use client";

import { Button, Card, Tag, Typography } from "@douyinfe/semi-ui";
import { Play, FileText, Zap, Shield } from "lucide-react";
import Link from "next/link";

import { AnnouncementModal } from "@/components/announcement/announcement-modal";
import { PartnerMarquee } from "@/components/marketing/partner-marquee";
import { usePortalLocale } from "@/contexts/portal-locale-context";

const { Title, Paragraph, Text } = Typography;

export default function HomePage() {
  const { messages: m } = usePortalLocale();

  const features = [
    {
      icon: Zap,
      title: m.home.feature1Title,
      desc: "基于协程池与连接复用，从容应对高并发 Token 请求，毫秒级负载均衡，为 AI 生产力提供稳定心跳。",
    },
    {
      icon: Shield,
      title: m.home.feature2Title,
      desc: "打破算力垄断，提供透明分组倍率与按量计费，让每一次模型调用都物超所值。",
    },
    {
      icon: FileText,
      title: m.home.feature3Title,
      desc: "OpenAI / Anthropic 兼容端点，一套密钥接入 Claude、GPT 等主流模型，降低集成成本。",
    },
  ];

  return (
    <div className="bg-white">
      <AnnouncementModal placement="home" />
      <section className="mx-auto max-w-5xl px-6 py-16 text-center md:py-24">
        <Tag
          color="blue"
          size="large"
          className="!mb-6 !rounded-full !px-4 !py-1"
        >
          Claude / Codex 核心接口已接入
        </Tag>
        <Title
          heading={1}
          className="portal-hero-gradient !text-4xl !font-bold !leading-tight md:!text-5xl lg:!text-6xl"
        >
          {m.home.heroTitle}
        </Title>
        <Title heading={2} className="!mt-3 !text-2xl !font-bold !text-[var(--semi-color-text-0)] md:!text-3xl">
          {m.home.heroSubtitle}
        </Title>
        <Paragraph className="!mx-auto !mt-6 !max-w-2xl !text-base !leading-relaxed !text-[var(--semi-color-text-2)]">
          专为 Claude、OpenAI Codex 等头部大模型提供高并发、高可用 API 接入。极致的稳定性保证，全网更具性价比的算力调用引擎。
        </Paragraph>
        <div className="mt-10 flex flex-wrap justify-center gap-4">
          <Link href="/register">
            <Button
              theme="solid"
              type="primary"
              size="large"
              icon={<Play className="h-4 w-4 fill-current" />}
              className="!rounded-full !px-8 !shadow-md"
            >
              立即开始使用
            </Button>
          </Link>
          <Link href="/pricing">
            <Button
              theme="light"
              size="large"
              icon={<FileText className="h-4 w-4" />}
              className="!rounded-full !border !border-[var(--semi-color-border)] !px-8"
            >
              探索模型
            </Button>
          </Link>
        </div>
      </section>

      <section className="bg-[var(--portal-bg)] py-16 md:py-20">
        <div className="mx-auto grid max-w-5xl gap-6 px-6 md:grid-cols-2 lg:grid-cols-3">
          {features.map((f) => (
            <Card
              key={f.title}
              shadows="hover"
              className="portal-card !border-0 !shadow-sm"
            >
              <span className="mb-4 inline-flex rounded-xl bg-white p-2 shadow-sm ring-1 ring-black/5">
                <f.icon className="h-6 w-6 text-[var(--portal-primary)]" />
              </span>
              <Title heading={4} className="!mb-2">
                {f.title}
              </Title>
              <Text type="secondary" className="!leading-relaxed">
                {f.desc}
              </Text>
            </Card>
          ))}
        </div>
      </section>

      <PartnerMarquee />
    </div>
  );
}
