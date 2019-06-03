#include <iostream>
#include <bits/stdc++.h> 
#include <chrono>

using namespace std; 
using namespace std::chrono;

#define SIDELENGTH 8
#define SIZE SIDELENGTH*SIDELENGTH

size_t getIndex(size_t x, size_t y) {
  return y*SIDELENGTH + x;
}

bool hasQueen(bitset<SIZE> const & board, size_t x, size_t y) {
  return board[getIndex(x, y)];
}

bool tryPlaceQueen(bitset<SIZE> & board, size_t x, size_t y) {
	if (x >= SIDELENGTH || y >= SIDELENGTH) {
		return false;
	}

	if (hasQueen(board, x, y)) {
		return false;
	}

	for (size_t i = 1; i < SIDELENGTH; i++) {
    auto right = x+i, left = x-i, top = y-i, down = y+i;
		auto rightInside = right < SIDELENGTH, leftInside = left < SIDELENGTH, topInside = top < SIDELENGTH, downInside = down < SIDELENGTH;
		if ((rightInside && (hasQueen(board, right, y) || (topInside && hasQueen(board, right, top)) || (downInside && hasQueen(board, right, down)))) ||
			(leftInside && (hasQueen(board, left, y) || (topInside && hasQueen(board, left, top)) || (downInside && hasQueen(board, left, down)))) ||
			(topInside && hasQueen(board, x, top)) || (downInside && hasQueen(board, x, down))) {
			return false;
		}
	}

  board.set(getIndex(x, y), 1);
	return true;
}

bool removeQueen(bitset<SIZE> & board, size_t x, size_t y) {
	if (x >= SIDELENGTH || y >= SIDELENGTH) {
		return false;
	}

  board.set(getIndex(x, y), 0);
	return true;
}

// Print prints the given chess board to the given writer
void Print(bitset<SIZE> const & board) {
	// Credits: Original code for this method see https://github.com/danrl/golibby/blob/9dd8757e94746578c5a9c0e4ca9d5a347fd7de32/queensboard/queensboard.go#L207
	//          Under MIT license (https://github.com/danrl/golibby/blob/master/LICENSE)

	// board framing
	string top ="┏";
	string middle = "┠";
	string bottom = "┗";
	for (size_t j = 1; j < SIDELENGTH; j++) {
		top = top.append("━━┯");
		middle = middle.append("──┼");
		bottom = bottom.append("━━┷");
	}
	top = top.append("━━┓\n");
	middle = middle.append("──┨\n");
	bottom = bottom.append("━━┛\n");

	// compile the field
  string out = top;
	for (size_t i = 0; i < SIDELENGTH; i++) {
		out = out.append("┃");
		for (size_t j = 0; j < SIDELENGTH; j++) {
			if (board[getIndex(j, i)]) {
				out = out.append(". ");
			} else {
        out = out.append("  ");
			}
			if (j < (SIDELENGTH - 1)) {
				out = out.append("│");
			}
		}
		out = out.append("┃\n");
		if (i < (SIDELENGTH - 1)) {
      out = out.append(middle);
		}
	}
  out = out.append(bottom);
  cout << out << endl;
}

void findSolutions(bitset<SIZE> & board, size_t x, vector< bitset<SIZE> > & solutions) {
	for (size_t i = 0; i < SIDELENGTH; i++) {
		if (tryPlaceQueen(board, x, i)) {
			if (x == (SIDELENGTH - 1)) {
        solutions.push_back(bitset<SIZE>(board));
				removeQueen(board, x, i);
			}

			findSolutions(board, x+1, solutions);
			removeQueen(board, x, i);
		}
	}
}

// FindSolutions finds all solutions of the queens problem on a chess board with the given side length
vector< bitset<SIZE> > FindSolutions() {
  vector< bitset<SIZE> > solutions;
  bitset<SIZE> board;
	findSolutions(board, 0, solutions);
  return solutions;
}

int main() {
  auto t1 = high_resolution_clock::now();
  auto solutions = FindSolutions();
  auto t2 = high_resolution_clock::now();

  auto duration = duration_cast<microseconds>(t2 - t1).count();

  for (size_t i = 0; i < solutions.size(); i++) {
    Print(solutions[i]);
  }

  cout << "Found " << solutions.size() << " solutions" << endl;
  cout << "Took " << (duration % 1000000000) / 1000000 << "." << (duration % 1000000) << " seconds" << endl;
  return 0;
}
