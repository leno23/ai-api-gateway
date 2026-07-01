"use client";

import { Card, Typography } from "@douyinfe/semi-ui";

import { useAuth } from "@/contexts/auth-context";

const { Title, Text } = Typography;

export default function SettingPage() {
  const { user } = useAuth();

  return (
    <div className="space-y-4">
      <Title heading={4}>个人设置</Title>
      <Card title="账户管理" className="!rounded-xl">
        <Text>用户名：{user?.username ?? "—"}</Text>
        <br />
        <Text>邮箱：{user?.email ?? "—"}</Text>
        <br />
        <Text>邀请码：{user?.invite_code ?? "—"}</Text>
      </Card>
      <Card title="其他设置" className="!rounded-xl">
        <Text type="secondary">通知方式、余额预警等将在 Phase 5 接入。</Text>
      </Card>
    </div>
  );
}
