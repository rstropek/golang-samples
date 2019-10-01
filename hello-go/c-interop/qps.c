#include <stdbool.h>
#include <string.h>
#include <stdio.h>

#define SIDELENGTH 8
#define SIZE (SIDELENGTH*SIDELENGTH)

int getIndex(int x, int y) {
  return y*SIDELENGTH + x;
}

bool hasQueen(const bool* board, int x, int y) {
  return board[getIndex(x, y)];
}

bool tryPlaceQueen(bool* board, int x, int y) {
  if (hasQueen(board, x, y)) {
    return false;
  }

  for (int i = 1; i < SIDELENGTH; i++) {
    int right = x+i, left = x-i, top = y-i, down = y+i;
    bool rightInside = right < SIDELENGTH, leftInside = left < SIDELENGTH && left >= 0, topInside = top < SIDELENGTH && top >= 0, downInside = down < SIDELENGTH;
    if ((rightInside && (hasQueen(board, right, y) || (topInside && hasQueen(board, right, top)) || (downInside && hasQueen(board, right, down)))) ||
      (leftInside && (hasQueen(board, left, y) || (topInside && hasQueen(board, left, top)) || (downInside && hasQueen(board, left, down)))) ||
      (topInside && hasQueen(board, x, top)) || (downInside && hasQueen(board, x, down))) {
      return false;
    }
  }

  board[getIndex(x, y)] = true;
  return true;
}

bool removeQueen(bool* board, int x, int y) {
  if (x >= SIDELENGTH || y >= SIDELENGTH || x < 0 || y < 0) {
    return false;
  }

  board[getIndex(x, y)] = false;
  return true;
}

int findSolutions(bool* board, int x) {
  int numberOfSolutions = 0;
  for (int i = 0; i < SIDELENGTH; i++) {
    if (tryPlaceQueen(board, x, i)) {
      if (x == (SIDELENGTH - 1)) {
        removeQueen(board, x, i);
        return 1;
      }

      numberOfSolutions += findSolutions(board, x+1);
      removeQueen(board, x, i);
    }
  }

  return numberOfSolutions;
}

int calculateNumberOfSolutions(int boardSize) {
  bool board[boardSize * boardSize];
  memset(board, 0, boardSize * boardSize);

  return findSolutions(board, 0);
}
