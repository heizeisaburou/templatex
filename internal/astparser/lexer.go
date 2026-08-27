package astparser

import (
	"github.com/heizeisaburou/templatex/internal/ast"
	"github.com/heizeisaburou/templatex/internal/document"
)

type Lexer struct{}

func (l Lexer) Lexerize(doc document.Document) []ast.Token {
	panic("Termina Lexerize")
}
