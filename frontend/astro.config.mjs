import { defineConfig } from 'astro/config';
import node from '@astrojs/node';
import react from '@astrojs/react';
import tailwindcss from '@tailwindcss/vite';

export default defineConfig({
  output: 'server',

  adapter: node({
    mode: 'standalone',
  }),

  integrations: [react()],

  vite: {
    plugins: [tailwindcss()],

    resolve: {
      // Dependencies that bundle their own React (recharts, @base-ui/react) can
      // otherwise end up with a second copy at runtime, which breaks hooks with
      // "Invalid hook call" and leaves islands unhydrated.
      dedupe: ['react', 'react-dom'],
    },

    optimizeDeps: {
      // Every React dependency reachable from an island, listed up front.
      // Astro renders islands lazily, so Vite discovers these subpaths one
      // page at a time and re-runs the optimizer on each discovery, reloading
      // the page each time. Declaring them here collapses that into a single
      // pass at startup. Dev-server ergonomics only — production builds do not
      // use the dep optimizer at all.
      include: [
        'react',
        'react-dom',
        'react-dom/client',
        'react/jsx-runtime',
        '@base-ui/react/button',
        '@base-ui/react/input',
        '@base-ui/react/select',
        '@base-ui/react/merge-props',
        '@base-ui/react/use-render',
        'recharts',
        'qrcode',
        'lucide-react',
        'class-variance-authority',
        'clsx',
        'tailwind-merge',
      ],
      // Test-only deps: the dep scanner walks src/**/*.test.* and would pull
      // these into the browser graph, where they are never used.
      exclude: ['vitest', '@testing-library/react'],
    },
  },
});
