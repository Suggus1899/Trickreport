# Autenticación — SSO / LDAP

## Flujo de Login

1. Usuario ingresa en la página de login
2. Opción A: **SSO** — Redirige a proveedor OAuth2 (Google, GitHub, Azure AD, etc.)
3. Opción B: **LDAP** — Ingresa credenciales corporativas, se validan contra LDAP
4. Backend genera JWT firmado
5. Frontend almacena token (httpOnly cookie o memory)
6. Middleware en Go valida JWT en cada request

## JWT

| Claim | Descripción |
|-------|-------------|
| sub | user_id |
| role | admin, agent, end_user |
| exp | Expiración (8h por defecto) |

## RBAC (Role-Based Access Control)

| Recurso | Admin | Agent | End User |
|---------|-------|-------|----------|
| Ver todos los tickets | ✓ | ✓ | Solo propios |
| Crear ticket | ✓ | ✓ | ✓ |
| Asignar ticket | ✓ | ✓ | ✗ |
| Cambiar estado | ✓ | ✓ | Solo cerrar propios |
| Admin usuarios | ✓ | ✗ | ✗ |
| Configurar SLA | ✓ | ✗ | ✗ |
| Escribir artículos KB | ✓ | ✓ | ✗ |
| Dashboard global | ✓ | ✓ (parcial) | ✗ |

## LDAP

- Configurable por variables de entorno (host, puerto, bind DN, base DN)
- Búsqueda de usuario por uid o mail
- Sincronización opcional: al primer login, se crea/actualiza el usuario local
