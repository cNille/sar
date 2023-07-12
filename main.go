package main

import (
	"bufio"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

// Color codes
var blueStart = "\033[1;34m"
var resetColor = "\033[0m"
var redStart = "\033[1;31m"

// StringSlice is a custom flag type that allows multiple values to be specified
type StringSlice []string

func (s *StringSlice) String() string {
	return fmt.Sprintf("%v", *s)
}
func (s *StringSlice) Set(value string) error {
	*s = append(*s, value)
	return nil
}

// Flags
var extensions StringSlice
var force bool
var verbose bool
var includeHidden bool

// init is a special function that is called before main
func init() {
	flag.BoolVar(&force, "force", false, "Force mode. If set, the script will make the replacements without asking for confirmation.")
	flag.BoolVar(&verbose, "verbose", false, "Verbose mode. If set, the script will print out all files it is processing.")
	flag.BoolVar(&includeHidden, "include-hidden", false, "Include hidden files. If set, the script will process hidden files and directories.")

	// Handle extension flag "-extension"
	flag.Var(&extensions, "extension", "Only process files with the given extension. Can be specified multiple times.")
}

func main() {
	flag.Parse()

	if len(flag.Args()) < 3 {
		fmt.Println("Usage: sar \"search-string\"  \"new-string\" ./my-folder")
		fmt.Println("\nOptional flags:")
		fmt.Println("  -force: Force mode. If set, the script will make the replacements without asking for confirmation.")
		fmt.Println("  -verbose: Verbose mode. If set, the script will print out all files it is processing.")
		fmt.Println("  -include-hidden: Include hidden files. If set, the script will process hidden files and directories.")
		fmt.Println("  -extension: Only process files with the given extension. Can be specified multiple times.")
		fmt.Println("\nExamples:")
		fmt.Println("\tsar \"my-app\" \"my-new-app\" ./my-folder -extension .go -extension .html")
		fmt.Println("\tsar -force \"my-app\" \"my-new-app\" ./my-folder")
		fmt.Println("\tsar -verbose \"my-app\" \"my-new-app\" ./my-folder")
		fmt.Println("\tsar -include-hidden \"my-app\" \"my-new-app\" ./my-folder")
		os.Exit(1)
	}

	search := flag.Args()[0]
	replace := flag.Args()[1]
	startDir := flag.Args()[2]
	writeCount := 0

	// Preview
	err := filepath.Walk(startDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			fmt.Printf("Error with file path %q: %v\n", path, err)
			return err
		}

		writeCount, err = parseFile(path, info, search, replace, force)
		if err != nil {
			return err
		}
		return nil
	})

	// Return if writes are done
	if force {
		if writeCount == 0 {
			fmt.Println("No files were changed.")
		} else {
			fmt.Printf("Write complete. %d files have been changed.\n", writeCount)
		}
		return
	}

	// Return if no matches found
	if writeCount == 0 {
		fmt.Println("No files would be changed.")
		return
	}

	// Preview complete,
	fmt.Printf("Preview complete. %d files would be changed.\n", writeCount)

	// Prompt for confirmation
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Do you want to proceed with the above changes? (y/n): ")
	text, _ := reader.ReadString('\n')
	if strings.TrimSpace(strings.ToLower(text)) != "y" {
		fmt.Println("Skipping...")
		os.Exit(0)
	}

	// Replace
	err = filepath.Walk(startDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			fmt.Printf("Error with file path %q: %v\n", path, err)
			return err
		}

		writeCount, err = parseFile(path, info, search, replace, true)
		if err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		fmt.Printf("Error walking the path %v: %v\n", startDir, err)
		return
	}
}

func parseFile(path string, info os.FileInfo, search string, replace string, write bool) (int, error) {
	writeCount := 0
	if !info.IsDir() {
		isBin, err := isBinary(path)
		if err != nil {
			fmt.Printf("Error reading file %q: %v\n", path, err)
			return writeCount, err
		}

		if !includeHidden {
			// Check if file is hidden or in a hidden directory
			// Split path into parts and check each part
			pathParts := strings.Split(path, "/")
			for _, part := range pathParts {
				if strings.HasPrefix(part, ".") {
					return writeCount, nil
				}
			}
		}

		if isBin {
			if verbose {
				fmt.Printf("Skipping binary file: %s\n", path)
			}
		} else {
			fileChanged, err := replaceInFile(path, search, replace, write)
			if err != nil {
				fmt.Printf("Error replacing in file %q: %v\n", path, err)
				return writeCount, err
			}
			if fileChanged {
				writeCount += 1
			}
		}
	}

	return writeCount, nil
}

func replaceInFile(path, search, replace string, write bool) (bool, error) {
	read, err := ioutil.ReadFile(path)
	if err != nil {
		return false, err
	}

	isChanged := false
	newContents := strings.Replace(string(read), search, replace, -1)
	if newContents != string(read) {
		if write {
			err := ioutil.WriteFile(path, []byte(newContents), 0)
			return isChanged, err
		} else {
			fmt.Printf("%s%s%s\n", blueStart, path, resetColor)
			oldLines := strings.Split(string(read), "\n")
			newLines := strings.Split(newContents, "\n")
			for i, line := range oldLines {
				if line != newLines[i] {
					fmt.Println(strings.Replace(line, search, redStart+search+resetColor, -1))
					isChanged = true
				}
			}
		}
	}
	return isChanged, nil
}

func isBinary(path string) (bool, error) {
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
