package main

import "testing"

func TestMaskValue(t *testing.T) {
	cases := map[string]string{
		"":         "****",
		"ab":       "****",
		"abcd":     "****",
		"abcdef":   "ab**ef",
		"abcdefgh": "ab****gh",
	}
	for in, want := range cases {
		got := maskValue(in)
		if got != want {
			t.Fatalf("maskValue(%q)=%q want %q", in, got, want)
		}
	}
}
