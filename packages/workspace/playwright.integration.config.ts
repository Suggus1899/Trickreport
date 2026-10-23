import { defineConfig } from '@playwright/test';
import { fileURLToPath } from 'node:url';

const databaseURL = process.env.E2E_DATABASE_URL;
if (!databaseURL) throw new Error('E2E_DATABASE_URL must point to a disposable PostgreSQL 16 database.');
const root = fileURLToPath(new URL('../..', import.meta.url));
const origin = 'http://127.0.0.1:3000';

export default defineConfig({
  testDir: '../../apps/it/tests/integration',
  workers: 1,
  timeout: 120_000,
  expect: { timeout: 15_000 },
  use: { baseURL: origin, headless: true, trace: 'retain-on-failure', screenshot: 'only-on-failure' },
  webServer: [
    {
      command: 'go run ./cmd/api',
      cwd: `${root}/backend`,
      url: 'http://127.0.0.1:8088/ready',
      timeout: 120_000,
      reuseExistingServer: false,
      env: {
        DATABASE_URL: databaseURL, PORT: '8088', ENV: 'development',
        JWT_SECRET: 'test-only-e2e-jwt-key-not-for-deployment', JWT_EXP_HOURS: '1', COOKIE_SECURE: 'false',
        APP_PUBLIC_URL: origin, CORS_ORIGINS: origin, LDAP_ENABLED: 'false',
        ADMIN_EMAIL: 'admin-e2e@example.test', ADMIN_PASSWORD: 'E2e.Admin.Only.2026!',
        SMTP_HOST: '127.0.0.1', SMTP_PORT: '1', SMTP_FROM: 'no-reply@example.test',
        SMTP_USERNAME: 'e2e', SMTP_PASSWORD: 'e2e', SMTP_USE_TLS: 'false',
      },
    },
    {
      command: 'npm run dev -w @trickreport/it -- --hostname 127.0.0.1',
      cwd: root, url: `${origin}/login`, timeout: 120_000, reuseExistingServer: false,
      env: { API_URL: 'http://127.0.0.1:8088', APP_PUBLIC_URL: origin },
    },
  ],
});
