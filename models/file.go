package models

import (
	"go/ast"
	"go/token"
	"log"
	"os"
)

// File holds a single parsed file and associated data.
type File struct {
	pkg  *Package  // Package to which this file belongs.
	file *ast.File // Parsed AST.

	// These fields are reset for each type being generated.
	typeName  string // Name of the constant type.
	typeFound bool
	fields    []Field // Accumulator for constant values of that type.
}

// reset sets the state for processing a new type.
func (f *File) reset(typeName string) {
	f.typeName = typeName
	f.typeFound = false
	f.fields = make([]Field, 0, 10)
}

// genDecl processes one declaration clause.
func (f *File) genDecl(node ast.Node) bool {
	// return early
	if f.typeFound {
		return true
	}

	decl, ok := node.(*ast.GenDecl)
	// ignore other nodes than type declarations
	if !ok || decl.Tok != token.TYPE {
		return true
	}

	for _, spec := range decl.Specs {
		// skip other declarations
		tspec, ok := spec.(*ast.TypeSpec)
		if !ok {
			continue
		} else if tspec.Name.Name != f.typeName {
			continue
		}

		// ensure that the type is a struct
		strct, ok := tspec.Type.(*ast.StructType)
		if !ok {
			log.Fatalf("type %s is not a struct", f.typeName)
			continue
		}

		f.typeFound = true
		if strct.Fields == nil { // Empty struct
			break
		}

		data, err := os.ReadFile(f.pkg.fset.Position(spec.Pos()).Filename)
		if err != nil {
			log.Fatalf("error reading file: %s", err)
		}

		// iterate over struct fields
		for _, field := range strct.Fields.List {
			if field == nil {
				continue
			}

			if len(field.Names) == 0 { // Embedded field
				// TODO: handle embedded fields
				continue
			}

			typePos := f.pkg.fset.Position(field.Type.Pos())
			fieldType := string(data[typePos.Offset : typePos.Offset+int(field.Type.End()-field.Type.Pos())])

			for _, name := range field.Names {
				// ignore unexported fields
				if name == nil {
					continue
				} else if !name.IsExported() {
					continue
				}

				f.fields = append(f.fields, Field{
					structName: f.typeName,
					fieldName:  name.Name,
					fieldType:  fieldType,
				})
			}
		}
		break
	}

	return f.typeFound
}
