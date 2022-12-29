package main

import (
	"fmt"
	"os"

	ansibleinventoryparser "github.com/ezratameno/ansible-inventory-parser/pkg/ansible-inventory-parser"
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

}

func run() error {
	return ansibleinventoryparser.Parse("inv.yaml")
}
