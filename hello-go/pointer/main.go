package main

import "fmt"

type person struct {
	firstName string
	lastName  string
}

func main() {
	// Let us create an int value and get it's address
	x := 42
	px := &x
	// Note the *dereferencing* with the * operator in the next line
	fmt.Printf("x is at address %v and it's value is %v\n", px, *px)

	// Let us double the value in x. Again, we use dereferencing
	*px *= 2
	fmt.Printf("x is at address %v and it's value is %v\n", px, *px)

	// We can allocate memory and retrieve a point to it using *new*.
	// The allocated memory is automatically set to zero.
	px = new(int)
	fmt.Printf("x is at address %v and it's value is %v\n", px, *px)

	// Go only knows call-by-value. In the following case, the value
	// is a pointer and the method can dereference it to write something
	// into the memory the pointer points to.
	func(val *int) {
		*val = 42
	}(px)
	fmt.Printf("x is at address %v and it's value is %v\n", px, *px)

	// We can also create pointers to structs. Note that although we
	// have a pointer, we can still access the struct's members using
	// a dot.
	pp := &person{"Foo", "Bar"}
	fmt.Printf("%s, %s\n", pp.lastName, pp.firstName)

	// We can pass a pointer to a struct to a method. In this case,
	// the method can change the struct's content.
	func(somebody *person) {
		somebody.firstName, somebody.lastName = somebody.lastName, somebody.firstName
	}(pp)
	fmt.Printf("%s, %s\n", pp.lastName, pp.firstName)
}
