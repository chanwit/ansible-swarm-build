package main

import "testing"
import "github.com/stretchr/testify/assert"

func TestMe(t *testing.T) {
	// assert.NoError(t, Import("files/ubuntu-14.10-server-amd64.ova", 512))
	// assert.NoError(t, Clone("base", "box", 4, "vboxnet0"))
	assert.NoError(t, Remove("box001", "box002", "box003", "box004", "base"))
}
