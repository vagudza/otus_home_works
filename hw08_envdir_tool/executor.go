package main

import (
	"errors"
	"os"
	"os/exec"
	"strings"
)

// RunCmd runs a command + arguments (cmd) with environment variables from env.
func RunCmd(cmd []string, env Environment) (returnCode int) {
	//nolint:gosec // it is OK
	command := exec.Command(cmd[0], cmd[1:]...)
	command.Stdin = os.Stdin
	command.Stdout = os.Stdout
	command.Stderr = os.Stderr
	command.Env = prepareEnv(env)

	err := command.Run()
	if err != nil {
		var exitError *exec.ExitError
		if errors.As(err, &exitError) {
			// https://pkg.go.dev/os/exec#ExitError
			return exitError.ExitCode()
		}

		// 1 is default exit code for any error
		return 1
	}

	// 0 is default exit code for success
	return 0
}

// prepareEnv prepares environment variables for the command.
func prepareEnv(env Environment) []string {
	currentEnv := os.Environ() // ["FOO=bar", "BAZ=qux"]
	newEnvMap := make(map[string]string, len(currentEnv))

	for _, e := range currentEnv {
		// Ex. parts=["FOO", "bar123"]; parts[0] is env variable name, parts[1] - env variable value
		parts := strings.SplitN(e, "=", 2)
		if env[parts[0]].NeedRemove {
			continue
		}

		newEnvMap[parts[0]] = parts[1]
	}

	// Add new env variables
	for name, val := range env {
		if !val.NeedRemove {
			newEnvMap[name] = val.Value
		}
	}

	newEnv := make([]string, 0, len(newEnvMap))
	for name, val := range newEnvMap {
		newEnv = append(newEnv, name+"="+val)
	}

	return newEnv
}
