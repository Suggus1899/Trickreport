/// <reference types="astro/client" />

import type { User } from './lib/api';

declare global {
  namespace App {
    interface Locals {
      token?: string;
      user: User | null;
    }
  }
}
