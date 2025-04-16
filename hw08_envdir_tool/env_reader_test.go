package main

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestReadDir(t *testing.T) {
	t.Run("success: get env data", func(t *testing.T) {
		tmpDir := t.TempDir()
		testCases := []struct {
			name     string
			content  string
			isDir    bool
			expected EnvValue
		}{
			{"FOO", "test\nvalue", false, EnvValue{"test", false}},
			{"BAR", "", false, EnvValue{"", true}},
			{"EMPTY_LINE", "\n\n", false, EnvValue{"", false}},
			{"WITH_NULL", "first\x00second", false, EnvValue{"first\nsecond", false}},
			{"DIR_SKIPPED", "", true, EnvValue{"", false}},
		}

		for _, tc := range testCases {
			path := filepath.Join(tmpDir, tc.name)

			if tc.isDir {
				err := os.Mkdir(path, 0o755)
				require.NoError(t, err)
				continue
			}

			err := os.WriteFile(path, []byte(tc.content), 0o644)
			require.NoError(t, err)
		}

		env, err := ReadDir(tmpDir)
		require.NoError(t, err)

		_, dirExists := env["DIR_SKIPPED"]
		require.False(t, dirExists, "Directory should be skipped")

		for _, tc := range testCases {
			if tc.isDir {
				continue
			}

			val, exists := env[tc.name]
			require.True(t, exists, "Variable %s not found", tc.name)
			require.Equal(t, tc.expected, val, "For %s", tc.name)
		}
	})

	t.Run("invalid directory", func(t *testing.T) {
		_, err := ReadDir("/non/existent/path")
		require.Error(t, err)
	})
}
