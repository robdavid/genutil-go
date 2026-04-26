/*
The opt package is intended to bring some null safety to Go. The idiomatic
approach to values that might be missing (i.e. nullable values) is to use a
pointer to that value, with a nil pointer indicating the value is absent. This
tends to overload two concerns, i.e. using pointers to implement access by
reference, and also to indicate the possibility of missing values. This package
tries to separate these concerns.

The core concepts are two concrete types, [opt.Val][T] and [opt.Ref][T], both
implementing the [Option] interface:

  - [opt.Val][T]: A value wrapper that holds a concrete value of type T, and a boolean
    flag to indicate a value is present (non-empty). Best suited for simple,
    non-reference data types like int, string, etc.
  - [opt.Ref][T]: A reference wrapper that holds a pointer to a value of type T, which
    when nil indicates the value is not present (empty). Used when the underlying
    type is expensive to copy or access by reference is desired, e.g. for mutability.

The package provides utility functions such as [opt.Value](v) and
[opt.Empty][T]() to create instances of [opt.Val][T], and opt.Reference(&v) and
[opt.EmptyRef][T]() to create instances of [opt.Ref][T].

# Usage

To use the option types, first determine if you need a simple value wrapper
([opt.Val][T]) or a reference wrapper ([opt.Ref][T]), and then use the
appropriate factory function:

	func Example() {
	    // Simple value (e.g., int)
	    valInt := opt.Value(123)   // Non-empty opt.Val[int]
	    var emptyInt opt.Option[int] = opt.Empty[int]() // Empty opt.Val[int]

	    // Referenced struct (requires address of the struct)
	    type MyStruct struct { Name string; Count int }
	    myInstance := MyStruct{}
	    structRef := opt.Reference(&myInstance) // Non-empty opt.Ref[MyStruct]
	    var emptyStructRef opt.Ref[MyStruct] = opt.EmptyRef[MyStruct]() // Empty opt.Ref[MyStruct]

	}

# Accessing the Wrapped Value

There are a number of methods by which the underlying value may be accessed,
like [opt.Option.Get], [opt.Option.Ref], [opt.Option.GetOK] or [opt.Option.Try]
and others, and also methods to check if the value is there to be accessed, like
[opt.Option.IsEmpty] or [opt.Option.HasValue]. The general concept in each case
is that in order to access the value, it must be accessed through one or other
of these methods. This hopefully forces the programmer to give some thought as
to how to handle the absent value case.

# Zero Value

The zero value for any option type ([opt.Val][T] or [opt.Ref][T]) is always
empty, indicating no stored value was provided.

	func Example() {
	    var zeroOpt opt.Val[int] // If Option is used as an alias for Val/Ref, this will be empty
	    fmt.Println(zeroOpt.IsEmpty()) // true
	}

# Comparisons

Two [opt.Val][T] of the same underlying type can be compared successfully using
==, provided the contained types also support comparison. An empty option will
always compare as unequal to a non-empty one. Comparison with == on a [opt.Ref]
instance will be comparing pointers rather than values.

Two additional functions, [opt.Equal] and [opt.DeepEqual] are provided for value
based comparison across both [opt.Val] and [opt.Ref] instances (or a mix of
each).

# Mutations

The package provides methods like [opt.Ensure] and [opt.Mutate] which allow for
in-place mutation, especially useful when wrapping large structs by reference
([opt.Ref][T]). The [opt.Ensure] method will ensure the Option is non-empty by
placing a zero value in the Option if it currently empty. The [opt.Mutate]
method passes a pointer to the value to a user provided function that may then
modify it.

	func Example() {
	    type nv struct{ name, value string }
	    optEmpty := opt.Empty[nv]()
	    // Mutates the underlying data structure in place.
	    optEmpty.Ensure().Mutate(func(n *nv) {
	        n.name = "new_name"
	        n.value = "new_value"
	    })
	}

# Marshalling and Unmarshaling

Option types implement JSON and YAML marshaling interfaces. This allows them to
be included in standard Go structures that are serialized or deserialized by
external packages. Note that this behavior is consistent for both [opt.Val][T]
and [opt.Ref][T].

Marshalling: A non-empty option will be marshalled as its contained value in
both JSON and YAML (assuming omitempty tag usage).

	type testOptMarshall struct {
	    Name  opt.Val[string] `json:"name,omitempty" yaml:"name,omitempty"`
	    Value opt.Val[int]    `json:"value,omitempty" yaml:"value,omitempty"`
	}

	func Example() {
	    testData := testOptMarshall{
	        Name:   opt.Value("a name"),
	        Value:  opt.Value(123),
	    }
	    // Marshal this struct into JSON/YAML bytes...
	}

Empty options will be rendered as null in JSON, and omitted entirely if the
omitempty tag is present when marshalling YAML. For the more recent JSON v2
package in the standard library, omitzero can be used to omit empty values
entirely from JSON output.

Unmarshaling: When unmarshaling data (JSON or YAML) to option types, keys that
are missing or explicitly set to null will result in an empty option value.

	func Example() {
	    var unmarshalledData testOptMarshall
	    // Unmarshal YAML containing null/missing fields into testOptMarshall
	    yaml.Unmarshal([]byte("name: a name\nvalue: null"), &unmarshalledData)

	    // Check the status of the fields:
	    fmt.Println(unmarshalledData.Name.HasValue())  // true (was present)
	    fmt.Println(unmarshalledData.Value.HasValue()) // false (was null/missing)
	}

YAML parser v3: Support for "gopkg.in/yaml.v2" is available in the standard opt
package, and is implemented without any explicit dependency on that library.
Support for the more recent "gopkg.in/yaml.3" package is also available, but via
the opt/yamlv3 package which pulls in the YAML library as a dependency. To use
this version of opt use the following import statement:

	import opt "github.com/robdavid/genutil-go/opt/yamlv3"
*/
package opt
