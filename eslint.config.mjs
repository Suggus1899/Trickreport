import { defineConfig, globalIgnores } from 'eslint/config';
import nextVitals from 'eslint-config-next/core-web-vitals';
import nextTs from 'eslint-config-next/typescript';

export default defineConfig([
  ...nextVitals,
  ...nextTs,
  {
    files: ['apps/**/*.{js,jsx,ts,tsx}', 'packages/**/*.{js,jsx,ts,tsx}'],
    settings: { next: { rootDir: ['apps/it/', 'apps/msp/', 'apps/customer-service/'] } },
    rules: { '@next/next/no-html-link-for-pages': 'off' },
  },
  globalIgnores(['**/.next/**', '**/node_modules/**', '**/next-env.d.ts', 'frontend/**', 'backend/**']),
]);
