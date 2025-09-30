//go:build !zos
// +build !zos

package pty

import (
	"testing"
)

func getNonBlockingFile(t *testing.T, ptmx Pty, path string) Pty {
	t.Helper()
	return ptmx
}
