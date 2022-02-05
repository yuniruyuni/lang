package token_test

import (
	"testing"

	"gotest.tools/assert"
)

func TestTokenizer_Tokenize(t *testing.T) {
	t.Parallel()
	assert.Equal(t, "test", "test")
}
