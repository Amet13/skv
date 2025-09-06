package main

import "testing"

func TestIndexOf(t *testing.T) {
	if i := indexOf([]string{"a", "b", "c"}, "b"); i != 1 {
		t.Fatalf("want 1 got %d", i)
	}
	if i := indexOf([]string{"a"}, "x"); i != -1 {
		t.Fatalf("want -1 got %d", i)
	}
}

func TestSplitCSVAndTrim(t *testing.T) {
	parts := splitCSV(" a, b , ,c ")
	if len(parts) != 3 || parts[0] != "a" || parts[1] != "b" || parts[2] != "c" {
		t.Fatalf("unexpected parts: %#v", parts)
	}
	if got := trim("  hello\t"); got != "hello" {
		t.Fatalf("trim got %q", got)
	}
}
