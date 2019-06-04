function getIndex(x, y, sideLength) {
	return y*sideLength + x;
}

function hasQueen(board, sideLength, x, y) {
    const result = board[getIndex(x, y, sideLength)];
	return result;
}

function tryPlaceQueen(board, sideLength, x, y) {
	if (x >= sideLength || y >= sideLength) {
		return false;
	}

	if (hasQueen(board, sideLength, x, y)) {
		return false;
	}

	for (let i = 1; i < sideLength; i++) {
		const [right, left, top, down] = [x+i, x-i, y-i, y+i];
		const [rightInside, leftInside, topInside, downInside] = [right < sideLength, left >= 0, top >= 0, down < sideLength];
		if ((rightInside && (hasQueen(board, sideLength, right, y) || (topInside && hasQueen(board, sideLength, right, top)) || (downInside && hasQueen(board, sideLength, right, down)))) ||
			(leftInside && (hasQueen(board, sideLength, left, y) || (topInside && hasQueen(board, sideLength, left, top)) || (downInside && hasQueen(board, sideLength, left, down)))) ||
			(topInside && hasQueen(board, sideLength, x, top)) || (downInside && hasQueen(board, sideLength, x, down))) {
			return false;
		}
	}

    board[getIndex(x, y, sideLength)] = 1;
	return true;
}

function removeQueen(board, sideLength, x, y) {
	if (x >= sideLength || y >= sideLength) {
		return false;
	}

    board[getIndex(x, y, sideLength)] = 0;
	return true;
}

function printBoard(board, sideLength) {
	// board framing
	var top, middle, bottom;
	top = '┏';
	middle = '┠';
	bottom = '┗';
	for (let j = 1; j < sideLength; j++) {
		top += "━━┯";
		middle += "──┼";
		bottom += "━━┷";
	}
	top += "━━┓\n";
	middle += "──┨\n";
	bottom += "━━┛\n";

    let output = top;
	for (let i = 0; i < sideLength; i++) {
		output += '┃';
		for (let j = 0; j < sideLength; j++) {
			if (board[getIndex(j, i, sideLength)]) {
				output += '. ';
			} else {
				output += '  ';
			}
			if (j < (sideLength - 1)) {
				output += '│';
			}
		}
		output += "┃\n";
		if (i < (sideLength - 1)) {
			output += middle;
		}
    }
    output += bottom;
    console.log(output);
}

function findSolutions(board, sideLength, x, solutions) {
	for(let i = 0; i < sideLength; i++) {
		if (tryPlaceQueen(board, sideLength, x, i)) {
			if (x == (sideLength - 1)) {
				solutions.push(new Uint8Array(board));
				removeQueen(board, sideLength, x, i);
				return solutions;
			}

			solutions = findSolutions(board, sideLength, x+1, solutions);
			removeQueen(board, sideLength, x, i);
		}
	}

	return solutions;
}

function FindSolutions(sideLength) {
	const start = performance.now();
	const solutions = findSolutions(new Uint8Array(sideLength * sideLength), sideLength, 0, []);
	return {
		solutions: solutions,
		calculationTime: performance.now() - start,
	}
}
