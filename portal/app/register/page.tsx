"use client";

import { Button, Card, Form, Toast, Typography } from "@douyinfe/semi-ui";
import Link from "next/link";
import { useRouter, useSearchParams } from "next/navigation";
import { Suspense, useMemo, useState } from "react";

import { BrandLogo } from "@/components/brand-logo";
import { ApiError, gatewayFetch } from "@/lib/api/client";

const { Title, Text } = Typography;

function RegisterForm() {
  const router = useRouter();
  const searchParams = useSearchParams();
  const aff = useMemo(() => searchParams.get("aff") ?? "", [searchParams]);
  const [loading, setLoading] = useState(false);

  return (
    <div className="flex flex-1 items-center justify-center bg-[var(--portal-bg)] p-6">
      <Card className="portal-card w-full max-w-md !rounded-xl" shadows="always">
        <div className="mb-6 flex justify-center">
          <BrandLogo href="/" />
        </div>
        <Title heading={3} className="!mb-6 text-center">
          注册
        </Title>
        <Form
          initValues={{ invite_code: aff }}
          onSubmit={async (values) => {
            const v = values as {
              username: string;
              email: string;
              password: string;
              invite_code?: string;
            };
            setLoading(true);
            try {
              await gatewayFetch("/auth/register", {
                method: "POST",
                body: JSON.stringify({
                  username: v.username,
                  email: v.email,
                  password: v.password,
                  invite_code: v.invite_code || undefined,
                }),
              });
              Toast.success("注册成功，请登录");
              router.push("/login");
            } catch (e) {
              if (e instanceof ApiError) {
                Toast.error(e.message || "注册失败");
              } else {
                Toast.error("网络异常，请稍后重试");
              }
            } finally {
              setLoading(false);
            }
          }}
        >
          <Form.Input
            field="username"
            label="用户名"
            rules={[
              { required: true, message: "请输入用户名" },
              { min: 3, message: "至少 3 个字符" },
            ]}
          />
          <Form.Input
            field="email"
            label="邮箱"
            type="email"
            rules={[
              { required: true, message: "请输入邮箱" },
              { type: "email", message: "邮箱格式不正确" },
            ]}
          />
          <Form.Input
            field="password"
            label="密码"
            type="password"
            mode="password"
            rules={[
              { required: true, message: "请输入密码" },
              { min: 8, message: "至少 8 位" },
            ]}
          />
          {aff ? (
            <Form.Input field="invite_code" label="邀请码" disabled />
          ) : (
            <Form.Input field="invite_code" label="邀请码（可选）" showClear />
          )}
          <Button
            htmlType="submit"
            type="primary"
            theme="solid"
            block
            size="large"
            loading={loading}
            className="!mt-2"
          >
            注册
          </Button>
        </Form>
        <Text type="secondary" className="mt-6 block text-center">
          已有账号？{" "}
          <Link href="/login" className="text-[var(--portal-primary)]">
            继续登录
          </Link>
        </Text>
      </Card>
    </div>
  );
}

export default function RegisterPage() {
  return (
    <Suspense
      fallback={
        <div className="flex flex-1 items-center justify-center p-6">加载中…</div>
      }
    >
      <RegisterForm />
    </Suspense>
  );
}
