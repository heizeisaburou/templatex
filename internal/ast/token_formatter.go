package ast

import (
	"fmt"
	"strings"

	"github.com/heizeisaburou/templatex/internal/document"
)

type TokenFormatter struct {
	doc document.Document
}

func NewTokenFormatter(doc document.Document) TokenFormatter {
	return TokenFormatter{doc: doc}
}

func (f TokenFormatter) Sprint(t Token) (string, error) {
	var b strings.Builder

	b.WriteString("Token{kind=")
	fmt.Fprint(&b, t.Kind())
	b.WriteString(", rng=")
	fmt.Fprint(&b, t.Range())
	b.WriteString(", src=")
	src, err := f.doc.Range(t.rng)
	if err != nil {
		return "", err
	}
	fmt.Fprintf(&b, "%q", string(src))
	b.WriteString("}")

	return b.String(), nil
}
