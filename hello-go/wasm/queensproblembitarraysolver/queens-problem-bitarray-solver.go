package queensproblembitarraysolver

import (
	"bytes"
	"io"
	"time"

	ba "github.com/golang-collections/go-datastructures/bitarray"
)

func getIndex(x byte, y byte, sideLength byte) uint64 {
	// Note: No semicolons ;-)
	return (uint64)(y*sideLength + x)
}

func hasQueen(board ba.BitArray, sideLength byte, x byte, y byte) bool {
	// Note error handling here. Go has no exceptions.
	// By convention, errors are the last return value and of type `error`
	// (see https://gobyexample.com/errors)
	result, _ := board.GetBit(getIndex(x, y, sideLength))
	return result
}

func tryPlaceQueen(board ba.BitArray, sideLength byte, x byte, y byte) bool {
	// Note: No parentheses ;-)
	if x >= sideLength || y >= sideLength {
		return false
	}

	if hasQueen(board, sideLength, x, y) {
		return false
	}

	var i byte
	for i = 1; i < sideLength; i++ {
		// The following calculations consider byte overflow (0-1 becoming 255)

		// Note declare-and-initialize syntax
		// (see https://gobyexample.com/variables)
		right, left, top, down := x+i, x-i, y-i, y+i
		rightInside, leftInside, topInside, downInside := right < sideLength, left < sideLength, top < sideLength, down < sideLength
		if (rightInside && (hasQueen(board, sideLength, right, y) || (topInside && hasQueen(board, sideLength, right, top)) || (downInside && hasQueen(board, sideLength, right, down)))) ||
			(leftInside && (hasQueen(board, sideLength, left, y) || (topInside && hasQueen(board, sideLength, left, top)) || (downInside && hasQueen(board, sideLength, left, down)))) ||
			(topInside && hasQueen(board, sideLength, x, top)) || (downInside && hasQueen(board, sideLength, x, down)) {
			return false
		}
	}

	board.SetBit(getIndex(x, y, sideLength))
	return true
}

func removeQueen(board ba.BitArray, sideLength byte, x byte, y byte) bool {
	if x >= sideLength || y >= sideLength {
		return false
	}

	board.ClearBit(getIndex(x, y, sideLength))
	return true
}

// Print prints the given chess board to the given writer
func Print(board ba.BitArray, sideLength byte, w io.Writer) {
	// Credits: Original code for this method see https://github.com/danrl/golibby/blob/9dd8757e94746578c5a9c0e4ca9d5a347fd7de32/queensboard/queensboard.go#L207
	//          Under MIT license (https://github.com/danrl/golibby/blob/master/LICENSE)

	// Note naming schema: Uppercase functions are exported, lowercase functions are local

	// board framing
	var top, middle, bottom bytes.Buffer
	top.WriteRune('┏')
	middle.WriteRune('┠')
	bottom.WriteRune('┗')
	var j byte
	for j = 1; j < sideLength; j++ {
		top.Write([]byte("━━┯"))
		middle.Write([]byte("──┼"))
		bottom.Write([]byte("━━┷"))
	}
	top.Write([]byte("━━┓\n"))
	middle.Write([]byte("──┨\n"))
	bottom.Write([]byte("━━┛\n"))

	// compile the field
	out := bytes.NewBuffer(top.Bytes())
	var i byte
	for i = 0; i < sideLength; i++ {
		out.WriteRune('┃')
		var j byte
		for j = 0; j < sideLength; j++ {
			if hasQueen, _ := board.GetBit(getIndex(j, i, sideLength)); hasQueen {
				out.WriteRune('.')
				out.WriteRune(' ')
			} else {
				out.WriteRune(' ')
				out.WriteRune(' ')
			}
			if j < (sideLength - 1) {
				out.WriteRune('│')
			}
		}
		out.Write([]byte("┃\n"))
		if i < (sideLength - 1) {
			out.Write(middle.Bytes())
		}
	}
	out.Write(bottom.Bytes())
	w.Write(out.Bytes())
}

func findSolutions(board ba.BitArray, sideLength byte, x byte, solutions []ba.BitArray) []ba.BitArray {
	var i byte
	for i = 0; i < sideLength; i++ {
		if tryPlaceQueen(board, sideLength, x, i) {
			if x == (sideLength - 1) {
				solutions = append(solutions, ba.NewBitArray((uint64)(sideLength*sideLength)).Or(board))
				removeQueen(board, sideLength, x, i)
				return solutions
			}

			solutions = findSolutions(board, sideLength, x+1, solutions)
			removeQueen(board, sideLength, x, i)
		}
	}

	return solutions
}

// Result represents the result of the search for solutions to an n queens problem
type Result struct {
	Solutions       []ba.BitArray
	CalculationTime time.Duration
}

// FindSolutions finds all solutions of the queens problem on a chess board with the given side length
func FindSolutions(sideLength byte) Result {
	start := time.Now()
	solutions := findSolutions(ba.NewBitArray((uint64)(sideLength*sideLength)), sideLength, 0, make([]ba.BitArray, 0))
	return Result{
		Solutions:       solutions,
		CalculationTime: time.Since(start),
	}
}
