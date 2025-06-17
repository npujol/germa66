package utils_test

import (
	"testing"

	"germa66/internal/utils"

	"github.com/stretchr/testify/assert"
)

func TestGetPathName(t *testing.T) {
	tests := []struct {
		name     string
		path     string
		expected string
	}{
		{
			name:     "unix path with extension",
			path:     "/path/to/file.txt",
			expected: "file",
		},
		{
			name:     "windows path with extension",
			path:     "C:\\path\\to\\file.txt",
			expected: "file",
		},
		{
			name:     "file without extension",
			path:     "/path/to/file",
			expected: "file",
		},
		{
			name:     "file with multiple extensions",
			path:     "/path/to/file.tar.gz",
			expected: "file.tar",
		},
		{
			name:     "just filename",
			path:     "file.txt",
			expected: "file",
		},
		{
			name:     "hidden file",
			path:     "/path/to/.hidden",
			expected: ".hidden",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := utils.GetPathName(tt.path)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestChangePathExt(t *testing.T) {
	tests := []struct {
		name     string
		path     string
		newExt   string
		expected string
	}{
		{
			name:     "change txt to csv",
			path:     "/path/to/file.txt",
			newExt:   ".csv",
			expected: "/path/to/file.csv",
		},
		{
			name:     "add extension to file without extension",
			path:     "/path/to/file",
			newExt:   ".txt",
			expected: "/path/to/file.txt",
		},
		{
			name:     "change extension without leading dot",
			path:     "/path/to/file.txt",
			newExt:   "csv",
			expected: "/path/to/file.csv",
		},
		{
			name:     "multiple extensions",
			path:     "/path/to/file.tar.gz",
			newExt:   ".zip",
			expected: "/path/to/file.tar.zip",
		},
		{
			name:     "just filename",
			path:     "file.txt",
			newExt:   ".csv",
			expected: "file.csv",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := utils.ChangePathExt(tt.path, tt.newExt)
			assert.Equal(t, tt.expected, result)
		})
	}
}