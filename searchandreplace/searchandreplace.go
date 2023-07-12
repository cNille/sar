package searchandreplace

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/cnille/sar/files"
	"github.com/cnille/sar/flags"
)

var blueStart = "\033[1;34m"
var resetColor = "\033[0m"
var redStart = "\033[1;31m"

func Preview(searchArgs []string) (bool, error) {
	files.MatchCount = 0
	err := files.WalkDirectory(searchArgs[2], func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		return parseFile(path, info, searchArgs[0], searchArgs[1], false)
	})
	if err != nil {
		return false, err
	}

	matchesFound := files.MatchCount > 0
	if matchesFound {
		fmt.Printf("Preview complete. %d files would be changed.\n", files.MatchCount)
	} else {
		fmt.Println("No files would be changed.")
	}

	return matchesFound, nil
}

func Replace(searchArgs []string) error {
	files.MatchCount = 0
	err := files.WalkDirectory(searchArgs[2], func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		return parseFile(path, info, searchArgs[0], searchArgs[1], true)
	})

	if err != nil {
		return err
	}

	if files.MatchCount == 0 {
		fmt.Println("No files were changed.")
	} else {
		fmt.Printf("Write complete. %d files have been changed.\n", files.MatchCount)
	}

	return nil
}

func ConfirmProceed() bool {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Do you want to proceed with the above changes? (y/n): ")
	text, _ := reader.ReadString('\n')
	return strings.TrimSpace(strings.ToLower(text)) == "y"
}

func replaceInFile(path, search, replace string, write bool) (bool, error) {
	read, err := files.ReadFile(path)
	if err != nil {
		return false, err
	}

	newContents := strings.Replace(string(read), search, replace, -1)
	if files.FileChanged(read, []byte(newContents)) {
		if write {
			err := files.WriteFile(path, []byte(newContents))
			if err != nil {
				return false, err
			}
			files.MatchCount++
		} else {
			fmt.Printf("%s%s%s\n", blueStart, path, resetColor)
			oldLines := strings.Split(string(read), "\n")
			newLines := strings.Split(newContents, "\n")
			matchFound := false
			for i, line := range oldLines {
				if line != newLines[i] {
					fmt.Println(strings.Replace(line, search, redStart+search+resetColor, -1))
					matchFound = true
				}
			}
			if matchFound {
				files.MatchCount++
			}
		}
		return true, nil
	}
	return false, nil
}

func parseFile(path string, info os.FileInfo, search string, replace string, write bool) error {
	if !info.IsDir() {
		isBin, err := files.IsBinary(path)
		if err != nil {
			return err
		}

		if !flags.IncludeHidden && files.IsHidden(path) {
			return nil
		}

		if !files.ExtensionMatches(path, flags.Extensions) {
			return nil
		}

		if isBin {
			if flags.Verbose {
				fmt.Printf("Skipping binary file: %s\n", path)
			}
		} else {
			_, err := replaceInFile(path, search, replace, write)
			if err != nil {
				return err
			}
		}
	}
	return nil
}
