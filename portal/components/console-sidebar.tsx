"use client";

import { Button, Layout, Nav } from "@douyinfe/semi-ui";
import {
  IconHistogram,
  IconKey,
  IconList,
  IconSetting,
  IconTicketCode,
  IconUser,
  IconChevronLeft,
} from "@douyinfe/semi-icons";
import Link from "next/link";
import { usePathname } from "next/navigation";
import { useState } from "react";

const { Sider } = Layout;

const MENU = [
  {
    itemKey: "chat",
    text: "聊天",
    items: [{ itemKey: "/console/playground", text: "操练场", icon: <IconUser /> }],
  },
  {
    itemKey: "console",
    text: "控制台",
    items: [
      { itemKey: "/console", text: "数据看板", icon: <IconHistogram /> },
      { itemKey: "/console/token", text: "令牌管理", icon: <IconKey /> },
      { itemKey: "/console/log", text: "使用日志", icon: <IconList /> },
      { itemKey: "/console/task", text: "任务日志", icon: <IconList /> },
    ],
  },
  {
    itemKey: "account",
    text: "个人中心",
    items: [
      { itemKey: "/console/wallet", text: "钱包管理", icon: <IconTicketCode /> },
      { itemKey: "/console/setting", text: "个人设置", icon: <IconSetting /> },
    ],
  },
];

function resolveSelectedKey(pathname: string): string {
  if (pathname.startsWith("/console/playground")) return "/console/playground";
  if (pathname.startsWith("/console/token")) return "/console/token";
  if (pathname.startsWith("/console/log")) return "/console/log";
  if (pathname.startsWith("/console/task")) return "/console/task";
  if (pathname.startsWith("/console/wallet")) return "/console/wallet";
  if (pathname.startsWith("/console/setting")) return "/console/setting";
  if (pathname.startsWith("/console")) return "/console";
  return "/console";
}

export function ConsoleSidebar() {
  const pathname = usePathname();
  const selectedKey = resolveSelectedKey(pathname);
  const [collapsed, setCollapsed] = useState(false);

  return (
    <Sider
      className="!bg-white !shadow-sm"
      style={{ width: collapsed ? 64 : 240 }}
    >
      <div className="flex h-full flex-col">
        <Nav
          style={{
            maxWidth: collapsed ? 64 : 240,
            height: "calc(100% - 48px)",
            flex: 1,
          }}
          selectedKeys={[selectedKey]}
          defaultOpenKeys={collapsed ? [] : ["chat", "console", "account"]}
          isCollapsed={collapsed}
          items={MENU.map((group) => ({
            itemKey: group.itemKey,
            text: group.text,
            items: group.items.map((item) => ({
              itemKey: item.itemKey,
              text: (
                <Link href={item.itemKey} className="text-inherit no-underline">
                  {item.text}
                </Link>
              ),
              icon: item.icon,
            })),
          }))}
        />
        <div className="border-t border-[var(--semi-color-border)] p-2">
          <Button
            theme="borderless"
            block
            icon={<IconChevronLeft className={collapsed ? "rotate-180" : ""} />}
            onClick={() => setCollapsed((c) => !c)}
          >
            {!collapsed ? "收起侧边栏" : null}
          </Button>
        </div>
      </div>
    </Sider>
  );
}
