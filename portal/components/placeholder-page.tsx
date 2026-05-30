"use client";

import { Card, Typography } from "@douyinfe/semi-ui";

const { Title, Text } = Typography;

type PlaceholderPageProps = {
  title: string;
  phase: number;
};

export function PlaceholderPage({ title, phase }: PlaceholderPageProps) {
  return (
    <div>
      <Title heading={4}>{title}</Title>
      <Card className="!mt-4 !rounded-xl">
        <Text type="secondary">该页面将在 OpenSpec Phase {phase} 实现。</Text>
      </Card>
    </div>
  );
}
