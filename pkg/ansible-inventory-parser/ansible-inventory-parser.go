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
	Name      string
	GroupName string
	Vars      map[string]string
	Indent    int
}

type Group struct {
	Name     string
	Hosts    []Host
	Children []Group
	Vars     map[string]string
	Indent   int
}

func readContent(filepath string) ([]string, error) {
	// Open the file.
	file, err := os.Open(filepath)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)

	var res []string
	for scanner.Scan() {
		res = append(res, scanner.Text())
	}

	return res, nil
}

func Parse(filepath string) error {
	inventoryFile, err := readContent(filepath)
	if err != nil {
		return err
	}
	for i := 0; i < len(inventoryFile); i++ {
		if isGroup(inventoryFile, i) {
			fmt.Printf("group name: %s\n", inventoryFile[i])
		}
		var g Group
		getHosts(inventoryFile, i, g)

	}
	return nil

}

// getHosts checks if the host is on the host level.
func getHosts(inventoryFile []string, i int, g Group) []Host {
	if !(strings.TrimSpace(inventoryFile[i]) == "hosts:") {
		return nil
	}
	// the line is the host file.
	hostIndent := indentationLevel(inventoryFile[i]) + 1

	// get all the hosts of this group.
	for j := i + 1; j < len(inventoryFile); j++ {
		indentLvl := indentationLevel(inventoryFile[j])

		// if it's the children then is the start of a sub group.
		if strings.Contains(inventoryFile[j], "children:") {
			getHosts(inventoryFile, j, g)
		}
		if indentLvl != hostIndent {
			break
		}

		// if we got heat then it's a host
		varsIndent := indentLvl + 1
		host := Host{
			Name:   strings.Split(inventoryFile[j], ":")[0],
			Indent: indentLvl,
			Vars:   make(map[string]string),
		}
		// get the vars
	inner:
		for k := j + 1; k < len(inventoryFile); k++ {
			indent := indentationLevel(inventoryFile[k])
			if indent != varsIndent {
				break inner
			}
			row := strings.TrimSpace(inventoryFile[k])

			s := strings.Split(row, ":")
			if len(s) == 2 {
				host.Vars[strings.Split(s[0], ":")[0]] = strings.Split(s[1], ":")[0]
			}
			// skip this lines because we already got the vars from them.
			j++
		}
		fmt.Printf("host %+v\n", host)

	}

	return nil
}

// isGroup test if the next line is the host row, meaning it's the start of a new group.
func isGroup(inventoryFile []string, i int) bool {
	for j := i + 1; j < len(inventoryFile); j++ {
		row := inventoryFile[j]
		row = strings.TrimSpace(row)
		// if the len is 1 then it's not a comment line.
		if len(strings.Split(row, "#")) == 1 {
			return row == "hosts:"
		}

	}
	return false

}

// hostDetails returns the details of the host.
func hostDetails(indentLvl int, scanner bufio.Scanner, row, groupName string) Host {
	host := Host{
		Name:      strings.Split(row, ":")[0],
		Indent:    indentLvl,
		GroupName: groupName,
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
	indentLvl := indentationLevel(nextRow) + 1
	nextRow = strings.TrimSpace(nextRow)
	nextRow = strings.Split(nextRow, "#")[0]

	if nextRow == "hosts:" {
		groups[indentLvl] = strings.Split(row, ":")[0]
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
