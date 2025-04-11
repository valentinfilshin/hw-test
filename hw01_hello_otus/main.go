package main

import (
	"fmt"

	"golang.org/x/example/hello/reverse"
)

func main() {
	s := "Hello, OTUS!"
	revertedString := reverse.String(s)
	fmt.Println(revertedString)
}
