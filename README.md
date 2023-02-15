# Go-Fields-Updater [![GitHub Release](https://img.shields.io/github/v/release/stanley-robotics/go-fields-updater)](https://github.com/Stanley-Robotics/go-fields-updater/releases)


Go-Fields-Updater is a tool to generate Go code that adds useful methods to update specific public structure members.


```
$ ./go-fields-updater --help
Go-Fields-Updater is a tool to generate Go code that adds useful method to update specific fields of specific type.
Usage of go-fields-updater:
        go-fields-updater [flags] -type T [directory]
        go-fields-updater [flags] -type T files... # Must be a single package
For more information, see:
        https://github.com/stanley-robotics/go-fields-updater
Flags:
  -output string
        output file name; default srcdir/<type>_updater.go
  -type string
        comma-separated list of type/structures names; must be set
```


## Generated functions and methods

When Go-Fields-Updater is applied to a structure, it will generate:

- The following basic methods:
  - Method `GetFieldsValues()`: returns a map representation of the current public members value.
  - Methods `UpdateField` and `UpdateFields`: update one or multiple fields with the specified values.

- Useful types to easily select fields to update:
  - Type `<Type>Field`: describes a public structure field as a `string`.
  - Type `<Type>Fields`: describes a map of `<Type>Field` with their respective value.
    This type comes with the following methods:
      - Method `Merge`: merges an other `<Type>Fields` into the current one.
      - Method `Contains`: returns whether or not any of the specified `<Type>Field` is part of the current `<Type>Fields`.
      - Method `Fields`: returns the ordered list of `<Type>Field` currently present in the map.

For example, if we have an structure type called `Person`,

```go
type gender int

const (
	Male gender = iota
	Female
)

type Person struct {
	Name      string
    Gender    gender
    age       int
    Relatives []Person
}
```

executing `go-fields-updater -type=Person` will generate a new file (`person_updater.go` by default) with basic methods:

```go
type PersonField string

type PersonFields map[PersonField]interface{}

// Merge is a convenient method to add elements from other to this map
func (this PersonFields) Merge(other PersonFields) {
	//...
}

// Contains returns whether the map contains at least one of the specified keys
func (this PersonFields) Contains(keys ...PersonField) bool {
	//...
}

// Fields returns the sorted list of fields composing the map
func (this PersonFields) Fields() []string {
	//...
}

const (
	PersonFieldName      PersonField = "Name"
	PersonFieldGender    PersonField = "Gender"
	PersonFieldRelatives PersonField = "Relatives"
)

// GetFieldsValues returns the Person struct value as PersonFields representation.
func (t Person) GetFieldsValues() PersonFields {
	//...
}

// UpdateField updates the specified field of the Person struct.
func (t *Person) UpdateField(field PersonField, value interface{}) error {
	//...
}

// UpdateFields updates the fields of the Person struct from the given fields map.
func (t *Person) UpdateFields(fields PersonFields) error {
	//...
}
```

From now on, we can:

```go
// Create an empty Person
var me Person = Person{}

// Update its Name
me.UpdateField(PersonFieldName, "Bob")
fmt.Println("My name is ", me.Name) // Will print "My name is Bob"

// Update both Name and Gender
me.UpdateFields(PersonFields{PersonFieldName: "Bob", PersonFieldGender: Male, })
fmt.Println("My name is ", me.Name, " (gender: ", me.Gender, ")") // Will print "My name is Bob (gender: 0)"

// Detect invalid update
if (err := me.UpdateField(PersonFieldGender, "Male"); err != nil ) {
    fmt.Println(err) // Will print "value for Gender is not main.gender (got string)"
}

// Get a PersonFields representation
fmt.Println(me.GetFieldsValues()) // Will print "map[Gender:0 Name:Bob Relatives:[]]"
```

Note that it does not allow updates on private fields.


## How to use

For a module-aware repo with `go-fields-updater` in the `go.mod` file, generation can be called by adding the following to a `.go` source file:

```golang
//go:generate go run go-fields-updater -type=YOURTYPE
```


## Inspiring projects

- [Dan Markham](https://github.com/dmarkham/enumer)
- [Stringer](https://godoc.org/golang.org/x/tools/cmd/stringer)
