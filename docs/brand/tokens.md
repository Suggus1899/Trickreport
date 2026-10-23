# Tokens de diseño Trickreport

Fuente de verdad de los tokens de marca para las tres aplicaciones. Los nombres de variable son los que ya usa el paquete compartido (`packages/workspace/src/styles.css`); los valores base son los fijados en el plan maestro. Los tokens de estado amplían la base y siempre se acompañan de texto o icono, nunca solo de color.

## Color

| Token | Tema claro | Tema oscuro | Uso |
| --- | --- | --- | --- |
| `--accent` | `#4F46E5` | `#818CF8` | Navegación y acciones principales |
| `--accent-foreground` | `#FFFFFF` | `#0B1220` | Texto sobre `--accent` |
| `--green` | `#047857` | `#34D399` | Acciones complementarias y acentos |
| `--background` | `#F8FAFC` | `#0B1220` | Fondo principal |
| `--foreground` | `#0F172A` | `#E5E7EB` | Texto principal |
| `--danger` | `#B91C1C` | `#F87171` | Errores y acciones destructivas |
| `--warning` | `#B45309` | `#FBBF24` | Advertencias y vencimientos próximos |
| `--success` | `#047857` | `#34D399` | Confirmaciones y estados cumplidos |
| `--info` | `#1D4ED8` | `#60A5FA` | Información neutral y ayudas |

Las superficies, bordes y textos atenuados (`--surface`, `--surface-secondary`, `--border`, `--muted`) se derivan de esta base en cada aplicación y deben mantener contraste suficiente en ambos temas.

## Tipografía

| Familia | Uso |
| --- | --- |
| Public Sans | Interfaz y texto general |
| Space Mono | Identificadores, folios y datos técnicos (p. ej. `TI-001`) |

## Publicación

[`tokens.css`](tokens.css) publica los tokens base como variables CSS para consumo de las aplicaciones Next.js (HeroUI 3 / Tailwind CSS 4). Si un valor cambia aquí, cambia en las tres aplicaciones; ningún componente debe escribir colores de marca directamente.
