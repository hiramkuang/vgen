// examples/main.go
package main

import (
	"fmt"
	"os"

	"github.com/hiramkuang/vgen/internal/generator"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: vgen <file_path>")
		os.Exit(1)
	}

	filePath := os.Args[1]

	if err := generator.GenerateValidator(filePath); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Successfully generated validator for %s\n", filePath)
}
