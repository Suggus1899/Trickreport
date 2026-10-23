import js from '@eslint/js';
import astro from 'eslint-plugin-astro';
import typescriptParser from '@typescript-eslint/parser';
export default [
  { ignores: ['dist/', '.astro/', 'node_modules/'] },
  js.configs.recommended,
  ...astro.configs.recommended,
  {
    files: ['**/*.{ts,tsx}'],
    languageOptions: {
      parser: typescriptParser,
      parserOptions: { ecmaFeatures: { jsx: true } },
    },
    rules: { 'no-undef': 'off' },
  },
  {
    files: ['**/*.astro'],
    languageOptions: { parserOptions: { parser: typescriptParser } },
  },
  {
    files: ['public/sw.js'],
    languageOptions: {
      globals: Object.fromEntries(['self', 'caches', 'fetch', 'Request', 'Response', 'URL'].map((name) => [name, 'readonly'])),
    },
  },
  { rules: { 'no-unused-vars': ['warn', { args: 'none', caughtErrors: 'none', varsIgnorePattern: '^_' }] } },
];
