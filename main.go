package main

import (
	"fmt"
	"os"

	"github.com/cnille/sar/flags"
	"github.com/cnille/sar/searchandreplace"
)

func main() {
	searchArgs, err := flags.Parse()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	if len(searchArgs) < 3 {
		flags.Usage()
		os.Exit(1)
	}

	if searchArgs[0] == searchArgs[1] {
		fmt.Println("The search string and the replace string are identical. Nothing to do.")
		os.Exit(1)
	}

	matchesFound, err := searchandreplace.Preview(searchArgs)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	if !matchesFound {
		os.Exit(0)
	}

	if !flags.Force {
		if !searchandreplace.ConfirmProceed() {
			fmt.Println("Skipping...")
			os.Exit(0)
		}
	}

	err = searchandreplace.Replace(searchArgs)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
