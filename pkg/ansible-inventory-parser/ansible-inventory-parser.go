package ansibleinventoryparser

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func New() {

}

type Inventory struct {
	All
}

type All struct {
}

func Parse(filepath string) error {
	// Open the file.
	file, err := os.Open(filepath)
	if err != nil {
		return err
	}
	defer file.Close()

	// Set up the scanner.
	scanner := bufio.NewScanner(file)
	groups := make(map[int]string)

	// Read and print the lines.
	for scanner.Scan() {
		row := scanner.Text()

		// get indentation level and remove trialing and leading spaces.
		indetLvl := IndentationLevel(row)
		row = strings.TrimSpace(row)

		// get the inline comment and remove it.
		// inlineComment := InlineComment(row)
		row = strings.Split(row, "#")[0]

		// update the group for the indentation level.
		updateGroup(groups, row, *scanner)

		fmt.Printf("indent level: %d, %s\n", indetLvl, row)
		// fmt.Printf("inlineComment: %s, %s\n", inlineComment, row)
		fmt.Printf("group: %s, %s\n", groups[indetLvl], row)

		if IsKeyVal(row) {
			fmt.Println("key", row)
		}
	}

	if err := scanner.Err(); err != nil {
		return err
	}
	fmt.Printf("%+v\n", groups)

	return nil

}

// updateGroup will update the group for this indentation level.
func updateGroup(groups map[int]string, row string, scanner bufio.Scanner) {
	// if there is no next row.
	if !scanner.Scan() {
		return
	}
	nextRow := scanner.Text()

	// +1 because we want to update the row after the "hosts" line.
	// read the next line to check if the it's the start of a new group.
	indetLvl := IndentationLevel(nextRow) + 1
	nextRow = strings.TrimSpace(nextRow)
	nextRow = strings.Split(nextRow, "#")[0]

	if nextRow == "hosts:" {
		groups[indetLvl] = strings.Split(row, ":")[0]
	}

}

// IsKeyVal checks if it's a key val row or only key.
func IsKeyVal(row string) bool {
	s := strings.Split(row, ":")
	// if the len is two then it's a key value row.
	return len(s) == 2 && s[1] != ""
}

// IsHeadComment checks if the row is a comment.
func IsHeadComment(row string) bool {
	return strings.HasPrefix(row, "#")
}

// IndentationLevel returns the indentation level of the row.
func IndentationLevel(row string) int {
	return (len(row) - len(strings.TrimLeft(row, " "))) / 2
}

// InlineComment will return the inline comment.
func InlineComment(row string) string {
	s := strings.Split(row, "#")
	if len(s) == 2 {
		return s[1]
	}
	return ""

}
