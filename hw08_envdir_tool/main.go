package main

import (
	"fmt"
	"os"
)

func main() {
	args := os.Args[1:]

	if len(args) < 2 {
		fmt.Println("Usage: go-envdir <path> <command> <args>")
		return
	}

	directory := args[0]
	command := args[1:]

	env, err := ReadDir(directory)
	if err != nil {
		fmt.Println(err)
		return
	}

	exitCode := RunCmd(command, env)

	os.Exit(exitCode)
}
