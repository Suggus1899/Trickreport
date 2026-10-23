import { defineConfig } from '@playwright/test';

export default defineConfig({
  testDir: '../../apps/it/tests/e2e',
  use: { baseURL: 'http://127.0.0.1:3000', headless: true },
  webServer: {
    command: 'npm run dev -w @trickreport/it -- --hostname 127.0.0.1',
    url: 'http://127.0.0.1:3000/login',
    reuseExistingServer: !process.env.CI,
    timeout: 120_000,
  },
});
