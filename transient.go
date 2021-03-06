// Manages transient output path.
//
// Performs:
//  1. Output to temporal file to the directory of desired output file.
//  2. Atomically rename it to the final output.
package main

import (
	"os"
	"path/filepath"

	"gopkg.in/myesui/uuid.v1"
)

// TransientOutputPath holds path for committing/aborting the result.
type TransientOutputPath struct {
	Output     string
	TempOutput string
	done       bool
}

// NewTransientOutput creates a context for atomic renaming output.
func NewTransientOutput(path string) *TransientOutputPath {
	result := new(TransientOutputPath)
	result.Output = filepath.Clean(path)
	d := filepath.Dir(result.Output)
	id := uuid.NewV4()
	result.TempOutput = filepath.Join(d, "gb-"+id.String()+".tmp")
	result.done = false
	return result
}

// Commit commits the result.
func (t *TransientOutputPath) Commit() error {
	if !t.done {
		t.done = true
		return os.Rename(t.TempOutput, t.Output)
	}
	return nil
}

// Abort discards the transient output.
func (t *TransientOutputPath) Abort() error {
	if !t.done {
		t.done = true
		return os.Remove(t.TempOutput)
	}
	return nil
}

// Done returns true if operation is done (Committed or Aborted).
func (t *TransientOutputPath) Done() bool {
	return t.done
}
