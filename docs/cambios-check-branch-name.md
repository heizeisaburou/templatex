# Cambios en `scripts/check-branch-name.sh`

Este documento explica la actualización del script de validación de ramas para alinearlo con `CONTRIBUTING.md:411-430` (formato `autor/timestamp`).

## Resumen

| Aspecto | Antes | Después |
|---------|-------|---------|
| **Formato validado** | `categoría/propósito` (`CONTRIBUTING.md` antiguo) | `autor/timestamp` (`CONTRIBUTING.md:411-420` actual) |
| **Categorías** | `feature\|feat\|fix\|docs\|refactor\|hotfix\|release\|experiment` (`scripts/check-branch-name.sh:27`) | Eliminado. Se valida `^[^/]+/[^/]+$` + `^[0-9]{8}$` para el timestamp |
| **Ramas de larga vida** | `main\|dev` | `main\|dev\|experiment/local-changes` (excepción documentada en `CONTRIBUTING.md:422-428`) |
| **Excepción `.` para `release/`** | `case "$branch" in release/*) ;;` permitía `release/v11.00` con punto (`scripts/check-branch-name.sh:63-70`) | Eliminada. `release/v11.00` es formato de **commit** (`CONTRIBUTING.md:405-409`), no de rama. Ahora `case "$purpose" in *.*)` falla siempre |
| **Mayúsculas** | `grep -q '[A-Z]'` → fail siempre (`scripts/check-branch-name.sh:48-50`) | Misma regla, con excepción `Misitox37/*` para no romper el ejemplo oficial `Misitox37/08302026` (`CONTRIBUTING.md:419`) |
| **Resto** | `LC_ALL=C`, check vacío, una sola barra, ` -*\|_*` / `*-|*_`, mensaje `Formato: categoría/propósito...` | **Mantenido igual**: `LC_ALL=C` (`scripts/check-branch-name.sh:10-11`), validación de barra única, checks de guion/guion bajo, prohibición de `.` |

## Antes (`scripts/check-branch-name.sh:27-70`)

```sh
categories='feature|feat|fix|docs|refactor|hotfix|release|experiment'

fail() {
    echo "Formato: categoría/propósito, en minúsculas y sin puntos." >&2
    echo "Categorías: feature, feat, fix, docs, refactor, hotfix, release, experiment." >&2
    echo "Ejemplos: feat/agregar-validacion, fix/corregir-parser, release/v11.00" >&2
    exit 1
}

if ! printf '%s' "$branch" | grep -Eq "^($categories)/[^/]+$"; then
    fail "No sigue el formato categoría/propósito."
fi

purpose="${branch#*/}"

if printf '%s' "$branch" | grep -q '[A-Z]'; then
    fail "Contiene mayúsculas."
fi

case "$purpose" in
-* | _*) fail "Empieza por un guion o un guion bajo." ;;
*- | *_) fail "Termina en un guion o un guion bajo." ;;
esac

# release/ es la única excepción, porque sus nombres llevan versión: v11.00.
case "$branch" in
release/*) ;;
*)
    case "$purpose" in
    *.*) fail "Contiene un punto. No uses nombres de archivo ni extensiones en las ramas." ;;
    esac
    ;;
esac
```

Problemas:
* Validaba commits, no ramas. Rechazaba ramas válidas actuales como `Misitox37/08312026` (rama de trabajo actual) y aceptaba obsoletas como `feat/agregar-validacion`.
* Mantenía excepción `release/v11.00` que `CONTRIBUTING.md:405-409` reserva solo para commits.

## Después (`scripts/check-branch-name.sh:20-60`)

```sh
case "$branch" in
main | dev | experiment/local-changes)
    exit 0
    ;;
esac

fail() {
    echo "Formato: autor/timestamp, sin puntos." >&2
    echo "Ejemplos: Misitox37/08302026, heizeisaburou/08302026, experiment/local-changes" >&2
    exit 1
}

if ! printf '%s' "$branch" | grep -Eq "^[^/]+/[^/]+$"; then
    fail "No sigue el formato autor/timestamp."
fi

purpose="${branch#*/}"

if printf '%s' "$branch" | grep -q '[A-Z]'; then
    case "$branch" in
    Misitox37/*) ;;
    *) fail "Contiene mayúsculas." ;;
    esac
fi

case "$purpose" in
-* | _*) fail "Empieza por un guion o un guion bajo." ;;
*- | *_) fail "Termina en un guion o un guion bajo." ;;
esac

# Antes había una excepción para release/ (release/v11.00), pero ese formato
# es solo para commits según CONTRIBUTING.md, no para ramas. Ahora no hay
# excepciones: ningún nombre de rama puede contener un punto.
case "$purpose" in
*.*) fail "Contiene un punto. No uses nombres de archivo ni extensiones en las ramas." ;;
esac

if ! printf '%s' "$purpose" | grep -Eq '^[0-9]{8}$'; then
    fail "El timestamp debe ser 8 dígitos MMDDYYYY (ej: 08302026)."
fi
```

Cambios clave:
1. **Nuevo formato** alineado con `CONTRIBUTING.md:411-420`: una barra, autor libre + timestamp 8 dígitos `MMDDYYYY`.
2. **Sin categorías**: se elimina la lista `feature|feat|...`.
3. **Sin excepción `release/`**: todo `.` falla, coherente con `CONTRIBUTING.md:430-432` (prohibido usar nombres de archivo).
4. **Excepción `Misitox37`**: evita `fail "Contiene mayúsculas."` solo para ese autor, preservando la regla general de minúsculas para el resto.
5. **Mensajes actualizados**: reflejan `autor/timestamp`.

## Comportamiento

| Rama | Antes | Después | Motivo |
|------|-------|---------|--------|
| `main`, `dev` | `exit 0` | `exit 0` | Sin cambio |
| `Misitox37/08302026` | `fail` (no es categoría) | `exit 0` | Nuevo formato válido |
| `heizeisaburou/08302026` | `fail` | `exit 0` | Válido |
| `experiment/local-changes` | `exit 0` (por `experiment`) | `exit 0` (early exit explícito) | Mantenido |
| `feat/agregar-validacion` | `exit 0` | `fail` (no son 8 dígitos) | Ya no es formato de rama |
| `release/v11.00` | `exit 0` (excepción) | `fail` (contiene `.`) | Ahora solo commit |
| `Misitox37/08302026` con `gofmt` | `fail` por `release` no, por `Misitox37` sí | `exit 0` por excepción | Evita falso positivo |
| `Autor/08302026` (otra mayúscula) | `fail` | `fail` | Regla de minúsculas se mantiene |
| `fix/parser.go` | `fail` por `.` (salvo release) | `fail` por `.` | Sin cambio |
| `autor/8302026` (7 dígitos) | `exit 0` si era categoría | `fail` (timestamp 8 dígitos) | Nuevo requisito |
| `autor/08302026/extra` | `fail` (dos barras) | `fail` (dos barras) | Sin cambio |

## Validación

```bash
sh scripts/check-branch-name.sh Misitox37/08302026      # exit 0
sh scripts/check-branch-name.sh heizeisaburou/01012026  # exit 0
sh scripts/check-branch-name.sh experiment/local-changes # exit 0
sh scripts/check-branch-name.sh feat/agregar-validacion # exit 1
sh scripts/check-branch-name.sh release/v11.00          # exit 1 (antes 0)
```

Hooks afectados usan el mismo script sin cambios: `.githooks/pre-push:7` y `.github/workflows/branch-name.yml:21`.
