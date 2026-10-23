# Catálogo de casos de uso Trickreport

El catálogo define alcance futuro para los tres productos. Tener un caso documentado no significa que esté implementado ni comprobado. Los casos sólo pueden salir de `planned` cuando se enlacen código, contrato, pruebas y evidencia operativa.

| Producto | Casos | Datos | Lectura |
| --- | ---: | --- | --- |
| Soporte TI | 200 | [ti.json](ti.json) | [ti.md](ti.md) |
| Servicios MSP | 200 | [msp.json](msp.json) | [msp.md](msp.md) |
| Atención al cliente | 200 | [cx.json](cx.json) | [cx.md](cx.md) |

Los ejemplos usan organizaciones, clientes, personas y pedidos ficticios. Los permisos del catálogo describen validaciones requeridas en el servidor; el rol de interfaz nunca concede acceso por sí solo. Los contratos API están marcados como propuestos y se definirán al implementar cada caso.

Para regenerar los JSON y documentos Markdown desde la fuente, ejecutar `python .author.py` dentro de este directorio.
