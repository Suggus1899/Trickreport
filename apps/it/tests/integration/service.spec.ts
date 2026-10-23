import { expect, test, type Page } from '@playwright/test';

const password = 'E2e.Requester.2026!';

async function register(page: Page, email: string) {
  await page.goto('/register');
  await page.getByLabel('Nombre completo').fill('Prueba solicitante');
  await page.getByLabel('Correo electrónico').fill(email);
  await page.getByLabel('Contraseña', { exact: true }).fill(password);
  await page.getByRole('button', { name: 'Crear cuenta', exact: true }).click();
  await expect(page).toHaveURL(/\/dashboard$/);
  await expect(page.getByRole('heading', { name: 'Hola, Prueba' })).toBeVisible();
}

test('requester creates, discusses, attaches, saves a draft and cannot read another account’s ticket', async ({ page, browser }) => {
  const stamp = crypto.randomUUID();
  await register(page, `requester-${stamp}@example.test`);
  await page.getByRole('link', { name: 'Nueva solicitud' }).click();
  await page.getByLabel('Título', { exact: true }).fill(`Acceso VPN ${stamp}`);
  await page.getByLabel('Descripción').fill('Necesito recuperar el acceso a la red de pruebas.');
  await page.getByLabel('Categoría', { exact: true }).fill('Acceso');
  await page.getByRole('button', { name: 'Crear solicitud', exact: true }).click();
  await expect(page).toHaveURL(/\/tickets\/[a-f0-9-]{36}$/);
  const ticketPath = new URL(page.url()).pathname;
  await expect(page.getByRole('heading', { name: `Acceso VPN ${stamp}` })).toBeVisible();
  await page.getByLabel('Tu respuesta').fill('Adjunto el detalle para que el equipo lo revise.');
  await page.getByRole('button', { name: 'Enviar respuesta' }).click();
  await expect(page.getByText('Adjunto el detalle para que el equipo lo revise.', { exact: true })).toBeVisible();
  await page.getByLabel('Seleccionar archivo').setInputFiles({
    name: 'evidence.txt', mimeType: 'text/plain', buffer: Buffer.from('Evidencia ficticia de la prueba de acceso.'),
  });
  await page.getByRole('button', { name: 'Subir archivo' }).click();
  const attachment = page.getByRole('link', { name: 'evidence.txt' });
  await expect(attachment).toBeVisible();
  const download = await page.request.get((await attachment.getAttribute('href'))!);
  expect(download.status()).toBe(200);
  expect(await download.text()).toBe('Evidencia ficticia de la prueba de acceso.');
  expect(download.headers()['content-disposition']).toContain('attachment');

  await page.goto('/tickets/new');
  await page.getByLabel('Título', { exact: true }).fill(`Borrador ${stamp}`);
  await page.getByLabel('Descripción').fill('Solicitud guardada de forma local para sincronizar después.');
  await page.getByLabel('Categoría', { exact: true }).fill('Software');
  await page.getByRole('button', { name: 'Guardar borrador local' }).click();
  await expect(page).toHaveURL(/\/tickets\/drafts$/);
  await expect(page.getByRole('heading', { name: `Borrador ${stamp}` })).toBeVisible();
  await page.getByRole('button', { name: 'Sincronizar', exact: true }).click();
  await expect(page.getByText('No tienes borradores pendientes.')).toBeVisible();
  await page.goto('/tickets');
  await expect(page.getByRole('link', { name: `Borrador ${stamp}` })).toBeVisible();
  await expect(page.getByRole('link', { name: 'Personas', exact: true })).toHaveCount(0);

  const other = await browser.newContext();
  try {
    const otherPage = await other.newPage();
    await register(otherPage, `other-${stamp}@example.test`);
    const denied = await otherPage.request.get(`/api/v1${ticketPath}`);
    expect([403, 404]).toContain(denied.status());
    const adminDenied = await otherPage.request.get('/api/v1/admin/users');
    expect(adminDenied.status()).toBe(403);
    await otherPage.goto('/tickets/drafts');
    await expect(otherPage.getByText('No tienes borradores pendientes.')).toBeVisible();
  } finally { await other.close(); }

  await page.getByRole('button', { name: 'Cerrar sesión', exact: true }).click();
  await expect(page).toHaveURL(/\/login$/);
  expect((await page.request.get(`/api/v1${ticketPath}`)).status()).toBe(401);
});

test('administrator publishes knowledge, creates a valid automation, manages SLA and uses the mobile navigation', async ({ page }) => {
  const stamp = crypto.randomUUID();
  await page.goto('/login');
  await page.getByLabel('Correo electrónico').fill('admin-e2e@example.test');
  await page.getByLabel('Contraseña', { exact: true }).fill('E2e.Admin.Only.2026!');
  await page.getByRole('button', { name: 'Entrar al espacio' }).click();
  await expect(page).toHaveURL(/\/dashboard$/);
  await page.goto('/articles/new');
  await page.getByLabel('Título', { exact: true }).fill(`Guía de prueba ${stamp}`);
  await page.getByLabel('Categoría', { exact: true }).fill('Acceso');
  await page.getByLabel('Contenido', { exact: true }).fill('Comprueba tu conexión y vuelve a intentarlo con tu cuenta de pruebas.');
  await page.getByLabel('Publicar para los solicitantes').check();
  await page.getByRole('button', { name: 'Guardar artículo' }).click();
  await expect(page.getByRole('heading', { name: `Guía de prueba ${stamp}` })).toBeVisible();

  await page.goto('/admin/automations');
  await page.getByRole('button', { name: 'Nueva regla' }).click();
  await page.getByLabel('Nombre', { exact: true }).fill(`Seguimiento ${stamp}`);
  await page.getByLabel('Evento', { exact: true }).selectOption('status_changed');
  await page.getByRole('button', { name: 'Guardar regla' }).click();
  await expect(page.getByRole('heading', { name: `Seguimiento ${stamp}` })).toBeVisible();

  await page.goto('/admin/sla');
  const high = page.locator('section.panel').filter({ has: page.getByRole('heading', { name: 'Alta', exact: true }) });
  await high.getByLabel('Primera respuesta', { exact: true }).fill('15');
  await high.getByLabel('Resolución', { exact: true }).fill('240');
  await high.getByLabel('Escalación', { exact: true }).fill('30');
  await high.getByRole('button', { name: 'Guardar política' }).click();
  await expect(page.getByText('Política SLA guardada.')).toBeVisible();
  await page.reload();
  await expect(high.getByLabel('Resolución', { exact: true })).toHaveValue('240');

  await page.setViewportSize({ width: 390, height: 844 });
  await page.getByRole('button', { name: 'Open navigation' }).click();
  await page.getByRole('link', { name: 'Solicitudes', exact: true }).click();
  await expect(page.getByRole('heading', { name: 'Solicitudes', exact: true })).toBeVisible();
  expect(await page.evaluate(() => document.documentElement.scrollWidth <= window.innerWidth)).toBe(true);
  await page.getByLabel('Idioma / Language').selectOption('en');
  await expect(page.getByRole('heading', { name: 'Requests', exact: true })).toBeVisible();
  await page.getByRole('button', { name: 'Switch to dark theme' }).click();
  await expect(page.locator('html')).toHaveClass(/dark/);
});
