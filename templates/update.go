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

type %[1]sField int

type %[1]sFields map[%[1]sField]interface{}

const (
%[2]s
)

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
		enumStr := fmt.Sprintf(`%sField%s`, name, field[0])
		caseStr := fmt.Sprintf(`case %s:`, enumStr)

		enums = append(enums,
			enumStr,
		)

		checkCases = append(checkCases,
			caseStr,
			fmt.Sprintf(`if _, ok := v.(%[2]s); !ok {
			                err = append(err, fmt.Sprintf("value for %[1]s is not %%%%s (got %%%%T)", reflect.TypeOf(&t.%[1]s).Elem(), v))
			             }`, field[0], field[1]),
		)

		updateCases = append(updateCases,
			caseStr,
			fmt.Sprintf(`t.%s = v.(%s)`, field[0], field[1]),
		)
	}

	if len(enums) > 0 {
		enums[0] = fmt.Sprintf(`%s %sField = iota`, enums[0], name)
	}

	return fmt.Sprintf(fieldsTypeFormat, name, strings.Join(enums, "\n"), strings.Join(checkCases, "\n"), strings.Join(updateCases, "\n"))
}
