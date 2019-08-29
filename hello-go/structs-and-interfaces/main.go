package main

import (
	"fmt"
	"math"
)

// Let us define two simple structs. In contrast to C#, there are no classes.
// Everything is a struct. Note that this does not say anything about allocation
// on stack or heap. The compiler will decide for you whether to allocate an
// instance of a struct on the stack or heap based on *escape analysis*
// See also https://en.wikipedia.org/wiki/Escape_analysis.

// ...and yes, dear C# developer, you can believe your eyes:
//    No semicolons at the end of lines in Go ;-)

type Point struct {
	X, Y float64
}

type Rect struct {
	LeftUpper, RightLower Point
}

type Circle struct {
	Center Point
	Radius float64
}

// Now let us add some functions to our structs

func (r Rect) Width() float64 {
	return r.RightLower.X - r.LeftUpper.X
}

func (r Rect) Height() float64 {
	return r.RightLower.Y - r.LeftUpper.Y
}

func (r Rect) Area() float64 {
	return float64(r.Width() * r.Height())
}

func (r *Rect) Enlarge(factor float64) {
	// Note that this function has a *pointer receiver type*.
	// That means that it can manipulate the content of r and
	// the caller will see the changed values. The other methods
	// of rect have a *value receiver type*, i.e. the struct
	// is copied and changes to it's values are not visible
	// to the caller.

	r.RightLower.X = r.LeftUpper.X + r.Width()*factor
	r.RightLower.Y = r.LeftUpper.Y + r.Height()*factor
}

func (c Circle) Area() float64 {
	return math.Pi * c.Radius * c.Radius
}

// Next, we define an interface. Note that the structs do not need to
// explicitly implement the interface. rect and circle fulfill the
// requirements of the interface (i.e. they have an area() method),
// therefore the implement the interface.

type Shape interface {
	Area() float64
}

const (
	WHITE int = 0xFFFFFF
	RED   int = 0xFF0000
	GREEN int = 0x00FF00
	BLUE  int = 0x0000FF
	BLACK int = 0x000000
)

// Go does not support inheritance of structs/classes. However, it
// supports embedding types. The following struct embeds Circle. Because of
// *member set promotion*, all members of Circle become available on
// ColoredCircle, too.

type ColoredCircle struct {
	Circle
	Color int
}

func (c ColoredCircle) GetColor() int {
	return c.Color
}

type Colored interface {
	GetColor() int
}

func main() {
	// Note the declare-and-assign syntax in the next line. The compiler
	// automatically determines the type of r, no need to specify it explicitly.
	r := Rect{LeftUpper: Point{X: 0, Y: 0}, RightLower: Point{X: 10, Y: 10}}
	c := Circle{Center: Point{X: 5, Y: 5}, Radius: 5}

	// Note that we can access the Radius of the ColoredCircle although Radius
	// is a member of the embedded type Circle.
	cc := ColoredCircle{c, RED}
	fmt.Printf("Colored circle has radius %f\n", cc.Radius)

	// Next, we create an array of shapes. As you can see, rect
	// and circle are compatible with the shape interface
	shapes := []Shape{r, c, cc}

	// Note the use of range in the for loop. In contrast to C#, the for-range
	// loop provides the value and an index.
	for ix, shape := range shapes {
		fmt.Printf("Area of shape %d (%T) is %f\n", ix, shape, shape.Area())

		// Note the syntax of the if statement in Go. You can write
		// declare-and-assign and boolean expression in a single line.
		// Very convenient once you got used to it.

		// Additionally note how we check whether shape is compatible
		// with the Colored interface. You get back a variable with the
		// correct type and a bool indicator indicating if the cast was ok.
		if colCirc, ok := shape.(Colored); ok {
			fmt.Printf("\thas color %x\n", colCirc.GetColor())
		}
	}

	r.Enlarge(2)
	fmt.Printf("Rectangle's area after enlarging it is %f\n", r.Area())
}
