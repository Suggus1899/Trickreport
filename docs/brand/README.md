# Identidad de marca Trickreport

La marca combina seguimiento claro de solicitudes con una señal de avance. El símbolo es una **T** geométrica —la inicial de Trickreport y la forma de un registro— con un trazo de verificación que sugiere progreso y resolución. Debe seguir siendo reconocible a tamaño pequeño.

## Símbolo y logotipo

| Archivo | Uso |
| --- | --- |
| `simbolo-color.svg` | Versión principal a color: fondo índigo, T en blanco roto y verificación en verde claro. |
| `simbolo-oscuro.svg` | Para temas y fondos oscuros. Es el mismo símbolo publicado en `apps/*/public/trickreport-mark.svg`. |
| `simbolo-claro.svg` | Para fondos claros: T índigo y verificación verde oscuro sobre fondo blanco roto. |
| `simbolo-mono.svg` / `simbolo-mono-blanco.svg` | Monocromo (oscuro o blanco), sin fondo: documentos y usos de un solo color. |
| `logotipo-claro.svg` / `logotipo-oscuro.svg` | Símbolo + nombre tipográfico `trickreport` en Public Sans, para fondos claros u oscuros. |
| `favicon.svg` | Favicon de las aplicaciones (copiado también en `apps/*/public/favicon.svg`). |

Las exportaciones PNG (símbolo a 512 px, logotipos, icono de aplicación de 180 px y favicon de 32 px) se generan desde los SVG cuando se necesitan; ver «Regeneración de PNG» al final de esta guía.

Los tres productos usan la misma marca con un descriptor funcional junto al logotipo: «Soporte TI», «MSP» y «Atención al cliente». El descriptor va en texto, nunca dentro del símbolo.

### Reglas de uso

- Conservar proporción y contraste; no deformar, no recolorear fuera de los tokens y no añadir texto dentro del símbolo.
- Tamaño mínimo del símbolo: 16 px en pantalla; por debajo, usar `favicon.svg`.
- Espacio libre alrededor del logotipo: al menos la altura de la barra de la T por cada lado.
- El logotipo se compone en Public Sans (600). Si la fuente no está disponible, el SVG sigue siendo la referencia válida; no sustituir la tipografía en usos oficiales.

## Color y tipografía

| Uso | Color | Token |
| --- | --- | --- |
| Primario claro | Indigo `#4F46E5` | `--accent` |
| Secundario claro | Verde `#047857` | `--green` |
| Fondo claro | `#F8FAFC` | `--background` |
| Texto claro | `#0F172A` | `--foreground` |
| Primario oscuro | Indigo claro `#818CF8` | `--accent` en tema oscuro |
| Secundario oscuro | Verde claro `#34D399` | `--green` en tema oscuro |
| Fondo oscuro | Azul tinta `#0B1220` | `--background` en tema oscuro |
| Texto oscuro | `#E5E7EB` | `--foreground` en tema oscuro |

Tipografía: Public Sans para la interfaz y Space Mono para identificadores técnicos. El sistema completo de tokens (incluidos estado, superficies y bordes) está en [`tokens.md`](tokens.md) y [`tokens.css`](tokens.css).

## Regeneración de PNG

Los PNG se generan desde los SVG, que son la fuente de verdad. Con una herramienta SVG→PNG (p. ej. `cairosvg`):

```bash
python -c "import cairosvg; cairosvg.svg2png(url='simbolo-oscuro.svg', write_to='simbolo-oscuro-512.png', output_width=512, output_height=512)"
```
