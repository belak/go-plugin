package plugin

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEnsure(t *testing.T) {
	t.Parallel()

	assert.NotPanics(t, func() {
		Ensure(nil)
	})

	assert.Panics(t, func() {
		Ensure(errors.New("Hello world"))
	})
}
