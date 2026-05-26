import { expect, test } from "@playwright/test";

const adminEmail = process.env.ADMIN_EMAIL ?? "admin@example.com";
const adminPassword = process.env.ADMIN_PASSWORD ?? "Admin@12345";

/** Ant Design 中文按钮文案可能带空格，如「登 录」 */
const loginButton = (page: import("@playwright/test").Page) =>
  page.getByRole("button", { name: /登\s*录/ });

async function loginAsAdmin(page: import("@playwright/test").Page) {
  await page.goto("/login");
  await expect(page.getByRole("heading", { name: "管理员登录" })).toBeVisible();
  await page.getByLabel("邮箱").fill(adminEmail);
  await page.getByLabel("密码").fill(adminPassword);
  await Promise.all([
    page.waitForURL(/\/admin\/?$/, { timeout: 15_000 }),
    loginButton(page).click(),
  ]);
  await expect(page.getByRole("heading", { name: "概览" })).toBeVisible();
}

test.describe("Admin console UI smoke", () => {
  test("login page renders", async ({ page }) => {
    await page.goto("/login");
    await expect(page.getByRole("heading", { name: "管理员登录" })).toBeVisible();
    await expect(loginButton(page)).toBeEnabled();
  });

  test("home redirects unauthenticated user to login", async ({ page }) => {
    await page.goto("/");
    await expect(page).toHaveURL(/\/login/);
  });

  test("admin login and navigate core pages", async ({ page }) => {
    await loginAsAdmin(page);

    await expect(page.getByRole("heading", { name: "概览" })).toBeVisible();
    await expect(page.getByText("健康检查")).toBeVisible();

    await page.getByRole("link", { name: "渠道管理" }).click();
    await expect(page).toHaveURL(/\/admin\/channels/);
    await expect(page.getByRole("heading", { name: "渠道管理" })).toBeVisible();
    await expect(page.getByRole("button", { name: "新建渠道" })).toBeVisible();

    await page.getByRole("link", { name: "兑换运营" }).click();
    await expect(page).toHaveURL(/\/admin\/redeem/);
    await expect(page.getByRole("heading", { name: "兑换运营" })).toBeVisible();

    await page.getByRole("link", { name: "用户治理" }).click();
    await expect(page).toHaveURL(/\/admin\/users/);
    await expect(page.getByRole("heading", { name: "用户治理" })).toBeVisible();
  });

  test("logout returns to login", async ({ page }) => {
    await loginAsAdmin(page);
    await page.getByRole("button", { name: /退\s*出/ }).click();
    await expect(page).toHaveURL(/\/login/);
    await expect(page.getByRole("heading", { name: "管理员登录" })).toBeVisible();
  });
});
