package main

import (
	"fmt"
	"os"
	"os/user"

	"github.com/smith-30/go-monkey/repl"
)

func main() {
	user, err := user.Current()
	if err != nil {
		panic(err)
	}

	fmt.Printf("Hello \x1b[32m%s\x1b[0m! This is the Monkey programming language!\n", user.Username)
	fmt.Printf("Feel free to type in commands\n")
	repl.Start(os.Stdin, os.Stdout)
}
