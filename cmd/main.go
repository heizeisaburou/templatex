package main

import (
	"container/list"
	"fmt"
)

func main() {
	mylista := list.New[int, int](10, 20)
	fmt.Printf("%#v", mylista)

}
