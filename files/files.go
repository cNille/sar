package files

import (
	"bufio"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

var MatchCount int

func IsHidden(path string) bool {
	// Check if file is hidden or in a hidden directory
	// Split path into parts and check each part
	pathParts := strings.Split(path, "/")
	for _, part := range pathParts {
		if part == "." || part == ".." {
			continue
		}
		if strings.HasPrefix(part, ".") {
			return true
		}
	}
	return false
}

func IsBinary(path string) (bool, error) {
	file, err := os.Open(path)
	if err != nil {
		return false, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanBytes)
	for scanner.Scan() {
		if scanner.Bytes()[0] == 0 {
			return true, nil
		}
	}

	if err := scanner.Err(); err != nil {
		return false, err
	}

	return false, nil
}

func ReadFile(path string) ([]byte, error) {
	return ioutil.ReadFile(path)
}

func WriteFile(path string, data []byte) error {
	return ioutil.WriteFile(path, data, 0)
}

func FileChanged(oldContents []byte, newContents []byte) bool {
	return string(oldContents) != string(newContents)
}

func WalkDirectory(path string, visit func(fp string, fi os.FileInfo, err error) error) error {
	return filepath.Walk(path, visit)
}
