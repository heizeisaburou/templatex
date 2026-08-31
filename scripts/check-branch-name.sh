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

# Ramas de larga vida: no siguen el formato autor/timestamp.
# experiment/local-changes es la única rama de experimentación permitida
# en local según CONTRIBUTING.md (no se sube a dev).
case "$branch" in
main | dev | experiment/local-changes)
	exit 0
	;;
esac

fail() {
	echo "Nombre de rama inválido: $branch" >&2
	echo "" >&2
	echo "  $1" >&2
	echo "" >&2
	echo "Formato: autor/timestamp, sin puntos." >&2
	echo "Ejemplos: Misitox37/08302026, heizeisaburou/08302026, experiment/local-changes" >&2
	exit 1
}

# Una sola barra: autor/timestamp. Las barras extra crean subdirectorios
# dentro de .git/refs y complican la limpieza.
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

# Git materializa el nombre de la rama como un archivo real dentro de .git/refs
# y .git/logs/refs. Un punto convierte ese archivo en algo que las herramientas
# que recorren el repositorio por extensión intentan parsear: `gofmt -l .` falla
# sobre una referencia de Git llamada, por ejemplo, "fix/parser.go".
#
# Antes había una excepción para release/ (release/v11.00), pero ese formato
# es solo para commits según CONTRIBUTING.md, no para ramas. Ahora no hay
# excepciones: ningún nombre de rama puede contener un punto.
case "$purpose" in
*.*) fail "Contiene un punto. No uses nombres de archivo ni extensiones en las ramas." ;;
esac

# El timestamp debe ser exactamente 8 dígitos MMDDYYYY, ej: 08302026.
if ! printf '%s' "$purpose" | grep -Eq '^[0-9]{8}$'; then
	fail "El timestamp debe ser 8 dígitos MMDDYYYY (ej: 08302026)."
fi
