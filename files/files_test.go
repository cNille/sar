package files

import (
	"io/ioutil"
	"os"
	"testing"
)

func TestIsHidden(t *testing.T) {
	// Test that it correctly identifies hidden files/directories
	if !IsHidden("./.hidden") {
		t.Error("Expected './.hidden' to be identified as a hidden path")
	}
	if !IsHidden("./normal/.hidden") {
		t.Error("Expected './normal/.hidden' to be identified as a hidden path")
	}
	// Test that it correctly identifies non-hidden files/directories
	if IsHidden("./normal") {
		t.Error("Expected './normal' to not be identified as a hidden path")
	}

	// Test that it correctly identifies hidden files/directories with parent directories
	if IsHidden("../normal") {
		t.Error("Expected './normal' to not be identified as a hidden path")
	}
}

func TestIsBinary(t *testing.T) {
	// Create tmp binary file
	tmpFile, err := ioutil.TempFile("", "tmp")
	if err != nil {
		t.Error(err)
	}
	defer os.Remove(tmpFile.Name())

	// Write binary data to tmp file
	_, err = tmpFile.Write([]byte{0x00, 0x01, 0x02})
	if err != nil {
		t.Error(err)
	}

	// Test that it correctly identifies binary files
	isBin, err := IsBinary(tmpFile.Name())
	if err != nil {
		t.Error(err)
	}
	if !isBin {
		t.Error("Expected tmp file to be identified as binary")
	}
}

func TestFileChanged(t *testing.T) {
	// Test that it correctly identifies changed files
	if !FileChanged([]byte("old"), []byte("new")) {
		t.Error("Expected file to be identified as changed")
	}

	// Test that it correctly identifies unchanged files
	if FileChanged([]byte("old"), []byte("old")) {
		t.Error("Expected file to not be identified as changed")
	}
}
