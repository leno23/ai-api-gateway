import { defineConfig, devices } from "@playwright/test";

const gatewayBase = process.env.GATEWAY_BASE_URL ?? "http://127.0.0.1:8080";
const frontendBase = process.env.FRONTEND_BASE_URL ?? "http://localhost:3000";

export default defineConfig({
  testDir: "./e2e",
  fullyParallel: true,
  forbidOnly: !!process.env.CI,
  retries: process.env.CI ? 1 : 0,
  workers: process.env.CI ? 1 : undefined,
  reporter: [["list"], ["html", { open: "never" }]],
  timeout: 30_000,
  use: {
    baseURL: frontendBase,
    trace: "on-first-retry",
    screenshot: "only-on-failure",
  },
  projects: [
    {
      name: "api",
      testMatch: /api\.smoke\.spec\.ts/,
      use: {
        baseURL: gatewayBase,
      },
    },
    {
      name: "ui",
      testMatch: /ui\.smoke\.spec\.ts/,
      use: {
        ...devices["Desktop Chrome"],
      },
    },
  ],
});
