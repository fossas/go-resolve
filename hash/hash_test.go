package hash_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/fossas/go-resolve/hash"
)

func TestHash(t *testing.T) {
	h, err := hash.Hash("testdata", []string{"example.go"})
	assert.NoError(t, err)
	assert.Equal(t, "MEaB+6yZ22ULP4JlU1dtsqlz+h4ZrwJhAkVeVgEIVik=", h)
}
