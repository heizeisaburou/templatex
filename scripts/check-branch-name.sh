#!/bin/sh
# Valida el nombre de una rama según la nomenclatura de CONTRIBUTING.md.
#
# Uso: scripts/check-branch-name.sh NOMBRE_DE_RAMA
set -eu

# Las clases como [a-z] dependen de la configuración regional: en en_US.UTF-8
# incluyen las letras acentuadas y en C.UTF-8 no. Se fija LC_ALL para que el
# hook local y el runner de CI apliquen exactamente el mismo criterio.
LC_ALL=C
export LC_ALL

branch="${1:-}"

if [ -z "$branch" ]; then
	echo "uso: $0 NOMBRE_DE_RAMA" >&2
	exit 2
fi

# Ramas de larga vida: no siguen el formato categoría/propósito.
case "$branch" in
main | dev)
	exit 0
	;;
esac

categories='feature|feat|fix|docs|refactor|hotfix|release|experiment'

fail() {
	echo "Nombre de rama inválido: $branch" >&2
	echo "" >&2
	echo "  $1" >&2
	echo "" >&2
	echo "Formato: categoría/propósito, en minúsculas y sin puntos." >&2
	echo "Categorías: feature, feat, fix, docs, refactor, hotfix, release, experiment." >&2
	echo "Ejemplos: feat/agregar-validacion, fix/corregir-parser, release/v11.00" >&2
	exit 1
}

# Una sola barra: categoría/propósito. Las barras extra crean subdirectorios
# dentro de .git/refs y complican la limpieza.
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

# Git materializa el nombre de la rama como un archivo real dentro de .git/refs
# y .git/logs/refs. Un punto convierte ese archivo en algo que las herramientas
# que recorren el repositorio por extensión intentan parsear: `gofmt -l .` falla
# sobre una referencia de Git llamada, por ejemplo, "fix/parser.go".
#
# release/ es la única excepción, porque sus nombres llevan versión: v11.00.
case "$branch" in
release/*) ;;
*)
	case "$purpose" in
	*.*) fail "Contiene un punto. No uses nombres de archivo ni extensiones en las ramas." ;;
	esac
	;;
esac
