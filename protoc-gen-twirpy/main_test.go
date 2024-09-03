package main

import (
	"io"
	"os"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func loadProto(t *testing.T) *os.File {
	file, err := os.Open("test_data/protoc_input.data")
	require.NoError(t, err)
	return file
}

func loadGenerated(t *testing.T) []byte {
	file, err := os.Open("test_data/protoc_output.data")
	defer file.Close()
	require.NoError(t, err)
	result, err := io.ReadAll(file)
	require.NoError(t, err)
	return result
}

func replaceGenerated(t *testing.T, data []byte) {
	file, err := os.Create("test_data/protoc_output.data")
	defer file.Close()
	require.NoError(t, err)
	_, err = file.Write(data)
	require.NoError(t, err)
}

func captureOutput(t *testing.T, f func()) []byte {
	orig := os.Stdout
	r, w, err := os.Pipe()
	require.NoError(t, err)
	os.Stdout = w
	f()
	os.Stdout = orig
	require.NoError(t, w.Close())
	out, err := io.ReadAll(r)
	require.NoError(t, err)
	return out
}

func TestMainFn(t *testing.T) {
	os.Stdin = loadProto(t)
	result := captureOutput(t, main)
	if os.Getenv("TEST_DATA_UPDATE") == "1" {
		replaceGenerated(t, result)
	}
	assert.Equal(t, string(loadGenerated(t)), string(result))
}
