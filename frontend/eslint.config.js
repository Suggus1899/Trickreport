import js from '@eslint/js';
import astro from 'eslint-plugin-astro';
export default [
  js.configs.recommended,
  ...astro.configs.recommended,
  { ignores: ['dist/', '.astro/', 'node_modules/'] },
  { rules: { 'no-unused-vars': 'warn', 'no-undef': 'off' } },
];
