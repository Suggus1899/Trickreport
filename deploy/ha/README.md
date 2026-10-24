# Referencia de alta disponibilidad (SaaS) — Soporte TI

Configuración **de referencia** para la topología SaaS del plan maestro
(§3, «Despliegue y alta disponibilidad»). Está versionada para revisión y
adaptación; **no está validada todavía sobre infraestructura real** y la
alta disponibilidad debe aprobarse en la topología real antes de publicar
el SaaS (condición de salida de la Etapa 2 del plan). La instalación de un
solo servidor (`docker-compose.yml` + `Caddyfile` de la raíz) sigue siendo
la modalidad autohospedada sencilla y **no** cumple el objetivo SaaS.

## Topología (tres dominios de fallo, seis nodos)

| Dominio | Nodo de aplicación | Nodo de datos |
|---|---|---|
| A | Swarm manager + HAProxy/Keepalived | PostgreSQL (Patroni) + etcd |
| B | Swarm manager + HAProxy/Keepalived | PostgreSQL (Patroni) + etcd |
| C | Swarm manager + HAProxy/Keepalived | PostgreSQL (Patroni) + etcd |

- Los tres nodos de aplicación mantienen el quórum de Docker Swarm y
  ejecutan las réplicas de la API y del frontend Next.js de Soporte TI
  (`docker-stack.yml`).
- Los tres nodos de datos mantienen PostgreSQL con Patroni y etcd
  (`patroni/patroni.yml`): replicación síncrona estricta, al menos una
  réplica confirma cada escritura; promoción automática al perder un nodo
  (objetivo de recuperación del servicio ≤ 60 s).
- El ingreso usa HAProxy + Keepalived con una IP virtual (VIP) sobre los
  nodos de aplicación (`haproxy/`, `keepalived/`). La infraestructura
  elegida debe admitir la VIP (red L2 o proveedor que lo permita).
- La API llega a PostgreSQL a través de HAProxy (puerto 5000, solo
  primario verificado contra la API REST de Patroni), nunca a un nodo
  concreto: así el failover de Patroni no exige reconfigurar la API.

## Orden de despliegue

1. **Datos**: instala etcd (3 nodos) y Patroni/PostgreSQL 16 con
   `patroni/patroni.yml` adaptado (scope `trickreport-it`, direcciones y
   secretos reales). Verifica con `patronictl list` que hay un primario y
   dos réplicas, una de ellas síncrona.
2. **Ingreso de datos**: HAProxy en los nodos de aplicación con
   `haproxy/haproxy.cfg`; comprueba que el puerto 5000 responde solo
   cuando Patroni marca el nodo como primario.
3. **Swarm**: `docker swarm init` en el primer nodo de aplicación y join
   de los otros dos como managers. Prepara las variables de entorno que
   la pila exige (`DATABASE_URL` apuntando al endpoint primario de
   HAProxy, `JWT_SECRET`, `ADMIN_EMAIL`, `ADMIN_PASSWORD`,
   `APP_PUBLIC_URL`, …) desde la gestión de secretos del operador; no se
   versionan en el repositorio.
4. **Aplicación**: publica las imágenes con la misma etiqueta de versión
   para API y frontend y despliega `docker-stack.yml` sustituyendo
   `VERSION`. Los archivos estáticos de la versión anterior de Next.js se
   conservan durante el despliegue gradual (ver la guía de alojamiento
   propio de Next.js enlazada en el plan maestro), por eso la pila usa
   `update_config` con `parallelism: 1` y `delay` entre réplicas.
5. **Ingreso público**: Keepalived con la VIP y HAProxy terminando TLS
   (o un terminador equivalente) hacia el puerto publicado de la pila.

## Validación exigida antes de publicar (plan maestro §5, «Operación»)

- Reinicio de una réplica de aplicación sin pérdida de sesiones.
- Caída de un nodo de aplicación: el servicio sigue y Swarm reprograma
  las réplicas.
- Caída del primario de PostgreSQL: promoción automática ≤ 60 s, sin
  pérdida de escrituras confirmadas ni duplicación de operaciones.
- Reconexión automática de clientes (incluido WebSocket `/api/v1/ws`).
- Despliegue gradual a una versión nueva y retorno a la anterior
  manteniendo el servicio disponible.
- Prueba de capacidad del plan (15.000 cuentas, 1.500 sesiones
  concurrentes) ejecutada sobre esta topología, con resultados publicados.

## Pendiente

- Adaptar direcciones, nombres de host, VIP y secretos a la infraestructura
  concreta que se elija.
- Automatizar la instalación de los nodos de datos (hoy es manual y
  guiada por esta referencia).
- Aprobar la topología con la validación anterior y registrar la evidencia
  en el registro de la etapa correspondiente.
