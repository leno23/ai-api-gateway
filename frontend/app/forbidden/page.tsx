"use client";

import Link from "next/link";

import { Button, Result } from "antd";

export default function ForbiddenPage() {
  return (
    <div className="flex min-h-full flex-1 items-center justify-center p-6">
      <Result
        status="403"
        title="无权限"
        subTitle="当前账号不是管理员，或管理权限未开通。请联系运维在数据库中设置 role。"
        extra={
          <Link href="/login">
            <Button type="primary">返回登录</Button>
          </Link>
        }
      />
    </div>
  );
}
