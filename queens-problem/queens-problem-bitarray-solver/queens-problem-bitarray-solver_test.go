package queensproblembitarraysolver

import (
	"testing"

	ba "github.com/golang-collections/go-datastructures/bitarray"
)

func TestGetIndex(t *testing.T) {
	if result := getIndex(0, 0, 8); result != 0 {
		t.Errorf("Index for 0/0 wrong, expected 0, got %d", result)
	}

	if result := getIndex(1, 0, 8); result != 1 {
		t.Errorf("Index for 1/0 wrong, expected 1, got %d", result)
	}

	if result := getIndex(0, 1, 8); result != 8 {
		t.Errorf("Index for 0/1 wrong, expected 8, got %d", result)
	}

	if result := getIndex(1, 1, 8); result != 9 {
		t.Errorf("Index for 1/1 wrong, expected 9, got %d", result)
	}
}

func TestTryPlaceQueenInvalidParameters(t *testing.T) {
	board := ba.NewBitArray(8 * 8)

	if tryPlaceQueen(board, 8, 8, 0) {
		t.Error("Did not recognize invalid index for x")
	}

	if tryPlaceQueen(board, 8, 0, 8) {
		t.Error("Did not recognize invalid index for y")
	}
}

func TestTryPlaceQueen(t *testing.T) {
	board := ba.NewBitArray(8 * 8)
	if !tryPlaceQueen(board, 8, 0, 0) {
		t.Error("Could not place queen on empty board")
	}
}

func TestRemoveQueen(t *testing.T) {
	board := ba.NewBitArray(8 * 8)
	tryPlaceQueen(board, 8, 0, 0)
	removeQueen(board, 8, 0, 0)
	if !tryPlaceQueen(board, 8, 0, 0) {
		t.Error("Could not place queen after removing it")
	}
}

func TestTryPlaceQueenSameLocation(t *testing.T) {
	board := ba.NewBitArray(8 * 8)
	tryPlaceQueen(board, 8, 0, 0)
	if tryPlaceQueen(board, 8, 0, 0) {
		t.Error("Successfully placed queen on same location tiwce")
	}
}

func TestTryPlaceDetectInvalidPosition(t *testing.T) {
	board := ba.NewBitArray(8 * 8)
	tryPlaceQueen(board, 8, 3, 3)
	if tryPlaceQueen(board, 8, 3, 0) || tryPlaceQueen(board, 8, 3, 7) ||
		tryPlaceQueen(board, 8, 0, 3) || tryPlaceQueen(board, 8, 7, 3) ||
		tryPlaceQueen(board, 8, 0, 6) || tryPlaceQueen(board, 8, 6, 0) ||
		tryPlaceQueen(board, 8, 0, 0) || tryPlaceQueen(board, 8, 7, 7) {
		t.Error("Invalid position not detected")
	}
}

func TestFindSolutions(t *testing.T) {
	if solutions := FindSolutions(8); len(solutions.Solutions) != 92 {
		t.Errorf("Expected 92 solutions, got %d", len(solutions.Solutions))
	}
}
