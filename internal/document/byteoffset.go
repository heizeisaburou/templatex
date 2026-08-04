package document

// ByteOffset representa un desplazamiento en bytes dentro de un documento.
type ByteOffset int

func (o ByteOffset) Compare(other ByteOffset) int {
	switch {
	case o < other:
		return -1
	case o > other:
		return 1
	default:
		return 0
	}
}

func (o ByteOffset) Clamp(end ByteOffset) ByteOffset {
	if o < 0 {
		return 0
	}
	if o > end {
		return end
	}

	return o
}
