# Métricas y alertas

El overlay `docker-compose.monitoring.yml` incorpora Prometheus 3.15.1, Grafana 13.2.4 y Alertmanager 0.34.1 fijados por versión y digest. Los puertos 9090, 3001 y 9093 escuchan únicamente en `127.0.0.1`; acceder mediante túnel SSH o un proxy administrativo autenticado. Caddy rechaza `/metrics` en el origen público. Prometheus consulta `app:8080` por la red privada.

Definir `GRAFANA_PASSWORD` y `ALERT_WEBHOOK_URL_FILE`. El segundo valor apunta a un archivo privado con una URL de webhook autorizada para recibir el formato genérico de Alertmanager. No colocar esa URL en el repositorio. Compose monta el archivo como secret; `alertmanager.yml` utiliza `url_file` y envía también la resolución de las alertas. No existe una contraseña `admin` predeterminada ni un receptor vacío que simule notificaciones configuradas.

```sh
docker compose -f docker-compose.yml -f docker-compose.monitoring.yml --env-file .env.onpremise config --quiet
docker compose -f docker-compose.yml -f docker-compose.monitoring.yml --env-file .env.onpremise up -d
```

Las reglas montadas en `/etc/prometheus/rules.yml` detectan scrape fallido o ausente durante un minuto, más de 5% de respuestas 5xx durante cinco minutos, worker sin progreso por más de tres minutos durante dos minutos y trabajos de correo fallidos persistentes. La tasa usa `http_requests_total` emitido por la API; los workers publican `worker_last_success_timestamp_seconds{worker}`. `outbox_pending_total` y `outbox_failed_total` son gauges de trabajos actuales, pese al sufijo histórico `total`.

Grafana provisiona el datasource y el panel **Trickreport operations**, con solicitudes, latencia p95, errores HTTP, conexiones de DB, progreso de workers y outbox. Las rutas HTTP instrumentadas deben utilizar el patrón de chi para evitar etiquetas con IDs o email. Las series de DB provienen de `database_pool_connections{state}`.

## Validación

```sh
promtool check rules monitoring/rules.yml
promtool test rules monitoring/rules.test.yml
```

Los casos prueban caída y recuperación del backend, tasa alta, ausencia de tráfico, worker detenido y recuperado, y fallos de correo resueltos. Promtool 3.15.1 ejecutó esos casos satisfactoriamente en esta sesión. CI repite la validación con la misma versión. Las reglas tienen sensibilidad inicial; ajustarla con tráfico observado y sin ocultar fallos de entrega.

Antes de arrancar el stack, ejecutar también `promtool check config` y `amtool check-config` con las rutas montadas y el secret autorizado. Después, ensayar una caída controlada y una recuperación, comprobar que el receptor autorizado recibe ambas notificaciones y registrar el resultado. No se proporcionó un destino de prueba en esta ejecución y no se enviaron notificaciones externas. La entrega de alertas queda pendiente de esa prueba; YAML válido no prueba recepción.
