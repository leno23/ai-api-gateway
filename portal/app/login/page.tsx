"use client";

import { Button, Card, Form, Toast, Typography } from "@douyinfe/semi-ui";
import Link from "next/link";
import { useRouter, useSearchParams } from "next/navigation";
import { Suspense, useEffect, useState } from "react";

import { BrandLogo } from "@/components/brand-logo";
import { useAuth } from "@/contexts/auth-context";
import { ApiError, gatewayFetch } from "@/lib/api/client";
import type { LoginResponse } from "@/lib/api/types";

const { Title, Text } = Typography;

function LoginForm() {
  const { setToken, setUser } = useAuth();
  const router = useRouter();
  const searchParams = useSearchParams();
  const [loading, setLoading] = useState(false);

  useEffect(() => {
    if (searchParams.get("expired") === "true") {
      Toast.warning("未登录或登录已过期");
    }
  }, [searchParams]);

  return (
    <div className="flex flex-1 items-center justify-center bg-[var(--portal-bg)] p-6">
      <Card className="portal-card w-full max-w-md !rounded-xl" shadows="always">
        <div className="mb-6 flex justify-center">
          <BrandLogo href="/" />
        </div>
        <Title heading={3} className="!mb-6 text-center">
          登录
        </Title>
        <Form
          onSubmit={async (values) => {
            const v = values as { login: string; password: string };
            setLoading(true);
            try {
              const loginField = v.login.includes("@")
                ? { email: v.login }
                : { username: v.login };
              const data = await gatewayFetch<LoginResponse>("/auth/login", {
                method: "POST",
                body: JSON.stringify({
                  ...loginField,
                  password: v.password,
                }),
              });
              setToken(data.access_token);
              if (data.data) {
                setUser(data.data);
              }
              Toast.success("登录成功");
              router.replace("/console");
            } catch (e) {
              if (e instanceof ApiError && e.status === 401) {
                Toast.error("用户名或密码不正确");
              } else if (e instanceof ApiError) {
                Toast.error(e.message || "登录失败");
              } else {
                Toast.error("网络异常，请稍后重试");
              }
            } finally {
              setLoading(false);
            }
          }}
        >
          <Form.Input
            field="login"
            label="用户名或邮箱"
            placeholder="用户名 / email@example.com"
            rules={[{ required: true, message: "请输入用户名或邮箱" }]}
            showClear
          />
          <Form.Input
            field="password"
            label="密码"
            type="password"
            mode="password"
            rules={[{ required: true, message: "请输入密码" }]}
          />
          <div className="mb-4 text-right">
            <Link href="/login" className="text-sm text-[var(--portal-primary)]">
              忘记密码
            </Link>
          </div>
          <Button
            htmlType="submit"
            type="primary"
            theme="solid"
            block
            size="large"
            loading={loading}
          >
            继续
          </Button>
        </Form>
        <Text type="secondary" className="mt-6 block text-center">
          还没有账号？{" "}
          <Link href="/register" className="text-[var(--portal-primary)]">
            注册
          </Link>
        </Text>
      </Card>
    </div>
  );
}

export default function LoginPage() {
  return (
    <Suspense
      fallback={
        <div className="flex flex-1 items-center justify-center p-6">加载中…</div>
      }
    >
      <LoginForm />
    </Suspense>
  );
}
