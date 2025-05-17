package main

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestRunCmd(t *testing.T) {
	t.Run("empty cmd", func(t *testing.T) {
		env := Environment{}
		var command []string

		result := RunCmd(command, env)
		assert.Equal(t, 1, result)
	})
}
