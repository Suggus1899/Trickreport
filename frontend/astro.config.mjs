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
    // Dependencies that bundle their own React (recharts, @base-ui/react) can
    // otherwise end up with a second copy at runtime, which breaks hooks with
    // "Invalid hook call" and leaves islands unhydrated.
    resolve: {
      dedupe: ['react', 'react-dom'],
    },
  },
});