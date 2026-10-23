# Migraciones y recuperación

La ejecución de estabilización no realizó backups ni restauraciones sobre datos existentes. Las pruebas de seguridad del script utilizan SQL sintético y Docker simulado. Un ensayo real de recuperación requiere un destino desechable y una autorización operativa propia.

## Migraciones

Aplicar migraciones mediante la API: `docker compose run --rm app --migrate-up`. El comando necesita las mismas variables y el mismo proyecto Compose del despliegue. El arranque de la API también aplica las migraciones pendientes. No ejecutar archivos SQL de migración manualmente ni editar una migración que ya esté aplicada.

Ensayar primero una instalación vacía y una actualización con el historial anterior de la instalación. Comparar versiones de migración, tablas, índices y relaciones de tenants. Las instalaciones anteriores a la migración inicial consolidada deben comprobarse expresamente; no asumir que una numeración similar acredita compatibilidad. El ensayo de upgrades y downgrades del esquema corresponde a las pruebas de base de datos, no al éxito de `compose config`.

`--migrate-down` revierte solamente la última migración y puede eliminar datos de esa capacidad. Ejecutarlo en el ensayo antes de aprobar una vuelta atrás y revisar la pareja de migraciones exacta. Nunca introducirlo como paso automático de rollback productivo.

## Copias programadas, solo cuando se habilitan

`docker-compose.backup.yml` es un overlay con perfil `backup`, excluido del despliegue por defecto. Usa el mismo `POSTGRES_DB`, `POSTGRES_USER` y secreto que la base, y una imagen PostgreSQL 16 fijada por digest. Conserva los intervalos existentes: diario, siete días, cuatro semanas y seis meses. El operador debe verificar permisos y espacio en el directorio montado antes de habilitarlo. La imagen requiere un filesystem POSIX con enlaces; una carpeta SMB/CIFS o exFAT no cumple ese contrato.

Antes de depender de una copia, comprobar su tamaño, `gzip -t`, fecha de la última copia completa y el resultado de un restore desechable. Un cron diario no acredita un RPO de 24 horas: las copias pueden fallar y el intervalo no cubre cambios posteriores al último dump íntegro. No se midieron RPO ni RTO en esta ejecución.

## Restore aislado

El script acepta únicamente dumps SQL de confianza generados por `pg_dump` sin `--create`, en texto o gzip. Un dump es código SQL ejecutable; comprobar su origen además de su integridad. La validación previa rechaza archivos vacíos, gzip corrupto, cambios de base e inclusiones de archivos/comandos. La restauración SQL completa valida sintaxis y datos dentro de una transacción.

La base destino debe ser **nueva**, tener un nombre explícito `trickreport_restore_<nombre>` de hasta 63 caracteres y ser distinta de la configurada en el servicio. `createdb` rechaza destinos existentes. El script no ejecuta DROP, no reemplaza un esquema y no toma credenciales ejecutando un archivo `.env` como shell.

Desde la raíz del repositorio, definir el proyecto y archivos del despliegue para que `docker compose exec` encuentre el servicio correcto:

```sh
export COMPOSE_FILE=docker-compose.yml:docker-compose.onpremise.yml
export COMPOSE_ENV_FILES=.env.onpremise
export COMPOSE_PROJECT_NAME=trickreport
bash scripts/restore.sh /ruta/dump.sql.gz trickreport_restore_rehearsal
```

En staging, seleccionar su overlay, archivo de entorno y nombre de proyecto. Si se usa Docker Compose en Windows, adaptar `COMPOSE_PATH_SEPARATOR` o los separadores de `COMPOSE_FILE`. El script descubre la red mediante el servicio `postgres`; no depende del nombre de una red Docker creada implícitamente.

Después de restaurar, comparar los conteos esperados, las relaciones entre usuarios/tickets/tenants, adjuntos y notificaciones, el historial de migraciones, sesiones y outbox. Probar el arranque de una API aislada y flujos de lectura y escritura usando ese destino. Registrar duración, tamaño, fecha del dump y pérdida temporal observada. Un SQL fallido deja el destino de ensayo para inspección; el origen permanece intacto.

Un corte posterior requiere detener escrituras, verificar una copia reciente si existe, acordar el impacto de datos y sesiones y cambiar el DSN de forma explícita. Este script termina en el ensayo y no sustituye automáticamente la base activa. La revisión de esos resultados precede cualquier aprobación del corte.

## Verificación local del script

`bash scripts/restore.test.sh` prueba rechazo antes de Docker de destino implícito o activo, archivo vacío, gzip corrupto y cambio de base; verifica el uso de una base nueva, transacción y propagación de un error SQL. No acredita permisos reales, espacio de disco, instalación Docker ni la integridad de una copia de negocio.
