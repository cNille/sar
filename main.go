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

// init is a special function that is called before main
func init() {
	flag.BoolVar(&force, "force", false, "Force mode. If set, the script will make the replacements without asking for confirmation.")
	flag.BoolVar(&verbose, "verbose", false, "Verbose mode. If set, the script will print out all files it is processing.")

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
		fmt.Println("  -extension: Only process files with the given extension. Can be specified multiple times.")
		fmt.Println("\nExamples:")
		fmt.Println("\tsar \"my-app\" \"my-new-app\" ./my-folder -extension .go -extension .html")
		fmt.Println("\tsar -force \"my-app\" \"my-new-app\" ./my-folder")
		fmt.Println("\tsar -verbose \"my-app\" \"my-new-app\" ./my-folder")
		os.Exit(1)
	}

	search := flag.Args()[0]
	replace := flag.Args()[1]
	startDir := flag.Args()[2]

	err := filepath.Walk(startDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			fmt.Printf("Error with file path %q: %v\n", path, err)
			return err
		}

		if !info.IsDir() {
			isBin, err := isBinary(path)
			if err != nil {
				fmt.Printf("Error reading file %q: %v\n", path, err)
				return err
			}

			if isBin {
				if verbose {
					fmt.Printf("Skipping binary file: %s\n", path)
				}
			} else {
				err := replaceInFile(path, search, replace)
				if err != nil {
					fmt.Printf("Error replacing in file %q: %v\n", path, err)
					return err
				}
			}
		}

		return nil
	})

	if err != nil {
		fmt.Printf("Error walking the path %v: %v\n", startDir, err)
		return
	}
}

func writeToFile(path string, content string) error {
	err := ioutil.WriteFile(path, []byte(content), 0)
	return err
}

func replaceInFile(path, search, replace string) error {
	read, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}

	newContents := strings.Replace(string(read), search, replace, -1)
	if newContents != string(read) {
		fmt.Printf("%s%s%s\n", blueStart, path, resetColor)
		oldLines := strings.Split(string(read), "\n")
		newLines := strings.Split(newContents, "\n")
		for i, line := range oldLines {
			if line != newLines[i] {
				fmt.Println(strings.Replace(line, search, redStart+search+resetColor, -1))
			}
		}

		if force {
			err = writeToFile(path, newContents)
			return err
		}

		reader := bufio.NewReader(os.Stdin)
		fmt.Print("Do you want to proceed with the above changes? (y/n): ")
		text, _ := reader.ReadString('\n')
		if strings.TrimSpace(strings.ToLower(text)) == "y" {
			err = writeToFile(path, newContents)
			if err != nil {
				return err
			}
		} else {
			fmt.Println("Skipping...")
		}
	}

	return nil
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
