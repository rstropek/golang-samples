package queensproblembitarraysolver

import (
	"bytes"
	"io"

	ba "github.com/golang-collections/go-datastructures/bitarray"
)

func getIndex(x byte, y byte) uint64 {
	return (uint64)(y*8 + x)
}

func hasQueen(board ba.BitArray, x byte, y byte) bool {
	result, _ := board.GetBit(getIndex(x, y))
	return result
}

func tryPlaceQueen(board ba.BitArray, x byte, y byte) bool {
	if x > 7 || y > 7 {
		return false
	}

	if hasQueen(board, x, y) {
		return false
	}

	var i byte
	for i = 1; i < 8; i++ {
		right, left, top, down := x+i, x-i, y-i, y+i
		rightInside, leftInside, topInside, downInside := right < 8, left < 8, top < 8, down < 8
		if (rightInside && (hasQueen(board, right, y) || (topInside && hasQueen(board, right, top)) || (downInside && hasQueen(board, right, down)))) ||
			(leftInside && (hasQueen(board, left, y) || (topInside && hasQueen(board, left, top)) || (downInside && hasQueen(board, left, down)))) ||
			(topInside && hasQueen(board, x, top)) || (downInside && hasQueen(board, x, down)) {
			return false
		}
	}

	board.SetBit(getIndex(x, y))

	return true
}

func removeQueen(board ba.BitArray, x byte, y byte) bool {
	if x > 7 || y > 7 {
		return false
	}

	board.ClearBit(getIndex(x, y))
	return true
}

func Print(board ba.BitArray, w io.Writer) {
	// board framing
	var top, middle, bottom bytes.Buffer
	top.WriteRune('┏')
	middle.WriteRune('┠')
	bottom.WriteRune('┗')
	for j := 1; j < 8; j++ {
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
	for i = 0; i < 8; i++ {
		out.WriteRune('┃')
		var j byte
		for j = 0; j < 8; j++ {
			if hasQueen, _ := board.GetBit(getIndex(j, i)); hasQueen {
				out.WriteRune('.')
				out.WriteRune(' ')
			} else {
				out.WriteRune(' ')
				out.WriteRune(' ')
			}
			if j < 7 {
				out.WriteRune('│')
			}
		}
		out.Write([]byte("┃\n"))
		if i < 7 {
			out.Write(middle.Bytes())
		}
	}
	out.Write(bottom.Bytes())
	w.Write(out.Bytes())
}

func FindSolution(board ba.BitArray, x byte, solutions []ba.BitArray) []ba.BitArray {
	var i byte
	for i = 0; i < 8; i++ {
		if tryPlaceQueen(board, x, i) {
			if x == 7 {
				solutions = append(solutions, ba.NewBitArray(8*8).Or(board))
				removeQueen(board, x, i)
			}

			solutions = FindSolution(board, x+1, solutions)
			removeQueen(board, x, i)
		}
	}

	return solutions
}
