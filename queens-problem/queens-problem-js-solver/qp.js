function getIndex(x, y) {
	return y*8 + x;
}

function hasQueen(board, x, y) {
    const result = board[getIndex(x, y)];
	return result;
}

function tryPlaceQueen(board, x, y) {
	if (x > 7 || y > 7) {
		return false;
	}

	if (hasQueen(board, x, y)) {
		return false;
	}

	for (let i = 1; i < 8; i++) {
		const [right, left, top, down] = [x+i, x-i, y-i, y+i];
		const [rightInside, leftInside, topInside, downInside] = [right < 8, left >= 0, top >= 0, down < 8];
		if ((rightInside && (hasQueen(board, right, y) || (topInside && hasQueen(board, right, top)) || (downInside && hasQueen(board, right, down)))) ||
			(leftInside && (hasQueen(board, left, y) || (topInside && hasQueen(board, left, top)) || (downInside && hasQueen(board, left, down)))) ||
			(topInside && hasQueen(board, x, top)) || (downInside && hasQueen(board, x, down))) {
			return false;
		}
	}

    board[getIndex(x, y)] = 1;
	return true;
}

function removeQueen(board, x, y) {
	if (x > 7 || y > 7) {
		return false;
	}

    board[getIndex(x, y)] = 0;
	return true;
}

function printBoard(board) {
	// board framing
	var top, middle, bottom;
	top = '┏';
	middle = '┠';
	bottom = '┗';
	for (let j = 1; j < 8; j++) {
		top += "━━┯";
		middle += "──┼";
		bottom += "━━┷";
	}
	top += "━━┓\n";
	middle += "──┨\n";
	bottom += "━━┛\n";

    let output = top;
	for (let i = 0; i < 8; i++) {
		output += '┃';
		for (let j = 0; j < 8; j++) {
			if (board[getIndex(j, i)]) {
				output += '. ';
			} else {
				output += '  ';
			}
			if (j < 7) {
				output += '│';
			}
		}
		output += "┃\n";
		if (i < 7) {
			output += middle;
		}
    }
    output += bottom;
    console.log(output);
}

function findSolution(board, x, solutions) {
	for(let i = 0; i < 8; i++) {
		if (tryPlaceQueen(board, x, i)) {
			if (x == 7) {
				solutions.push(new Uint8Array(board));
				removeQueen(board, x, i);
			}

			solutions = findSolution(board, x+1, solutions);
			removeQueen(board, x, i);
		}
	}

	return solutions;
}
