# GUÍA DE CONTRIBUCIÓN

<img src="" alt="" width="" height="">

Esta guía de contribución contiene las específicaciones a manejar con GIT. Incluye a detalle el
flujo de trabajo que se maneja, la función de las ramas y uso de comandos.

Además, se soporta aquí la guía para manejar la documentación del proyecto.

## Índice

<details>
<summary>
1. [Configuración de cuenta de GIT](#configuración-de-cuenta-de-git)
</summary>
</details>

<details>
<summary>
2. [Ramas del proyecto](ramas-del-proyecto)
</summary>
</details>

<details>
<summary>
3. [Flujo de trabajo diario](flujo-de-trabajo-diario)
</summary>
</details>

<details>
<summary>
4. [Resolver un conflicto de merge](resolver-un-conflicto-de-merge)
</summary>
</details>

<details>
<summary>
5. [Nomenclatura de commits y ramas](nomenclatura-de-commits-y-ramas)
</summary>
</details>

<details>
<summary>
6. [Documentación del proyecto](documentación-del-proyecto)
</summary>
</details>

<details open>
<summary> 
# Configuración de cuenta de GIT
</summary>
Tener su cuenta de git es esencial, **ESTÁ** prohibido cambiar el correo o usuario
de su cuenta global de git para cualquiera de los espacios de trabajo (laptop, desktop PC, 
VM, Sandbox, directory, workspace, enviroment) desde
donde vaya a trabajar con el repositorio.

Para evitar estos problemas, se recomienda manejar configuración local, es decir, que tenga su
cuenta configurada por directorios/espacios de trabajo.

> [!IMPORTANT]
>
> Sí configuras una cuenta de git a nivel local, esta tendra prioridad sobre configuración
> globales o de sistema (`local > global > system`).

<details>
<summary>
## Configurar cuenta de git
</summary>
1. Configuración local (recomendada)
   Usuario:
   ~~~
   git config user.name "kebab-case"
   ~~~
   Correo:
   ~~~
   git config user.email "user@gmail.com"
   ~~~
2. Configuración global
   Usuario:
   ~~~
   git config --global user.name "kebab-case"
   ~~~
   Correo:
   ~~~
   git config --global user.email "user@gmail.com"
   ~~~
</details>

<details open>
<summary>
# Ramas del proyecto
</summary>
Para garantizar la integridad del proyecto, se han creado dos ramas de trabajo, a las que se suman las ramas de los contribuidores (5 en total).
Trabajar con varias ramas en git suena a una odísea, sin embargo, con el debido cuidado y seguimiento de las normas establecidas, garantizará que
si se llega a aceptar un pull request innecesario, exista una rama previa a main que aún se pueda tratar.

```mermaid
flowchart TB

	%% ORIGEN - de las ramas dev
	M((MAIN))
	O((STAGE))
	M --> O

	%% MAIN — arriba
	subgraph MAIN["main"]
		direction LR
		M1((●)) --> M2((●)) --> M3((●))
	end

	%% STAGE — segunda capa
	subgraph STAGE["stage"]
		direction LR
		S1((●)) --> S2((●)) --> S3((●))
	end

	%% DEVELOPERS — abajo
	subgraph DEVELOPERS["developers"]
		direction TB

		subgraph DEV1["dev-1"]
			direction LR
			D1_1((●)) --> D1_2((●)) --> D1_3((●)) --> D1_4((●))
		end

		subgraph DEV2["dev-2"]
			direction LR
			D2_1((●)) --> D2_2((●)) --> D2_3((●)) --> D2_4((●))
		end

		subgraph DEV3["dev-3"]
			direction LR
			D3_1((●)) --> D3_2((●)) --> D3_3((●)) --> D3_4((●))
		end
	end

	%% ÚNICO ORIGEN
	O --> D1_1
	O --> D2_1
	O --> D3_1

	%% DEV → STAGE
	D1_3 --> S1
	D2_2 --> S2
	D3_3 --> S3

	%% STAGE → MAIN
	S2 --> M1
	S3 --> M2
```

La rama de Stage sirve de puente entre las ramas dev y main. El objetivo de su existencia es
tener una rama previa en el repositorio sin ser la principal, permitiendo evaluar a todos los
usuarios el contenido subido aquí. La particularidad de esta rama es que se requiere hacer Pull
Request para combinar un commit de dev con la rama.

Una vez que se haya aprovado o desaprovado el commit en stage por todos los miembros, está puede
ser subida a main con un Pull Request, y aquí el dueño del repositorio tiene la responsabilidad
de resolver conflictos de merge si se presentan.

[!WARNING]

> Está prohibido enviar un pull request directo a main o hacer push si es el dueño del
> repositorio, y en caso de que suceda esto, se sancionará al responsable.

Las ramas dev se manejan de manera sencilla, y su flujo será explicado en la sección a
continuación, cabe resaltar que únicamente se puede tener una única rama dev a la vez, ya que si
se tienen más, manejar la sincronización con stage puede ser fatal.
</details>

<details>
<summary>
# Flujo de trabajo diario
</summary>
Los merges conflict son demasiado comunes, estresantes y justos. Sin embargo no son necesarios, siempre y cuando se maneja un correcto flujo de
trabajo, que será explicado a continuación:

1. Actualizar main local respecto al main y stage de github.
   ~~~
   git pull origin main
   git pull origin stage
   ~~~
2. Crear una rama y moverse a ella, el nombre va en kebab-case en base al problema que van a resolver y siguiendo la nomenclatura.
   ~~~
   git checkout -b "categoría/ propósito"
   ~~~
3. Luego de trabajar con ella, se pasa al área de staging y de ahí commitea.
   ~~~
   git add .
   ~~~
   git commit -m "categoría/ problema tratado"
3. Procede a subir su rama al repositorio
   ~~~
   git push -u origin [rama]
   ~~~
4. Desde github, hace un pull request, evita que tenga varios commits, para reducir tiempos de redición.
   ![el gatoso](docs/assets/github_pull_request.png)
5. Una vez en stage, tras la aprobación de todos, se resuelven merge conflicts (si existen).

6. El que envió el pull request a stage, ahora hace uno a main, el cuál será aceptado por el dueño.

7. Elimina la rama creada moviendose a stage (nunca estar sobre main)
   ~~~
   git switch stage
   git branch -d [rama]
   ~~~
</details>

<details>
<summary>
# Nomenclatura de commits y ramas
</summary>
La revisión rápida de commits 
</details>
