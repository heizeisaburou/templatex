# Convenio de errores

## Brief

Cómo se propagan los errores entre `pkg/` e `internal/`, y por qué `internal/document` envuelve los
errores de los paquetes en los que se apoya en lugar de dejarlos pasar tal cual.

## El problema

Los paquetes de `pkg/` son librerías genéricas e independientes entre sí, y cada una define sus
propios centinelas. Eso está bien: acoplarlas a un paquete de errores común sería peor. Pero implica
que en el repositorio conviven varios errores distintos con el mismo nombre y el mismo significado:

| Paquete | Centinela |
| --- | --- |
| `pkg/list` | `ErrOutOfBounds` |
| `pkg/cursor` | `ErrOutOfBounds` |
| `pkg/rnge` | `ErrInvalidRange` |
| `internal/document` | `ErrOutOfBounds` |

Son valores diferentes, así que `errors.Is` los distingue. Cuando `document.NewCursorAt` devolvía el
error de `pkg/cursor` sin tocarlo, esto era `false`:

```go
_, err := document.NewCursorAt(doc, -1)
errors.Is(err, document.ErrOutOfBounds) // false: el error venía de pkg/cursor
```

Quien usara `document` tenía que saber en qué paquete interno se había originado el fallo para poder
comprobarlo. Eso convierte un detalle de implementación en parte del contrato: el día que `document`
cambie de paquete de apoyo, el código de sus usuarios deja de detectar el error.

## La solución

`internal/document` actúa como frontera. Toda operación que por dentro se apoye en `pkg/rnge` o
`pkg/cursor` envuelve el error recibido con el centinela equivalente del propio paquete:

```go
func outOfBounds(err error) error {
	if err == nil {
		return nil
	}

	return fmt.Errorf("%w: %w", ErrOutOfBounds, err)
}
```

El doble `%w` es la clave. `errors.Is` acierta con el centinela de `document`, y el error original
sigue en la cadena, así que su mensaje no se pierde al depurar.

Devolver `nil` cuando `err` es `nil` permite envolver directamente el resultado de la llamada, sin
un `if` en cada punto:

```go
func (c *Cursor) Seek(position ByteOffset) error {
	return outOfBounds(c.cur.Seek(position))
}
```

## La regla

**Usa siempre los centinelas del paquete que llamas.** Si trabajas con `document`, comprueba contra
`document.ErrOutOfBounds` y `document.ErrInvalidRange`.

El error de origen permanece en la cadena por comodidad al depurar, pero **no forma parte del
contrato**: no escribas código que dependa de él.

Los tests de `internal/document/error_test.go` fijan este contrato para cada frontera del paquete.
Si añades una operación que llame a `pkg/`, envuélvela y añade ahí su caso.

## Errores frente a pánicos

El proyecto migró de pánicos a errores: `MustNewRange` y los pánicos de `pkg/rnge` y del antiguo
`pkg/clamp` se eliminaron, y los constructores y asignadores devuelven `error`.

Dos consecuencias que conviene tener presentes al escribir código nuevo:

- Un asignador que falla **no modifica el valor**. `Range.Set`, `Region.SetStart` y sus equivalentes
  validan antes de asignar, así que tras un error el rango sigue como estaba.
- En los tests, en lugar de `Must...`, se usan helpers locales que abortan la prueba: `mustNew`,
  `mustRange` y `mustRegion` en `internal/document`, `newToken` y `newRange` en `internal/ast`.

Queda una excepción pendiente de migrar: `NewPosition`, `Position.SetLine` y `Position.SetCol`
todavía entran en pánico con valores negativos.
