package models

import (
	"go/ast"
	"go/token"
	"go/types"
)

// Package holds information about a Go package
type Package struct {
	name  string
	defs  map[*ast.Ident]types.Object
	files []*File
	fset  *token.FileSet
}

// Name returns the name of the package
func (p Package) Name() string {
	return p.name
}
