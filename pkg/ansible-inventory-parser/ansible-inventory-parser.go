package ansibleinventoryparser

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func New() {

}

type All struct {
	Groups []Group
	Vars   map[string]string
}
type Host struct {
	Name   string
	Vars   map[string]string
	Indent int
}

type Group struct {
	Hosts     []Host
	SubGroups []Group
	Vars      map[string]string
	Indent    int
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

	// indicate the hosts indent level
	var hostIndentLvl int
	// Read and print the lines.
	for scanner.Scan() {
		row := scanner.Text()

		// get indentation level and remove trialing and leading spaces.
		indentLvl := indentationLevel(row)
		row = strings.TrimSpace(row)

		// get the inline comment and remove it.
		// inlineComment := InlineComment(row)
		row = strings.Split(row, "#")[0]

		// the hosts are on the level of the hosts key word plus 1.
		if row == "hosts:" {
			hostIndentLvl = indentLvl + 1
		}

		// update the group for the indentation level.
		updateGroup(groups, row, *scanner)

		if IsHost(indentLvl, hostIndentLvl) {
			host := getHostDetails(indentLvl, *scanner, row)
			fmt.Printf("%+v\n", host)
		}

		// fmt.Printf("indent level: %d, %s\n", indentLvl, row)
		// // fmt.Printf("inlineComment: %s, %s\n", inlineComment, row)
		// fmt.Printf("group: %s, %s\n", groups[indentLvl], row)
		// fmt.Printf("host group: %s, %s\n", groups[hostIndentLvl], row)

	}

	if err := scanner.Err(); err != nil {
		return err
	}
	// fmt.Printf("%+v\n", groups)

	return nil

}

// checks if the host is on the host level.
func IsHost(indentLvl, hostIndentLvl int) bool {
	return indentLvl == hostIndentLvl

}

func getHostDetails(indentLvl int, scanner bufio.Scanner, row string) Host {
	host := Host{
		Name:   strings.Split(row, ":")[0],
		Indent: indentLvl,
	}
	host.Vars = make(map[string]string)
	for scanner.Scan() {
		row := scanner.Text()

		// get indentation level and remove trialing and leading spaces.
		rowIndentLvl := indentationLevel(row)

		// TODO: check if i need to support more types other then key value.
		// check if the level of the indentation is different from 1 level above the indentation of the host.
		if rowIndentLvl != indentLvl+1 {
			break
		}
		row = strings.TrimSpace(row)

		s := strings.Split(row, ":")
		if len(s) == 2 {
			host.Vars[strings.Split(s[0], ":")[0]] = strings.Split(s[1], ":")[0]
		}
	}

	return host

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
	indetLvl := indentationLevel(nextRow) + 1
	nextRow = strings.TrimSpace(nextRow)
	nextRow = strings.Split(nextRow, "#")[0]

	if nextRow == "hosts:" {
		groups[indetLvl] = strings.Split(row, ":")[0]
	}

}

// isKeyVal checks if it's a key val row or only key.
func isKeyVal(row string) bool {
	s := strings.Split(row, ":")
	// if the len is two then it's a key value row.
	return len(s) == 2 && s[1] != ""
}

// isHeadComment checks if the row is a comment.
func isHeadComment(row string) bool {
	return strings.HasPrefix(row, "#")
}

// indentationLevel returns the indentation level of the row.
func indentationLevel(row string) int {
	return (len(row) - len(strings.TrimLeft(row, " "))) / 2
}

// inlineComment will return the inline comment.
func inlineComment(row string) string {
	s := strings.Split(row, "#")
	if len(s) == 2 {
		return s[1]
	}
	return ""

}
