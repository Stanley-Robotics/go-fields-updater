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

type %[1]sFields int

type %[1]sFieldUpdates map[%[1]sFields]interface{}

const (
%[2]s
)

// UpdateFromFields updates the fields of the %[1]s struct from the given fields map.
func (t *%[1]s) UpdateFromFields(fields %[1]sFieldUpdates) error {
	// ensure consistency of fields first
	err := make([]string, 0, len(fields))
	for k, v := range fields {
		switch k {
		%[3]s
		}
	}

	if len(err) > 0 {
		return errors.New(strings.Join(err, ", "))
	}

	// proceed with updating fields
	for k, v := range fields {
		switch k {
		%[4]s
		}
	}

	return nil
}
`

// GenerateFields generates the fields enum and the updateFromFields function.
func GenerateUpdateFromFields(name string, fields [][2]string) string {
	var enums []string
	var checkCases []string
	var updateCases []string
	for _, field := range fields {
		enum := fmt.Sprintf(`%sField%s`, name, field[0])
		enums = append(enums, enum)

		checkCases = append(checkCases, fmt.Sprintf(`case %s:`, enum))
		checkCases = append(checkCases, fmt.Sprintf(`if _, ok := v.(%s); !ok {err = append(err, fmt.Sprint("invalid type for %s: ", reflect.TypeOf(v), " != ", reflect.TypeOf(t.%s)))}`, field[1], field[0], field[0]))

		updateCases = append(updateCases, fmt.Sprintf(`case %s:`, enum))
		updateCases = append(updateCases, fmt.Sprintf(`t.%s = v.(%s)`, field[0], field[1]))
	}

	if len(enums) > 0 {
		enums[0] = fmt.Sprintf(`%s %sFields = iota`, enums[0], name)
	}

	return fmt.Sprintf(fieldsTypeFormat, name, strings.Join(enums, "\n"), strings.Join(checkCases, "\n"), strings.Join(updateCases, "\n"))
}
