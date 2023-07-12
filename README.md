# SAR - Search And Replace

SAR is a tiny command-line utility written in Go that allows you to search for a
string in all text files in a specified directory (and its subdirectories), and
replace it with another string.

## Features

- Search and replace strings in multiple files
- Preview mode to show where replacements would be made, without making the changes
- Skip preview mode with the force-flag
- Skip binary files automatically
- Ability to include hidden files. (default is to ignore hidden directories and files)
- Ability to specify multiple file extensions to limit the search

## Usage

Build the SAR utility with:

```bash
go build -o sar
```

Basic usage of the SAR utility:

```bash
./sar "search-string" "new-string" ./directory-path
```

To skip preview and change files directly mode:

```bash
./sar -force "search-string" "new-string" ./directory-path
```

To specify file extensions:

```bash
./sar -extension ".txt" -extension ".go" "search-string" "new-string" ./my-dir
```

To include hidden files:

```bash
./sar -include-hidden "search-string" "new-string" .
```

## Installation

After building the SAR utility, you can install it to your system's PATH for
easy access:

```bash
sudo cp ./sar /usr/local/bin/
```

Starting from Go 1.16, go install can be used with a path to a Go program
within a module, and it will build and install that program in your GOPATH/bin
or GOBIN directory. Use:

```bash
go install .
```

## Run tests

```bash
go test ./...
```

## Future Features

Here are some features planned for future versions of SAR:

- Undo Feature: SAR will save all changes it makes in a temporary file. If you
  regret the last changes, SAR will be able to revert the files to their
  previous states.
- Handle regex.
- Allow to "split" and only confirm one file-change at a time.
- Only use one argument for search and replace strings, like sed. Example: "s/search-str/new-str"
  - Would enable for more features than just substitute.
