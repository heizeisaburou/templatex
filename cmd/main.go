package main

import (
	"fmt"

	"github.com/heizeisaburou/templatex/internal/document"
	"github.com/heizeisaburou/templatex/pkg/list"
)

func main() {
	mylista := list.New[int, int](10, 20)
	fmt.Printf("%#v", mylista)

	src := []byte("este es mi archivo")
	myDocumento := document.New(src)

	_ = myDocumento

}
