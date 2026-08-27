package document

import "fmt"

// Position representa una posición en formato humano dentro de un documento.
type Position struct {
	line int
	col  int
}

// ByteOffset
// En la linea X columna X pasa X

func NewPosition(line int, col int) Position {
	if line < 0 || col < 0 {
		panic(fmt.Sprintf("invalid position: %d:%d", line, col))
	}

	return Position{
		line: line,
		col:  col,
	}
}

// Compare primero retorna si las lineas no son iguales y cual es mayor o menor.
//
// y si el primer test pasa retorna cual columna es mayor a la otra o si son iguales.
func (p Position) Compare(other Position) int {
	if p.line != other.line {
		if p.line < other.line {
			return -1
		}
		return 1
	}

	if p.col < other.col {
		return -1
	}
	if p.col > other.col {
		return 1
	}

	return 0
}

func (p Position) Line() int {
	return p.line
}

func (p Position) Col() int {
	return p.col
}

func (p *Position) SetLine(value int) {
	if value < 0 {
		panic(fmt.Sprintf("minimum value = 0; got %d", value))
	}

	p.line = value
}
func (p *Position) SetCol(value int) {
	if value < 0 {
		panic(fmt.Sprintf("minimum value = 0; got %d", value))
	}

	p.col = value
}

func (p Position) String() string {
	return fmt.Sprintf("%d:%d", p.line, p.col)
}
