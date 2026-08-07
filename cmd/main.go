package main

import (
	"fmt"

	"github.com/heizeisaburou/templatex/internal/document"
)

func main() {
	// mylista := list.New[int, int](10, 20)
	// fmt.Printf("%#v", mylista)
	mydocumento := document.New([]byte("hola"))

	fmt.Printf("%v", mydocumento.Len())
}
