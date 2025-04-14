package main

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
)

type Environment map[string]EnvValue

// EnvValue helps to distinguish between empty files and files with the first empty line.
type EnvValue struct {
	Value      string
	NeedRemove bool
}

// ReadDir reads a specified directory and returns map of env variables.
// Variables represented as files where filename is name of variable, file first line is a value.
func ReadDir(dir string) (Environment, error) {
	env := make(Environment)

	files, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	for _, f := range files {
		if f.IsDir() {
			continue
		}

		// Skip files with '=' in their name to avoid conflicts with environment variable names (FOO=BAR=value)
		name := f.Name()
		if strings.Contains(name, "=") {
			continue
		}

		content, err := os.ReadFile(filepath.Join(dir, name))
		if err != nil {
			return nil, err
		}

		// split the content by new line and take the first line
		lines := bytes.SplitN(content, []byte("\n"), 2)
		// replace all null bytes in first line with new line
		firstLine := bytes.ReplaceAll(lines[0], []byte{0x00}, []byte("\n"))
		// clear all trailing spaces and tabs
		value := strings.TrimRight(string(firstLine), " \t")

		env[name] = EnvValue{
			Value:      value,
			NeedRemove: len(content) == 0,
		}
	}

	return env, nil
}
