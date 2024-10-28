package turm

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTurm_Calculate(t *testing.T) {
	turm, err := NewTurm(5, 3)
	assert.NoError(t, err)
	expectedResults := []TurmIntermediateResult{
		{OldValue: 5, Operation: '*', Operand: 2, NewValue: 10},
		{OldValue: 10, Operation: '*', Operand: 3, NewValue: 30},
		{OldValue: 30, Operation: '/', Operand: 2, NewValue: 15},
		{OldValue: 15, Operation: '/', Operand: 3, NewValue: 5},
	}

	results := turm.Calculate()

	assert.Equal(t, expectedResults, results)
}

func TestTurm_InvalidConstructorParameters(t *testing.T) {
	turm, err := NewTurm(1, 3)
	assert.Nil(t, turm)
	assert.Error(t, err)

	turm, err = NewTurm(5, 1)
	assert.Nil(t, turm)
	assert.Error(t, err)
}

func TestTurm_CalculateIterative(t *testing.T) {
	turm, err := NewTurm(5, 3)
	assert.NoError(t, err)
	expectedResults := []TurmIntermediateResult{
		{OldValue: 5, Operation: '*', Operand: 2, NewValue: 10},
		{OldValue: 10, Operation: '*', Operand: 3, NewValue: 30},
		{OldValue: 30, Operation: '/', Operand: 2, NewValue: 15},
		{OldValue: 15, Operation: '/', Operand: 3, NewValue: 5},
	}

	iter := turm.CalculateIterative()
	results := make([]TurmIntermediateResult, 0)
	iter(func(r TurmIntermediateResult) bool {
		results = append(results, r)
		return true
	})

	assert.Equal(t, expectedResults, results)
}
