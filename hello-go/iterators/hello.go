package main

// Read more at https://go.dev/blog/range-functions

import (
	"fmt"
	"iter"
)

func main() {
	simple_iterator()
	fmt.Println("---")
	breaking_from_loop()
	fmt.Println("---")
	pull_iterator()
	fmt.Println("---")
	pull_iterator_from_push()
	fmt.Println("---")
	adapters()
}

func simple_iterator() {
	// iter is of type Seq[int]
	iter := func(yield func(int) bool) {
		for i := 0; i < 10; i++ {
			if !yield(i) {
				return
			}
		}
	}

	for v := range iter {
		fmt.Println(v)
	}
}

func breaking_from_loop() {
	// iter is of type Seq[int]
	iter := func(yield func(int) bool) {
		for i := 0; true; i++ {
			fmt.Println("yielding", i)
			if !yield(i) {
				fmt.Println("breaking")
				return
			}
		}
	}

	for v := range iter {
		fmt.Println(v)
		if v == 15 {
			break
		}
	}
}

func pull_iterator() {
	pull_iter := func() (func() (int, bool), func()) {
		current := 0

		next := func() (int, bool) {
			if current >= 10 {
				return 0, false
			}
			current++
			return current, true
		}

		stop := func() {
			current = 0
		}

		return next, stop
	}

	next, stop := pull_iter()
	defer stop()

	for v, ok := next(); ok; v, ok = next() {
		fmt.Println(v)
	}
}

func pull_iterator_from_push() {
	// simple push iterator
	push_iter := func(yield func(int) bool) {
		for i := 1; i <= 10; i++ {
			if !yield(i) {
				return
			}
		}
	}

	next, stop := iter.Pull(push_iter)
	defer stop()

	for v, ok := next(); ok; v, ok = next() {
		fmt.Println(v)
	}
}

func Filter[V any](f func(V) bool, s iter.Seq[V]) iter.Seq[V] {
	return func(yield func(V) bool) {
		for v := range s {
			if f(v) {
				if !yield(v) {
					return
				}
			}
		}
	}
}

func Map[V any, W any](f func(V) W, s iter.Seq[V]) iter.Seq[W] {
	return func(yield func(W) bool) {
		for v := range s {
			if !yield(f(v)) {
				return
			}
		}
	}
}

func adapters() {
	iter := func(yield func(int) bool) {
		for i := 1; i <= 10; i++ {
			if !yield(i) {
				return
			}
		}
	}

	for v := range Map(func(i int) int { return i * 2 }, Filter(func(i int) bool { return i%2 == 0 }, iter)) {
		fmt.Println(v)
	}
}
