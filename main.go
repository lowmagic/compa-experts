// Command compa-experts: Expert and catalog packages for compa, each a Go package exporting an fs.FS, and compose, which merges the chosen ones into the bundle a compa build ships.
package main

import (
	"fmt"
	"os"
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func run() error {
	fmt.Println("Hello from compa-experts!")
	return nil
}
