package main

import "fmt"

func printArray(label string, elements ...int) {
	fmt.Printf("%s = %v | len = %d | cap = %d\n", label, elements, len(elements), cap(elements))
}

func main() {
	// Crea un slice de enteros con los valores 10 y 20.
	var miArray = []int{10, 20}

	// Imprime sus elementos, longitud y capacidad.
	printArray("miArray", miArray...)

	// Declara un slice de enteros nil.
	var miArray2 []int

	// Imprime si es nil, su longitud y su capacidad.
	fmt.Printf("miArray2 == nil? %t | len %d | cap %d\n", miArray2 == nil, len(miArray2), cap(miArray2))

	// Reserva capacidad para 2 elementos, manteniendo longitud 0.
	var miArray3 = make([]int, 0, 2)

	// Añade 1, 2 y 3 de una sola vez.
	miArray3 = append(miArray3, 1, 2, 3)

	// Imprime sus elementos, longitud y capacidad.
	printArray("miArray3", miArray3...)

	// Añade 4 y 5.
	miArray3 = append(miArray3, 4, 5)

	// Imprime de nuevo sus elementos, longitud y capacidad.
	printArray("miArray3", miArray3...)

	// Modifica el segundo elemento usando su índice.
	miArray3[2] = 69

	// Copia el slice en otra variable.
	miArray4 := append([]int(nil), miArray3...)

	// Modifica un elemento desde la nueva variable.
	miArray4[3] = 70

	// Comprueba si el slice original también ha cambiado.
	printArray("miArray3", miArray3...)
	printArray("miArray4", miArray4...)

	// Crea un subslice desde el índice 1 hasta el 4, sin incluir el 4.
	miArray5 := miArray4[1:4]

	// Imprime sus elementos, longitud y capacidad.
	printArray("miArray5", miArray5...)

	// Crea otro subslice con tres índices y capacidad limitada a 3.
	low := 1
	high := 4

	desiredCapacity := 3

	max := low + desiredCapacity

	miArray6 := miArray4[low:high:max]

	// Añade un elemento al subslice limitado.
	miArray6 = append(miArray6, 128)

	// Comprueba si el append reutilizó el array original o creó uno nuevo.
	printArray("miArray4", miArray4...)
	printArray("miArray6", miArray6...)
}
