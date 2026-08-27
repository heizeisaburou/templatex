package ast

import (
	"errors"
	"fmt"

	"github.com/heizeisaburou/templatex/internal/document"
)

type TokenKind int

const (
	TokenUnknown TokenKind = iota
	TokenWhitespace
	tokenLimit
)

func (k TokenKind) String() string {
	switch k {
	case TokenWhitespace:
		return "TokenWhitespace"
	default:
		return fmt.Sprintf("TokenUnknown(%d)", k)
	}
}

var (
	ErrNotValidToken = errors.New("El token utilizado no es valido")
)

type Token struct {
	kind TokenKind
	rng  document.Range
}

func NewToken(kind TokenKind, rng document.Range) (Token, error) {

	if kind < TokenUnknown || kind > tokenLimit {
		return Token{}, ErrNotValidToken
	}

	return Token{kind: kind, rng: rng}, nil
}

func (t Token) Range() document.Range {
	return t.rng
}

func (t Token) Kind() TokenKind {
	return t.kind
}
