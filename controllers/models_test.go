package controllers

import (
	"strings"
	"testing"
)

func TestModelName(t *testing.T) {
	cases := map[string]string{
		"x-wing.glb":                     "x-wing",
		"  Falcon.GLB ":                  "Falcon",
		`C:\models\tie.glb`:              "tie",
		"../../etc/passwd":               "passwd",
		"":                               "",
		strings.Repeat("a", 80) + ".glb": strings.Repeat("a", 60),
	}
	for input, want := range cases {
		if got := modelName(input); got != want {
			t.Errorf("modelName(%q) = %q, want %q", input, got, want)
		}
	}
}
