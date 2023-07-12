package searchandreplace

import (
	"io/ioutil"
	"os"
	"testing"
)

func TestReplaceInFile(t *testing.T) {
	// Create a temporary file
	file, err := ioutil.TempFile("", "")
	if err != nil {
		t.Fatalf("Failed to create temporary file: %v", err)
	}
	defer os.Remove(file.Name())

	// Write some content to the file
	_, err = file.WriteString("Hello, world!")
	if err != nil {
		t.Fatalf("Failed to write to temporary file: %v", err)
	}

	// Close the file
	err = file.Close()
	if err != nil {
		t.Fatalf("Failed to close temporary file: %v", err)
	}

	// Test replaceInFile
	changed, err := replaceInFile(file.Name(), "world", "golang", true)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	if !changed {
		t.Error("Expected file contents to change")
	}

	// Check the file's new content
	content, err := ioutil.ReadFile(file.Name())
	if err != nil {
		t.Fatalf("Failed to read temporary file: %v", err)
	}
	if string(content) != "Hello, golang!" {
		t.Errorf("Unexpected file content: %s", string(content))
	}
}
