"use client";

import { App, Button, Card, Form, Input, Typography } from "antd";
import Link from "next/link";
import { useRouter } from "next/navigation";
import { useState } from "react";

import { useAuth } from "@/contexts/auth-context";
import { ApiError, gatewayFetch } from "@/lib/api/client";
import type { LoginResponse } from "@/lib/api/types";
import { isGatewayConfigured } from "@/lib/gateway-config";

export default function LoginPage() {
  const { setToken } = useAuth();
  const router = useRouter();
  const { message } = App.useApp();
  const [loading, setLoading] = useState(false);

  if (!isGatewayConfigured()) {
    return null;
  }

  return (
    <div className="flex min-h-full flex-1 items-center justify-center bg-neutral-50 p-6">
      <Card className="w-full max-w-md shadow-sm" styles={{ body: { padding: 28 } }}>
        <Typography.Title level={3} className="!mb-6 text-center">
          管理员登录
        </Typography.Title>
        <Form
          layout="vertical"
          requiredMark={false}
          onFinish={async (v: { email: string; password: string }) => {
            setLoading(true);
            try {
              const data = await gatewayFetch<LoginResponse>("/auth/login", {
                method: "POST",
                body: JSON.stringify({
                  email: v.email,
                  password: v.password,
                }),
              });
              setToken(data.access_token);
              message.success("登录成功");
              router.replace("/admin");
            } catch (e) {
              if (e instanceof ApiError && e.status === 401) {
                message.error("邮箱或密码不正确");
              } else if (e instanceof ApiError) {
                message.error(e.message || "登录失败");
              } else {
                message.error("网络异常，请稍后重试");
              }
            } finally {
              setLoading(false);
            }
          }}
        >
          <Form.Item
            label="邮箱"
            name="email"
            rules={[{ required: true, type: "email", message: "请输入有效邮箱" }]}
          >
            <Input autoComplete="email" placeholder="admin@example.com" />
          </Form.Item>
          <Form.Item
            label="密码"
            name="password"
            rules={[{ required: true, message: "请输入密码" }]}
          >
            <Input.Password autoComplete="current-password" />
          </Form.Item>
          <Form.Item className="mb-0">
            <Button type="primary" htmlType="submit" loading={loading} block size="large">
              登录
            </Button>
          </Form.Item>
        </Form>
        <Typography.Paragraph type="secondary" className="!mb-0 !mt-6 text-center text-sm">
          管理员账号需在数据库中授予 role；普通用户无法访问管理接口。{" "}
          <Link href="/">返回首页</Link>
        </Typography.Paragraph>
      </Card>
    </div>
  );
}
