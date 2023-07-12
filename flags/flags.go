package flags

import (
	"flag"
	"fmt"
)

type StringSlice []string

func (s *StringSlice) String() string {
	return fmt.Sprintf("%v", *s)
}

func (s *StringSlice) Set(value string) error {
	*s = append(*s, value)
	return nil
}

var (
	Extensions    StringSlice
	Force         bool
	Verbose       bool
	IncludeHidden bool
)

func init() {
	flag.BoolVar(&Force, "force", false, "Force mode. If set, the script will make the replacements without asking for confirmation.")
	flag.BoolVar(&Verbose, "verbose", false, "Verbose mode. If set, the script will print out all files it is processing.")
	flag.BoolVar(&IncludeHidden, "include-hidden", false, "Include hidden files. If set, the script will process hidden files and directories.")

	// Handle extension flag "-extension"
	flag.Var(&Extensions, "extension", "Only process files with the given extension. Can be specified multiple times.")
}

func Parse() ([]string, error) {
	flag.Parse()

	if len(flag.Args()) < 3 {
		return nil, fmt.Errorf("not enough arguments provided")
	}

	return flag.Args(), nil
}

func Usage() {
	fmt.Println("Usage: sar [OPTIONAL FLAGS] \"search-string\"  \"new-string\" ./my-folder")
	fmt.Println("Note that the flags have to be specified before the search and replace strings.")
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
}
