package main

import (
	"fmt"
	"os"
	"os/exec"
	"syscall"
)

// RunCmd runs a command + arguments (cmd) with environment variables from env.
func RunCmd(cmd []string, env Environment) (returnCode int) {
	for key, value := range env {
		if value.NeedRemove {
			err := syscall.Unsetenv(key)
			if err != nil {
				return 1
			}
		} else {
			err := syscall.Setenv(key, value.Value)
			if err != nil {
				return 1
			}
		}
	}

	command := cmd[0]
	args := cmd[1:]

	commandExec := exec.Command(command, args...)

	commandExec.Stdin = os.Stdin
	commandExec.Stdout = os.Stdout
	commandExec.Stderr = os.Stderr

	err := commandExec.Run()
	if err != nil {
		fmt.Println("command:", command, "error:", err)
		return 1
	}

	return commandExec.ProcessState.ExitCode()
}
