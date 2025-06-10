package main

import (
	"testing"
	"time"
)

func TestBuildSearchQuery(t *testing.T) {
	from := parseDate(t, "2023-01-01")
	to := parseDate(t, "2023-01-31")
	q := buildSearchQuery(&from, &to, "alice", "cli/cli", []string{"bug"})
	expected := "created:2023-01-01..2023-01-31 author:alice repo:cli/cli bug"
	if q != expected {
		t.Fatalf("expected %q, got %q", expected, q)
	}
}

func parseDate(t *testing.T, v string) time.Time {
	t.Helper()
	dt, err := time.Parse("2006-01-02", v)
	if err != nil {
		t.Fatal(err)
	}
	return dt
}
