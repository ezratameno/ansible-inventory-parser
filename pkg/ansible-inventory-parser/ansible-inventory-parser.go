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

	var g Group
	runner := &g

	for i := 0; i < len(inventoryFile); i++ {
		hosts, groupName := getHosts(inventoryFile, i)
		if hosts != nil {
			runner.Name = groupName
			runner.Hosts = hosts
			fmt.Printf("runner: %+v\n", runner)
			// fmt.Printf("groupName: %+v\n", groupName)
		}

	}
	return nil

}

// getHosts checks if the host is on the host level.
func getHosts(inventoryFile []string, i int) ([]Host, string) {
	if !(strings.TrimSpace(inventoryFile[i]) == "hosts:") {
		return nil, ""
	}
	// the line is the host file.
	hostIndent := indentationLevel(inventoryFile[i]) + 1
	var hosts []Host
	var groupName string

	// ====================================================
	// get group name
	for j := i; j >= 0; j-- {
		if parseName(inventoryFile[j]) == "hosts" {
			groupName = parseName(inventoryFile[j-1])
			break
		}
		// indentLvl := indentationLevel(inventoryFile[j])
		// if indentLvl+2 == hostIndent-1 {
		// 	groupName = parseName(inventoryFile[j])
		// 	break
		// }

	}

	// get all the hosts of this group.
	for j := i + 1; j < len(inventoryFile); j++ {
		indentLvl := indentationLevel(inventoryFile[j])

		// if it's the children then is the start of a sub group.
		if strings.Contains(inventoryFile[j], "children:") {
			h, _ := getHosts(inventoryFile, j+1)

			hosts = append(hosts, h...)

			break

		}
		if indentLvl != hostIndent {
			break
		}

		// if we got heat then it's a host
		varsIndent := indentLvl + 1
		host := Host{
			Name:   parseName(inventoryFile[j]),
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

		hosts = append(hosts, host)
		// fmt.Printf("host: %+v\n", host)

	}
	// fmt.Printf("hosts %+v\n", hosts)

	return hosts, groupName
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
