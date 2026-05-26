import { expect, test } from "@playwright/test";

const adminEmail = process.env.ADMIN_EMAIL ?? "admin@example.com";
const adminPassword = process.env.ADMIN_PASSWORD ?? "Admin@12345";

test.describe("Gateway API smoke", () => {
  test("health returns ok", async ({ request }) => {
    const res = await request.get("/health");
    expect(res.ok()).toBeTruthy();
    await expect(res.json()).resolves.toEqual({ status: "ok" });
  });

  test("openapi spec is served", async ({ request }) => {
    const res = await request.get("/openapi.yaml");
    expect(res.ok()).toBeTruthy();
    const body = await res.text();
    expect(body).toContain("openapi: 3.0.3");
  });

  test("admin login and protected admin routes", async ({ request }) => {
    const loginRes = await request.post("/auth/login", {
      data: { email: adminEmail, password: adminPassword },
    });
    expect(loginRes.ok()).toBeTruthy();
    const { access_token: token } = (await loginRes.json()) as {
      access_token: string;
    };
    expect(token.length).toBeGreaterThan(10);

    const headers = { Authorization: `Bearer ${token}` };

    const channelsRes = await request.get("/admin/channels", { headers });
    expect(channelsRes.ok()).toBeTruthy();
    const channels = (await channelsRes.json()) as { items: unknown[] };
    expect(Array.isArray(channels.items)).toBeTruthy();

    const redeemStatsRes = await request.get("/admin/redeem/stats", { headers });
    expect(redeemStatsRes.ok()).toBeTruthy();

    const redeemListRes = await request.get("/admin/redeem/codes?page=1&page_size=10", {
      headers,
    });
    expect(redeemListRes.ok()).toBeTruthy();
  });

  test("unauthenticated admin route returns 401", async ({ request }) => {
    const res = await request.get("/admin/channels");
    expect(res.status()).toBe(401);
  });
});
