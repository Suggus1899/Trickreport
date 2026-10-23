import { getViteConfig } from 'astro/config';
import react from '@vitejs/plugin-react';
import { fileURLToPath } from 'node:url';

export default getViteConfig({
  plugins: [react()],
  resolve: {
    alias: {
      '@': fileURLToPath(new URL('./src', import.meta.url)),
      'astro:middleware': 'astro/middleware',
    },
  },
  test: {
    environment: 'node',
    globals: true,
    include: ['src/**/*.test.ts', 'src/**/*.test.tsx'],
    setupFiles: ['./vitest.setup.ts'],
    maxWorkers: 1,
  },
});
