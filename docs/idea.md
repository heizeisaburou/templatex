# Idea original

## Brief

En este documento describimos lo que padawan y yo pensamos en un principio para el lenguaje. Lo
que debería tener el primer MVP y cómo se vería.

## Sintaxis de ejemplo

### Mostrar campo

```text
<% fld fieldName:arg1:arg2 %>
```

### Comentarios

```text
<!Aqui mostrariamos comentarios!>

<! Comentario
   Multilinea
!>
```

### Condicionales

```text
<% if ! (palabra || significado) %>
<@ import archivo.ts @>
<% elseif lectura || sonido %>

<% else %>

<% end % >
```

### Imports

```text
<@ import archivo.ts @> <! se importaría compilado como javascript !>
<@ import archivo.ts @>
```
