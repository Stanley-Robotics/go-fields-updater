package models

import (
	"fmt"
)

// Field represents a single field in a struct.
type Field struct {
	structName     string // The name of the struct to which this field belongs
	fieldName      string // The name of the field
	fieldType      string // The type of the field
	fieldUpdatable bool   // Whether the field is updatable, or not
}

func (f Field) String() string {
	return fmt.Sprintf("%s.%s %s", f.structName, f.fieldName, f.fieldType)
}
