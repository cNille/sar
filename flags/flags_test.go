package flags

import (
	"os"
	"testing"
)

func TestParse(t *testing.T) {
	os.Args = []string{"cmd", "-force", "-verbose", "-include-hidden", "-extension", ".txt", "-extension", ".go", "search-string", "replace-string", "./test-dir"}
	flagArgs, err := Parse()
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	if len(flagArgs) != 3 || !Force || !Verbose || !IncludeHidden {
		t.Errorf("Incorrect flag parsing")
	}
}
