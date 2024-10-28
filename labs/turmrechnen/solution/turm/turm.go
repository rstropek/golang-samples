package turm

import (
	"fmt"
	"iter"
)

type Turm struct {
	startValue int
	height     int
}

func NewTurm(startValue int, height int) (*Turm, error) {
	if startValue <= 1 || height <= 2 {
		return nil, fmt.Errorf("invalid parameters: startValue > 1 and height > 2 required")
	}

	return &Turm{startValue, height}, nil
}

type TurmIntermediateResult struct {
	OldValue  int
	Operation rune
	Operand   int
	NewValue  int
}

func (t Turm) Calculate() []TurmIntermediateResult {
	results := make([]TurmIntermediateResult, 0, (t.height-1)*2)

	value := t.startValue
	for phase := 0; phase < 2; phase++ {
		for i := 2; i <= t.height; i++ {
			var r TurmIntermediateResult
			switch phase {
			case 0:
				r = TurmIntermediateResult{value, '*', i, value * i}
			case 1:
				r = TurmIntermediateResult{value, '/', i, value / i}
			}
			results = append(results, r)
			value = r.NewValue
		}
	}

	return results
}

func (t Turm) CalculateIterative() iter.Seq[TurmIntermediateResult] {
	return func(yield func(TurmIntermediateResult) bool) {
		value := t.startValue
		for phase := 0; phase < 2; phase++ {
			for i := 2; i <= t.height; i++ {
				var r TurmIntermediateResult
				switch phase {
				case 0:
					r = TurmIntermediateResult{value, '*', i, value * i}
				case 1:
					r = TurmIntermediateResult{value, '/', i, value / i}
				}

				if !yield(r) {
					return
				}

				value = r.NewValue
			}
		}
	}
}
