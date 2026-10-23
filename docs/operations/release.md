# Release

El release contiene binarios Go con archivos de configuración SLSA y un runtime Next.js standalone para el producto TI. El workflow de release comprueba el tag, reutiliza CI, genera el artefacto del frontend, lo arranca en un directorio limpio y publica checksum y procedencia.

El artefacto de frontend es Linux amd64 y requiere Node 24.22.1. Extrae el tarball, configura `API_URL`, `HOSTNAME` y `PORT`, y ejecuta:

```sh
node apps/it/server.js
```

Las migraciones se ejecutan desde el binario de API según el procedimiento de despliegue; el tarball de UI no modifica la base de datos. Conserva el release anterior hasta validar `/health`, `/ready`, login, una escritura autenticada y WebSocket en el entorno destino.

El workflow se actualizó para empaquetar Next.js, pero el artefacto todavía requiere una ejecución de CI/release sobre GitHub Actions. MSP y atención al cliente no forman parte de este runtime porque sus backends y despliegues independientes siguen pendientes.
