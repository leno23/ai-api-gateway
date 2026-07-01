"use client";

import { Badge, Button, Dropdown, Nav } from "@douyinfe/semi-ui";
import { IconBell, IconDesktop, IconGlobe } from "@douyinfe/semi-icons";
import Link from "next/link";
import { usePathname, useRouter } from "next/navigation";
import { useMemo } from "react";

import { BrandLogo } from "@/components/brand-logo";
import { useAuth } from "@/contexts/auth-context";
import { usePortalLocale } from "@/contexts/portal-locale-context";
import type { PortalLocale } from "@/lib/i18n/messages";

export function SiteHeader() {
  const pathname = usePathname();
  const router = useRouter();
  const { token, user, logout } = useAuth();
  const { locale, setLocale, messages: m } = usePortalLocale();

  const navItems = useMemo(
    () => [
      { itemKey: "/", text: m.nav.home },
      { itemKey: "/console", text: m.nav.console },
      { itemKey: "/pricing", text: m.nav.pricing },
      { itemKey: "/docs", text: m.nav.docs },
      { itemKey: "/about", text: m.nav.about },
    ],
    [m],
  );

  const selectedKeys = navItems
    .filter((item) =>
      item.itemKey === "/" ? pathname === "/" : pathname.startsWith(item.itemKey),
    )
    .map((item) => item.itemKey);

  const setLanguage = (next: PortalLocale) => {
    setLocale(next);
  };

  const avatarLetter = (user?.display_name || user?.username || "U")
    .charAt(0)
    .toUpperCase();

  return (
    <header className="sticky top-0 z-50 border-b border-[var(--semi-color-border)] bg-white/95 backdrop-blur">
      <div className="mx-auto flex h-14 max-w-[1600px] items-center gap-4 px-4">
        <BrandLogo />
        <Nav
          mode="horizontal"
          selectedKeys={selectedKeys.length ? selectedKeys : ["/"]}
          className="min-w-0 flex-1 !border-none"
          items={navItems.map((item) => ({
            itemKey: item.itemKey,
            text: (
              <Link href={item.itemKey} className="text-inherit no-underline">
                {item.text}
              </Link>
            ),
          }))}
        />
        <div className="flex shrink-0 items-center gap-1">
          <Badge count={3} type="danger" overflowCount={9}>
            <Button
              theme="borderless"
              icon={<IconBell />}
              aria-label={m.nav.notifications}
            />
          </Badge>
          <Button
            theme="borderless"
            icon={<IconDesktop />}
            aria-label="显示模式"
          />
          <Dropdown
            trigger="click"
            position="bottomRight"
            render={
              <Dropdown.Menu>
                <Dropdown.Item
                  active={locale === "zh-CN"}
                  onClick={() => setLanguage("zh-CN")}
                >
                  简体中文
                </Dropdown.Item>
                <Dropdown.Item active={locale === "en"} onClick={() => setLanguage("en")}>
                  English
                </Dropdown.Item>
              </Dropdown.Menu>
            }
          >
            <Button theme="borderless" icon={<IconGlobe />} aria-label={m.nav.language} />
          </Dropdown>
          {token && user ? (
            <Dropdown
              trigger="click"
              position="bottomRight"
              render={
                <Dropdown.Menu>
                  <Dropdown.Item disabled>
                    {user.display_name || user.username}
                  </Dropdown.Item>
                  <Dropdown.Divider />
                  <Dropdown.Item onClick={() => router.push("/console/setting")}>
                    {m.nav.settings}
                  </Dropdown.Item>
                  <Dropdown.Item
                    onClick={() => {
                      logout();
                      router.push("/login");
                    }}
                  >
                    {m.nav.logout}
                  </Dropdown.Item>
                </Dropdown.Menu>
              }
            >
              <Button theme="borderless" className="!gap-2 !px-2">
                <span className="flex h-8 w-8 items-center justify-center rounded-full bg-teal-500 text-sm font-medium text-white">
                  {avatarLetter}
                </span>
                <span className="hidden max-w-[120px] truncate text-sm sm:inline">
                  {user.username}
                </span>
              </Button>
            </Dropdown>
          ) : (
            <>
              <Link href="/login">
                <Button theme="light">{m.nav.login}</Button>
              </Link>
              <Link href="/register">
                <Button theme="solid" type="primary">
                  {m.nav.register}
                </Button>
              </Link>
            </>
          )}
        </div>
      </div>
    </header>
  );
}
