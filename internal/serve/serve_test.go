package serve

import (
	"os"
	"testing"
)

func TestHasErrors(t *testing.T) {
	noErrors := []Result{
		{Host: "a", Status: "started"},
		{Host: "b", Status: "skipped (already running)"},
	}
	if HasErrors(noErrors) {
		t.Error("expected no errors")
	}

	withErrors := []Result{
		{Host: "a", Status: "started"},
		{Host: "b", Status: "fail", Err: os.ErrNotExist},
	}
	if !HasErrors(withErrors) {
		t.Error("expected errors")
	}
}
