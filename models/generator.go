package models

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/format"
	"go/token"
	"log"

	"golang.org/x/tools/go/packages"
	"golang.org/x/tools/imports"

	"github.com/aurelbec/go-fields-updater/templates"
)

// Generator holds the state of the analysis. Primarily used to buffer
// the output for format.Source.
type Generator struct {
	buf bytes.Buffer // Accumulated output.
	pkg *Package     // Package we are scanning.
}

// Printf prints the string to the output
func (g *Generator) Printf(format string, args ...interface{}) {
	_, _ = fmt.Fprintf(&g.buf, format, args...)
}

// ParsePackage analyzes the single package constructed from the patterns and tags.
// ParsePackage exits if there is an error.
func (g *Generator) ParsePackage(patterns []string, tags []string) *Package {
	fset := token.NewFileSet()
	cfg := &packages.Config{
		Mode:  packages.NeedName | packages.NeedTypes | packages.NeedSyntax | packages.NeedTypesInfo,
		Fset:  fset,
		Tests: false,
	}
	pkgs, err := packages.Load(cfg, patterns...)
	if err != nil {
		log.Fatal(err)
	}
	if len(pkgs) != 1 {
		log.Fatalf("error: %d packages found", len(pkgs))
	}
	return g.addPackage(pkgs[0], fset)
}

// addPackage adds a type checked Package and its syntax files to the generator.
func (g *Generator) addPackage(pkg *packages.Package, fset *token.FileSet) *Package {
	g.pkg = &Package{
		fset:  fset,
		name:  pkg.Name,
		defs:  pkg.TypesInfo.Defs,
		files: make([]*File, len(pkg.Syntax)),
	}

	for i, file := range pkg.Syntax {
		g.pkg.files[i] = &File{
			file: file,
			pkg:  g.pkg,
		}
	}
	return g.pkg
}

// Format returns the gofmt-ed contents of the Generator's buffer.
func (g *Generator) Format() []byte {
	quit := func(err error) []byte {
		// Should never happen, but can arise when developing this code.
		// The user can compile the output to see the error.
		log.Printf("warning: internal error: invalid Go generated: %s", err)
		log.Printf("warning: compile the package to analyze the error")
		return g.buf.Bytes()
	}

	src, err := format.Source(g.buf.Bytes())
	if err != nil {
		return quit(err)
	}

	src, err = imports.Process("", src, nil)
	if err != nil {
		return quit(err)
	}

	return src
}

// Generate
func (g *Generator) Generate(typeName string) {
	fields := make([]Field, 0, 20)
	for _, file := range g.pkg.files {
		// Set the state for this run of the walker.
		if file.file != nil {
			file.reset(typeName)
			ast.Inspect(file.file, file.genDecl)
			fields = append(fields, file.fields...)
		}
	}

	if len(fields) == 0 {
		log.Fatalf("no fields defined for %s", typeName)
	}

	list := make([][2]string, len(fields))
	for i, field := range fields {
		// log.Print("found ", field)
		list[i] = [2]string{field.fieldName, field.fieldType}
	}
	g.Printf(templates.GenerateUpdateFromFields(typeName, list))
}
