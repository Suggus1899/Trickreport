import { expect, test } from '@playwright/test';

test('IT login page is usable without a backend session', async ({ page }) => {
  await page.goto('/login');
  await expect(page.getByRole('heading', { name: 'Tu equipo, conectado.' })).toBeVisible();
  await expect(page.getByLabel('Correo electrónico')).toBeVisible();
  await expect(page.getByLabel('Contraseña')).toBeVisible();
  await expect(page.getByRole('link', { name: 'Olvidé mi contraseña' })).toHaveAttribute('href', '/password-reset');
});
