package main

import (
	"fmt"
	"os"

	"ascii-art/internal/args"
	"ascii-art/internal/parser"
	"ascii-art/internal/renderer"
)

func main() {
	output, text, banner, err := args.ParseArgs(os.Args[1:])
	if err != nil {
		fmt.Println(err)
		return
	}

	if text == "" {
		return
	}

	for _, ch := range text {
		if (ch < ' ' || ch > '~') && ch != '\n' && ch != '\r' {
			fmt.Println("Usage: go run . [OPTION] [STRING] [BANNER]\n\nEX: go run . --output=<fileName.txt> something standard")
			return
		}
	}

	font, err := parser.ParseBanner("ascii-style/" + banner + ".txt")
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	result, err := renderer.RenderText(text, font)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	if output != "" {
		err = os.WriteFile(output, []byte(result), 0644)
		if err != nil {
			fmt.Printf("Error: %v\n", err)
		}
		return
	}

	fmt.Print(result)
}
