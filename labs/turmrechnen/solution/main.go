package main

import (
	"fmt"
	"os"
	"strconv"

	"github.com/rstropek/golang-samples/turmrechnen/turm"
)

func main() {
	t, err := TurmFromArguments()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	results := t.Calculate()

	maxLengthOldValue, maxLengthOperand := getMaxLength(results)

	for _, r := range results {
		fmt.Printf("%*d %c %*d = %d\n", maxLengthOldValue, r.OldValue, r.Operation, -maxLengthOperand, r.Operand, r.NewValue)
	}

	for r := range t.CalculateIterative() {
		fmt.Printf("%*d %c %*d = %d\n", maxLengthOldValue, r.OldValue, r.Operation, -maxLengthOperand, r.Operand, r.NewValue)
		break
	}
}

func getMaxLength(results []turm.TurmIntermediateResult) (int, int) {
	maxLengthOldValue := 0
	maxLengthOperand := 0
	for _, r := range results {
		length := len(strconv.Itoa(r.OldValue))
		if length > maxLengthOldValue {
			maxLengthOldValue = length
		}
		length = len(strconv.Itoa(r.Operand))
		if length > maxLengthOperand {
			maxLengthOperand = length
		}
	}
	return maxLengthOldValue, maxLengthOperand
}

func TurmFromArguments() (*turm.Turm, error) {
	if len(os.Args) != 3 {
		return nil, fmt.Errorf("usage: turmrechnen <start value> <height>")
	}

	startValue, err := strconv.Atoi(os.Args[1])
	if err != nil {
		return nil, fmt.Errorf("invalid start value, it must be an integer > 1")
	}

	height, err := strconv.Atoi(os.Args[2])
	if err != nil {
		return nil, fmt.Errorf("invalid height, it must be an integer > 2")
	}

	return turm.NewTurm(startValue, height)
}
