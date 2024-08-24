package admin

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMakeToken(t *testing.T) {
	t.Parallel()

	adminUsername = "foo"
	adminPassword = "bar"

	token, err := makeToken()
	assert.Nil(t, err)

	matched, err := parseToken(token)
	assert.Nil(t, err)
	assert.True(t, matched)
}
