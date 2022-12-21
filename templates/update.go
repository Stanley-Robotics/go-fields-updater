package templates

import (
	"fmt"
	"strings"
)

// Arguments to format are:
//	[1]: type name
//	[2]: fields enums
//	[3]: check cases
//	[4]: update cases
const fieldsTypeFormat = `

type %[1]sField string

type %[1]sFields map[%[1]sField]interface{}

// Merge is a convenient method to add elements from other to this map
func (this %[1]sFields) Merge(other %[1]sFields) {
	for field, value := range other {
		this[field] = value
	}
}

// Contains returns whether the map contains at least one of the specified keys
func (this %[1]sFields) Contains(keys ...%[1]sField) bool {
	for _, key := range keys {
		if _, contains := this[key]; contains {
			return true
		}
	}
	return false
}

// Fields returns the sorted list of fields composing the map
func (this %[1]sFields) Fields() []string {
	fields := make(sort.StringSlice, 0, len(this))
	for field := range this {
		fields = append(fields, string(field))
	}
	fields.Sort()
	return fields
}

const (
%[2]s
)

// GetFieldsValues returns the %[1]s struct value as %[1]sFields representation.
func (t %[1]s) GetFieldsValues() %[1]sFields {
	return %[1]sFields {
		%[3]s
	}
}

// UpdateField updates the specified field of the %[1]s struct.
func (t *%[1]s) UpdateField(field %[1]sField, value interface{}) error {
	return t.UpdateFields(%[1]sFields{field: value})
}

// UpdateFields updates the fields of the %[1]s struct from the given fields map.
func (t *%[1]s) UpdateFields(fields %[1]sFields) error {
	// ensure consistency of fields first
	err := make([]string, 0, len(fields))
	for k, v := range fields {
		switch k {
		%[4]s
		}
	}

	if len(err) > 0 {
		return errors.New(strings.Join(err, ", "))
	}

	// proceed with updating fields
	for k, v := range fields {
		switch k {
		%[5]s
		}
	}

	return nil
}
`

// GenerateFields generates the fields enum and the updateFromFields function.
func GenerateUpdateFromFields(name string, fields [][2]string) string {
	var enums []string
	var values []string
	var checkCases []string
	var updateCases []string
	for _, field := range fields {
		enumStr := fmt.Sprintf(`%[1]sField%[2]s %[1]sField = "%[2]s"`, name, field[0])
		caseStr := fmt.Sprintf(`case %[1]sField%[2]s:`, name, field[0])

		enums = append(enums,
			enumStr,
		)

		values = append(values,
			fmt.Sprintf("%[1]sField%[2]s: t.%[2]s,", name, field[0]),
		)

		checkCases = append(checkCases,
			caseStr,
			fmt.Sprintf(`if _, ok := v.(%[2]s); !ok && v != nil {
			                err = append(err, fmt.Sprintf("value for %[1]s is not %%%%s (got %%%%T)", reflect.TypeOf(&t.%[1]s).Elem(), v))
			             }`, field[0], field[1]),
		)

		updateCases = append(updateCases,
			caseStr,
			fmt.Sprintf(`t.%s, _ = v.(%s)`, field[0], field[1]),
		)
	}

	return fmt.Sprintf(fieldsTypeFormat, name, strings.Join(enums, "\n"), strings.Join(values, "\n"), strings.Join(checkCases, "\n"), strings.Join(updateCases, "\n"))
}
