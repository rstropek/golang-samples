package validator

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEmptyValidatorIsValid(t *testing.T) {
	v := New()
	assert.True(t, v.Valid())
}

func TestCheck(t *testing.T) {
	v := New()
	v.Check(false, "test", "testmessage")
	assert.False(t, v.Valid())
	assert.Equal(t, "testmessage", v.Errors["test"])
}

func TestIn(t *testing.T) {
	assert.True(t, In("a", []string{"a", "b"}...))
	assert.False(t, In("c", []string{"a", "b"}...))
}

func TestMatches(t *testing.T) {
	assert.True(t, Matches("john.doe@somewhere.com", EmailRX))
	assert.False(t, Matches("john.doe@", EmailRX))
}

func TestUnique(t *testing.T) {
	assert.True(t, Unique([]string{"a", "b"}))
	assert.False(t, Unique([]string{"a", "a"}))
}
