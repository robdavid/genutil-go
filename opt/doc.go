/*
The opt package is attempting to bring some nil safety to Go. Traditionally, the
idiomatic approach to values that might be missing has been to use a pointer to
a value, with a nil pointer indicating that there is no value available. This
tends to overload two concerns onto one feature, i.e. using pointers to
implement call by reference, such as for when passing a large structure by value
is expensive, and also to indicate the possibility of missing values. This
package tries to separate these concerns.

This package provides generic types that can be used to wrap either values or
references to values that can be empty, or not present. It tries to do this in a
way that allows the programmer to more safely handle nullability and reduce the
risk of inadvertently triggering a panic.

# Usage

The package provides the following types:

  - [Option][T] An interface that defines an abstraction over the two concrete types.
  - [Val][T]    A type that holds a value of type T, or holds nothing.
  - [Ref][T]    A type that holds a reference to value of type T, or holds nothing.

Typically, you would use an instance of Val to hold a simple type, such as an
`int`, that may be nullable. You can use the [Value] or [Empty] functions to
create instances of [Val]:

	int_value := opt.Value(123)   # Creates a Val[Int] containing 123
	int_empty := opt.Empty[int]() # Creates a Val[int] with no value

If you have a larger structure, then you might want to access it by reference.
For this use the [Reference] or [EmptyRef] functions to create instance of
[Ref].

	var myInstance MyStruct = MyStruct {}
	struct_ref   := opt.Reference(&myInstance) # Creates a Ref[MyStruct] pointing to myInstance
	struct_empty := opt.EmptyRef[MyStruct]()   # Creates a Ref[MyStruct] with no value.

Both [Val] and [Ref] instances contain the same collection of methods (defined
by [Option]) for checking for the presence of, and for accessing, the underlying
value. Note that by requiring access to the value via methods, the possibility
of inadvertently accessing a null member is reduced.

E.g.

	func addOne(optVal opt.Val[int]) opt.Val[int] {
		if optVal.HasValue() {
			return opt.Value(optVal.Get()+1)
		} else {
		 	return optVal
		}
	}

or more succinctly

	func addOne(optVal opt.Val[int]) opt.Val[int] {
	    return optVal.Morph(func(x int) { return x+1 }
	}
*/

package opt
