"use client";

import {
  ApiOutlined,
  DashboardOutlined,
  GiftOutlined,
  LogoutOutlined,
  UserOutlined,
} from "@ant-design/icons";
import { Button, Layout, Menu, Typography } from "antd";
import Link from "next/link";
import { usePathname } from "next/navigation";

const { Header, Sider, Content } = Layout;

const menuItems = [
  {
    key: "/admin",
    icon: <DashboardOutlined />,
    label: <Link href="/admin">概览</Link>,
  },
  {
    key: "/admin/channels",
    icon: <ApiOutlined />,
    label: <Link href="/admin/channels">渠道管理</Link>,
  },
  {
    key: "/admin/redeem",
    icon: <GiftOutlined />,
    label: <Link href="/admin/redeem">兑换运营</Link>,
  },
  {
    key: "/admin/users",
    icon: <UserOutlined />,
    label: <Link href="/admin/users">用户治理</Link>,
  },
];

export function AdminShell({
  children,
  onLogout,
}: {
  children: React.ReactNode;
  onLogout: () => void;
}) {
  const pathname = usePathname();
  const selected = menuItems.some((i) => i.key === pathname)
    ? pathname
    : pathname.startsWith("/admin/channels")
      ? "/admin/channels"
      : pathname.startsWith("/admin/redeem")
        ? "/admin/redeem"
        : pathname.startsWith("/admin/users")
          ? "/admin/users"
          : "/admin";

  return (
    <Layout className="h-screen min-h-screen overflow-hidden bg-neutral-50">
      <Sider
        breakpoint="lg"
        collapsedWidth={0}
        theme="light"
        width={220}
        className="!h-full !min-h-0 border-r border-neutral-100"
      >
        <div className="flex h-14 items-center justify-center border-b border-neutral-100 px-3">
          <Typography.Text strong className="truncate">
            Gateway 控制台
          </Typography.Text>
        </div>
        <Menu mode="inline" selectedKeys={[selected]} items={menuItems} className="border-none" />
      </Sider>
      <Layout className="flex min-h-0 flex-1 flex-col">
        <Header className="flex shrink-0 items-center justify-end gap-3 border-b border-neutral-100 bg-white px-6">
          <Button type="text" icon={<LogoutOutlined />} onClick={onLogout}>
            退出
          </Button>
        </Header>
        <Content className="m-4 flex-1 overflow-auto rounded-lg bg-white p-6 shadow-sm">
          {children}
        </Content>
      </Layout>
    </Layout>
  );
}
