package main

import (
	"testing"
)

func TestRegExpMatch(t *testing.T) {
	text := `UserDB user's db
	@ifmeasure`
	if !checkComment(text) {
		t.Fatal("should match")
	}

}
