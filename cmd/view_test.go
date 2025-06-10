package main

import "testing"

func TestViewRequiresArg(t *testing.T) {
	cmd := newViewCmd()
	cmd.SetArgs([]string{})
	if err := cmd.Execute(); err == nil {
		t.Fatal("expected error when no argument provided")
	}
}
