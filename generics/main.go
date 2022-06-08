package main

import (
	"fmt"
	"reflect"

	"github.com/samber/lo"
	"github.com/thoas/go-funk"
	"golang.org/x/exp/constraints"
)

type eatOrKeep interface {
	// Returns true if item should be eaten
	shouldEat() bool
}

//#region Lentil
type lentil struct {
	isGood bool
}

func (l lentil) shouldEat() bool { return !l.isGood }

//#endregion

//#region Snail
type snail struct {
	hasHouse bool
}

func (s snail) shouldEat() bool { return !s.hasHouse }

//#endregion

//#region Bird
type bird struct{}

// Typed version of processing function.
//
// Removes all items from the slice that should be eaten.
func (p bird) process(items []eatOrKeep) []eatOrKeep {
	result := []eatOrKeep{}
	for _, item := range items {
		if !item.shouldEat() {
			result = append(result, item)
		}
	}

	return result
}

//#endregion

//#region Reflection

// Reflection-based implementation of the filtering.
//
// First parameter must be the slice to be filtered. Second parameter
// must be the filter function returning boolean. Function removes all items
// from slice for which filter function returns false.
func processInterface(itemsSlice interface{}, filter interface{}) interface{} {
	// Get reflection value for slice to process. Note that error handling
	// (e.g. if itemsSlice is not a slice) is not implemented to keep example simple
	// as this article is not a tutorial on Go Runtime Reflection.
	items := reflect.ValueOf(itemsSlice)

	// Create slice to store results
	result := reflect.MakeSlice(items.Type(), 0, 0)

	// Get reflection value for filter function.
	funcValue := reflect.ValueOf(filter)

	// Iterate over all items
	for i := 0; i < items.Len(); i++ {
		// Get item on index i using reflection
		item := items.Index(i)

		// Call filter function using reflection and check result.
		keep := funcValue.Call([]reflect.Value{item})
		if keep[0].Interface().(bool) {
			// Append item to result slice using reflection
			result = reflect.Append(result, item)
		}
	}

	// Return result slice as interface{}
	return result.Interface()
}

//#endregion

//#region Generics
func process[I any](items []I, filter func(i I) bool) []I {
	result := []I{}
	for _, item := range items {
		if filter(item) {
			result = append(result, item)
		}
	}

	return result
}

//#endregion

//#region Typed struct
type itemsGroup struct {
	eatOrKeep
	count int
}

type itemsBag struct {
	bag []itemsGroup
}

func newItemsBag() *itemsBag {
	return &itemsBag{
		bag: make([]itemsGroup, 0),
	}
}

func (b *itemsBag) append(item eatOrKeep) {
	if len(b.bag) == 0 || item.shouldEat() != b.bag[len(b.bag)-1].shouldEat() {
		b.bag = append(b.bag, itemsGroup{eatOrKeep: item, count: 1})
	} else {
		b.bag[len(b.bag)-1].count++
	}
}

func (b itemsBag) getItems() []eatOrKeep {
	result := make([]eatOrKeep, 0)
	for _, group := range b.bag {
		for i := 0; i < group.count; i++ {
			result = append(result, group.eatOrKeep)
		}
	}

	return result
}

//#endregion

//#region Generic struct
type genericItemsGroup[T any] struct {
	item  T
	count int
}

type genericItemsBag[T any] struct {
	bag              []genericItemsGroup[T]
	equalityComparer func(T, T) bool
}

func newGenericItemsBag[T any](comparer func(T, T) bool) *genericItemsBag[T] {
	return &genericItemsBag[T]{
		bag:              make([]genericItemsGroup[T], 0),
		equalityComparer: comparer,
	}
}

func (b *genericItemsBag[T]) append(item T) {
	if len(b.bag) == 0 || !b.equalityComparer(item, b.bag[len(b.bag)-1].item) {
		b.bag = append(b.bag, genericItemsGroup[T]{item: item, count: 1})
	} else {
		b.bag[len(b.bag)-1].count++
	}
}

func (b genericItemsBag[T]) getItems() []T {
	result := make([]T, 0)
	for _, group := range b.bag {
		for i := 0; i < group.count; i++ {
			result = append(result, group.item)
		}
	}

	return result
}

//#endregion

//#region Generic channels
func processChannel[I any](items <-chan I, filter func(i I) bool) <-chan I {
	out := make(chan I)
	go func() {
		defer close(out)
		for item := range items {
			if filter(item) {
				out <- item
			}
		}
	}()
	return out
}

//#endregion

//#region Type constraints
const (
	SMALL  = 1
	MEDIUM = 2
	LARGE  = 3
)

type sizedLentil struct {
	lentil
	lentilSize int
}

func (l sizedLentil) size() int { return l.lentilSize }

type sized interface {
	size() int
}

type sizedEatOrKeep interface {
	sized
	eatOrKeep
}

func processAndSort[I sizedEatOrKeep](items []I, filter func(i I) bool) []I {
	result := []I{}
	for _, item := range items {
		if filter(item) {
			result = append(result, item)
		}
	}

	bubblesort(result, func(item I) int { return item.size() })
	return result
}

func bubblesort[I any, O constraints.Ordered](items []I, toOrdered func(item I) O) {
	for itemCount := len(items) - 1; ; itemCount-- {
		hasChanged := false
		for index := 0; index < itemCount; index++ {
			if toOrdered(items[index]) > toOrdered(items[index+1]) {
				items[index], items[index+1] = items[index+1], items[index]
				hasChanged = true
			}
		}
		if !hasChanged {
			break
		}
	}
}
//#endregion

func main() {
	items := []eatOrKeep{
		lentil{isGood: true},
		lentil{isGood: false},
		lentil{isGood: false},
		snail{hasHouse: true},
		snail{hasHouse: false},
	}
	processedItems := bird{}.process(items)
	fmt.Println("Eaten:", len(items)-len(processedItems), "Kept:", len(processedItems))

	interfacedItems := processInterface(items, func(item eatOrKeep) bool { return !item.shouldEat() })
	processedItems = interfacedItems.([]eatOrKeep)
	fmt.Println("Eaten:", len(items)-len(processedItems), "Kept:", len(processedItems))

	processedItems = process(items, func(item eatOrKeep) bool { return !item.shouldEat() })
	//processedItems = process[eatOrKeep](items, func(item eatOrKeep) bool { return !item.shouldEat() })
	fmt.Println("Eaten:", len(items)-len(processedItems), "Kept:", len(processedItems))

	processedItems = lo.Filter(items, func(item eatOrKeep, _ int) bool { return !item.shouldEat() })
	fmt.Println("Eaten:", len(items)-len(processedItems), "Kept:", len(processedItems))

	interfacedItems = funk.Filter(items, func(item eatOrKeep) bool { return !item.shouldEat() })
	processedItems = interfacedItems.([]eatOrKeep)
	fmt.Println("Eaten:", len(items)-len(processedItems), "Kept:", len(processedItems))

	bag := newItemsBag()
	bag.append(lentil{isGood: true})
	bag.append(lentil{isGood: true})
	bag.append(lentil{isGood: false})
	bag.append(lentil{isGood: false})
	processedItems = process(bag.getItems(), func(item eatOrKeep) bool { return !item.shouldEat() })
	fmt.Println("Eaten:", len(items)-len(processedItems), "Kept:", len(processedItems))

	genericBag := newGenericItemsBag(func(lhs eatOrKeep, rhs eatOrKeep) bool { return lhs.shouldEat() == rhs.shouldEat() })
	genericBag.append(lentil{isGood: true})
	genericBag.append(lentil{isGood: true})
	genericBag.append(lentil{isGood: false})
	genericBag.append(lentil{isGood: false})
	processedItems = process(genericBag.getItems(), func(item eatOrKeep) bool { return !item.shouldEat() })
	fmt.Println("Eaten:", len(items)-len(processedItems), "Kept:", len(processedItems))

	// Fill buffered channel with lentils
	in := make(chan eatOrKeep, 5)
	in <- lentil{isGood: true}
	in <- lentil{isGood: true}
	in <- lentil{isGood: false}
	in <- lentil{isGood: false}
	in <- lentil{isGood: false}
	close(in)

	// Use generic channel processing
	total := len(in)
	remaining := 0
	for range processChannel(in, func(item eatOrKeep) bool { return !item.shouldEat() }) {
		remaining++
	}
	fmt.Println("Eaten:", total-remaining, "Kept:", remaining)

	sizedItems := []sizedEatOrKeep{
		sizedLentil{lentilSize: LARGE, lentil: lentil{isGood: true}},
		sizedLentil{lentilSize: MEDIUM, lentil: lentil{isGood: false}},
		sizedLentil{lentilSize: SMALL, lentil: lentil{isGood: true}},
	}
    processedOrderd := processAndSort(sizedItems, func(item sizedEatOrKeep) bool { return !item.shouldEat() })
	fmt.Println("Eaten:", len(sizedItems)-len(processedOrderd), "Kept:", len(processedOrderd))
    for _, sortedItem := range processedOrderd {
        fmt.Println("Size:", sortedItem.size())
    }
}
