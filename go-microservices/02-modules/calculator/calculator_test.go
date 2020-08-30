package calculator

import "testing"

// Note: This samples show very basic unit tests without any supporting packges.
// Consider packages like *testify* in practice (https://pkg.go.dev/github.com/stretchr/testify/assert?tab=doc).

func TestAdd(t *testing.T) {
	if result := Add(21, 21); result != 42 {
		t.Errorf("Expected 42, got %d", result)
	}
}

func TestSub(t *testing.T) {
	if result := Sub(63, 21); result != 42 {
		t.Errorf("Expected 42, got %d", result)
	}
}

func TestDiv(t *testing.T) {
	if result, err := Div(84, 2); result != 42 || err != nil {
		t.Errorf("Expected 42, got %d", result)
	}
}

func TestDivByZero(t *testing.T) {
	if _, err := Div(42, 0); err == nil {
		t.Errorf("Expected error, got nil")
	}
}
