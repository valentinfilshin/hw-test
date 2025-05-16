package main

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestReadDir(t *testing.T) {
	expectedEnv := Environment{
		"BAR":    {Value: "bar", NeedRemove: false},
		"BARUHA": {Value: "baruha u uha!", NeedRemove: false},
		"EMPTY":  {NeedRemove: false},
		"FOO":    {Value: "   foo\nwith new line", NeedRemove: false},
		"HELLO":  {Value: "\"hello\"", NeedRemove: false},
		"UNSET":  {NeedRemove: true},
	}

	t.Run("check test cases", func(t *testing.T) {
		testDir := "testdata/env"
		testDirEnv, err := ReadDir(testDir)

		require.Nil(t, err)
		require.Equal(t, expectedEnv, testDirEnv)
	})

	t.Run("no directory", func(t *testing.T) {
		_, err := ReadDir("test")
		require.NotNil(t, err)
	})

	t.Run("wrong directory", func(t *testing.T) {
		_, err := ReadDir("/dev/random")
		require.NotNil(t, err)
	})
}
