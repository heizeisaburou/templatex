package document

import "errors"

var (
	ErrInvalidUTF8                 = errors.New("invalid utf8 sequence")
	ErrOutOfBounds                 = errors.New("out of bounds")
	ErrByteOffsetNotAtRuneBoundary = errors.New("el offset de bytes no apunta al inicio de un carácter codificado en UTF-8")
)
