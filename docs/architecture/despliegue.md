# Despliegue y alta disponibilidad

Dos topologías documentadas para los tres productos. La instalación de un solo servidor **no** cumple el objetivo SaaS acordado en el plan maestro.

## Topología autohospedada sencilla

- Docker Compose con Caddy, frontend, API y PostgreSQL en un servidor.
- Adecuada para una organización que se autohospeda o para evaluación.
- Caddy actúa como ingreso y terminación TLS.

## Topología SaaS con alta disponibilidad

- **Docker Swarm** con al menos dos réplicas por servicio de aplicación.
- **Tres dominios de fallo**: un nodo de aplicación y uno de datos en cada uno (seis servidores o VPS en la referencia).
  - Los tres nodos de aplicación mantienen el quorum de Swarm.
  - Los tres nodos de datos mantienen PostgreSQL y etcd.
- **Ingreso redundante** con HAProxy y Keepalived sobre una IP virtual; la infraestructura elegida debe admitir esa configuración de red.
- **PostgreSQL con Patroni y etcd**, replicación síncrona estricta con al menos una réplica confirmando escrituras: se conservan las escrituras confirmadas y el servicio se recupera automáticamente tras perder un nodo.
- El tamaño de las máquinas se fijará con las pruebas de capacidad del plan (15.000 cuentas, 1.500 sesiones concurrentes; objetivo de recuperación ≤ 60 s ante la caída de un nodo).

## Estado compartido entre réplicas

Las sesiones, los trabajos y las notificaciones dependen de estado persistente, no de memoria local de una réplica:

- Coordinación con PostgreSQL: trabajos con *leases*, límites que requieren coordinación entre réplicas y avisos mediante `LISTEN/NOTIFY`.
- El estado recuperable vive en tablas; una réplica puede reiniciarse sin perder trabajo.
- Las páginas privadas se sirven sin caché compartida; todas las réplicas de Next.js usan la misma imagen por versión, con identificación de despliegue y conservación temporal de los recursos estáticos de la versión anterior para no romper sesiones en curso durante un despliegue gradual.

## Gestión comercial inicial

Altas, planes y límites se gestionan de forma manual en la primera etapa comercial; los pagos y las suscripciones automáticas quedan para una etapa posterior.

## Referencia versionada (Etapa 2, Soporte TI)

- **Sencilla**: `docker-compose.yml` y `Caddyfile` de la raíz ya sirven el
  frontend Next.js de Soporte TI (`apps/it`) junto a la API y PostgreSQL;
  Caddy enruta `/api/v1/ws` directamente a la API.
- **SaaS HA**: [`deploy/ha/`](../../deploy/ha/README.md) contiene la
  referencia versionada y comentada: pila Swarm de la capa de aplicación
  (`docker-stack.yml`), Patroni con replicación síncrona estricta
  (`patroni/patroni.yml`), HAProxy con enrutado por rol de Patroni
  (`haproxy/haproxy.cfg`) y Keepalived con la IP virtual
  (`keepalived/keepalived.conf`), con el orden de despliegue y la
  validación exigida. Es configuración de referencia: no está validada
  sobre infraestructura real y la alta disponibilidad debe aprobarse en
  la topología real antes de publicar el SaaS.
