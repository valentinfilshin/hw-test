package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRunCmd(t *testing.T) {
	t.Run("empty cmd", func(t *testing.T) {
		env := Environment{}
		var command []string

		result := RunCmd(command, env)
		assert.Equal(t, 1, result)
	})

	t.Run("wrong command", func(t *testing.T) {
		env := Environment{}
		command := []string{"wrong command", "hi"}

		result := RunCmd(command, env)
		assert.Equal(t, 1, result)
	})

	t.Run("good command", func(t *testing.T) {
		env := Environment{}
		command := []string{"ls", "-la"}

		result := RunCmd(command, env)
		assert.Equal(t, 0, result)
	})
}
