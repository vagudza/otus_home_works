package main

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCopy(t *testing.T) {
	testDir := t.TempDir()
	inputPath := filepath.Join(testDir, "input.txt")

	err := os.WriteFile(inputPath, []byte("123456"), 0o644)
	require.NoError(t, err)

	dirPath := filepath.Join(testDir, "testdir")
	err = os.Mkdir(dirPath, 0o755)
	require.NoError(t, err)

	t.Run("invalid parameters", func(t *testing.T) {
		tests := []struct {
			name          string
			fromPath      string
			toPath        string
			offset        int64
			limit         int64
			expectedError error
		}{
			{
				name:          "non-existent source file",
				fromPath:      filepath.Join(testDir, "non-existent.txt"),
				toPath:        filepath.Join(testDir, "out.txt"),
				offset:        0,
				limit:         0,
				expectedError: os.ErrNotExist,
			},
			{
				name:          "offset exceeds file size",
				fromPath:      inputPath,
				toPath:        filepath.Join(testDir, "out.txt"),
				offset:        1000000,
				limit:         0,
				expectedError: ErrOffsetExceedsFileSize,
			},
			{
				name:          "unsupported file type (directory)",
				fromPath:      dirPath,
				toPath:        filepath.Join(testDir, "out.txt"),
				offset:        0,
				limit:         0,
				expectedError: ErrUnsupportedFile,
			},
			{
				name:          "invalid destination path",
				fromPath:      inputPath,
				toPath:        "/non-existent-dir/out.txt",
				offset:        0,
				limit:         0,
				expectedError: os.ErrNotExist,
			},
		}

		for _, tc := range tests {
			t.Run(tc.name, func(t *testing.T) {
				err := Copy(tc.fromPath, tc.toPath, tc.offset, tc.limit)
				require.ErrorIs(t, err, tc.expectedError)
			})
		}
	})

	t.Run("successful copy operations", func(t *testing.T) {
		tests := []struct {
			name       string
			offset     int64
			limit      int64
			resultSize int64
		}{
			{
				name:       "copy entire file",
				offset:     0,
				limit:      0,
				resultSize: 6,
			},
			{
				name:       "copy with offset",
				offset:     3,
				limit:      0,
				resultSize: 3,
			},
			{
				name:       "copy with limit",
				offset:     0,
				limit:      4,
				resultSize: 4,
			},
			{
				name:       "copy with offset and limit",
				offset:     1,
				limit:      3,
				resultSize: 3,
			},
		}

		for _, tc := range tests {
			t.Run(tc.name, func(t *testing.T) {
				outPath := filepath.Join(testDir, "out.txt")
				err := Copy(inputPath, outPath, tc.offset, tc.limit)
				require.NoError(t, err)

				stat, err := os.Stat(outPath)
				require.NoError(t, err)
				require.Equal(t, tc.resultSize, stat.Size())

				err = os.Remove(outPath)
				require.NoError(t, err)
			})
		}
	})
}
