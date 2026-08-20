package ast

import (
	"github.com/heizeisaburou/templatex/internal/document"
)

type TokenFormatter struct {
	dc document.Document
}

// func (f TokenFormatter) String(t Token) string {
// 	return fmt.Sprintf("Slice(%v),\n Kind(%d), \n Range(%v)",
// 		f.dc.Slice(t.rng), t.kind, t.rng)
// }

// Inventar una forma de mostrar un Token modo debug
