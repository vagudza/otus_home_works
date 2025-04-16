package main

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestRunCmd(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		env := Environment{
			"TEST_VAR": EnvValue{"test_value", false},
		}
		code := RunCmd([]string{"sh", "-c", "exit 0"}, env)
		if code != 0 {
			t.Errorf("Expected exit code 0, got %d", code)
		}
	})

	t.Run("failure", func(t *testing.T) {
		env := Environment{}
		code := RunCmd([]string{"sh", "-c", "exit 42"}, env)
		if code != 42 {
			t.Errorf("Expected exit code 42, got %d", code)
		}
	})

	t.Run("success: check env var replacement", func(t *testing.T) {
		env := Environment{
			"FOO": EnvValue{"bar", false},
			"BAR": EnvValue{"", true},
		}

		// Set existing variable to test removal
		err := os.Setenv("BAR", "should_be_removed")
		require.NoError(t, err)

		defer func() {
			err = os.Unsetenv("BAR")
			require.NoError(t, err)
		}()

		code := RunCmd([]string{"sh", "-c", "echo $FOO; echo $BAR"}, env)
		if code != 0 {
			t.Errorf("Command failed with code %d", code)
		}
	})
}
