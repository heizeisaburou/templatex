# Configuración del clon

## Brief

Lo que hay que configurar **una vez por clon** después de clonar el repositorio, y por qué existe
cada cosa. Ninguno de estos ajustes avisa cuando falta: simplemente algo deja de funcionar más
tarde, a veces en la máquina de otra persona, y cuesta relacionar una cosa con la otra.

## Activar los hooks

```bash
git config core.hooksPath .githooks
```

Git solo ejecuta hooks desde `.git/hooks/`, y `.git/` no se clona nunca. El directorio `.githooks/`
del repositorio se descarga con el resto del código, pero Git no lo mira hasta que se lo indicas.

No es un descuido de Git, es deliberado: un hook es código ejecutable, y ejecutar automáticamente el
de un repositorio recién clonado convertiría `git clone` en «ejecuta código ajeno en tu máquina».
Por el mismo motivo, `core.hooksPath` no puede fijarse desde el propio repositorio.

El hook `pre-push` valida el nombre de la rama antes de subirla. Es una comodidad, no una barrera:
se salta con `--no-verify` y depende de que cada persona lo active. La comprobación que no se salta
nadie es la de la Pull Request, descrita más abajo.

### En Windows

Funciona igual. Git for Windows incluye su propio entorno POSIX con `sh`, `grep` y `printf`, y
ejecuta los hooks a través de él, así que un script `#!/bin/sh` corre sin necesidad de WSL.

El repositorio incluye un `.gitattributes` que fuerza finales de línea LF en los scripts. Es
necesario porque el instalador de Git for Windows deja `core.autocrlf = true`, y un script de shell
con finales CRLF no llega ni a arrancar:

```text
bad interpreter: /bin/sh^M: no such file or directory
```

Desconcierta porque el archivo se ve correcto al abrirlo. No hace falta hacer nada: `.gitattributes`
gana sobre la configuración local de cada clon.

## Limpiar las referencias de ramas borradas

```bash
git config --global fetch.prune true
```

Al fusionar una Pull Request, la rama desaparece del servidor, pero tu clon conserva su referencia
de seguimiento (`origin/nombre-de-la-rama`) para siempre. Un `git fetch` normal solo añade y
actualiza referencias; nunca borra ninguna.

Esas referencias no son inofensivas: cada una es un archivo real dentro de `.git/refs/remotes/` y
`.git/logs/refs/remotes/`, y se acumulan con cada rama fusionada.

Con `fetch.prune` activado, cada `fetch` y cada `pull` eliminan solos las referencias cuya rama ya
no existe en el servidor. Puntualmente, lo mismo se consigue con `git fetch --prune`.

## Comprobaciones automáticas

| Qué | Cuándo | Dónde |
| --- | --- | --- |
| Nombre de la rama | Antes de cada `push` | `.githooks/pre-push` (requiere activarlo) |
| Nombre de la rama | En cada Pull Request | `.github/workflows/branch-name.yml` |

Ambas usan el mismo script, `scripts/check-branch-name.sh`, para que no puedan discrepar. Puedes
comprobar un nombre a mano en cualquier momento:

```bash
sh scripts/check-branch-name.sh feat/agregar-validacion
```

La nomenclatura completa está en [CONTRIBUTING.md](../CONTRIBUTING.md), sección 5.

## Formatear el código

Usa siempre:

```bash
go fmt ./...
```

No uses `gofmt -l .`. Recorre el directorio entero, incluido `.git/`, e intenta parsear como Go
cualquier archivo acabado en `.go`, incluidas las referencias de ramas cuyo nombre termina así. El
error resultante es desconcertante, porque señala un archivo que no es código:

```text
expected 'package', found 0000000000000000000000000000000000000000
```

De ahí viene la prohibición de usar puntos en los nombres de rama. Si necesitas el modo «solo
listar, sin modificar», pasa las rutas explícitamente:

```bash
gofmt -l ./cmd ./internal ./pkg
```

`gofmt -l $(git ls-files '*.go')` no sirve como alternativa: falla con `lstat: no such file or
directory` en cuanto hay un archivo borrado y todavía sin confirmar.
