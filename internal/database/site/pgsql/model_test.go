package pgsql

import (
	"database/sql"
	"testing"
)

func TestNotNullString(t *testing.T) {
	var nulls sql.NullString
	var s = ""
	blankString := NotNullString(nulls)
	if blankString != s {
		t.Fatalf("Want %s, got %v", s, blankString)
	}
}
