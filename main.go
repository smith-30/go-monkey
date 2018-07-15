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

// let map = fn(arr, f) {
// 	let iter = fn(arr, accumlated) {
// 		if (len(arr) == 0) {
// 			accumlated
// 		} else {
// 			iter(rest(arr), push(accumlated, f(first(arr))))
// 		}
// 	};

// 	iter(arr, []);
// }

// let reduce = fn(arr, initial, f) {
// 	let iter = fn(arr, result) {
// 		if (len(arr) == 0) {
// 			result
// 		} else {
// 			iter(rest(arr), f(result, first(arr)))
// 		}
// 	};

// 	iter(arr, initial);
// }

// let sum = fn(arr) {
// 	reduce(arr, 0, fn(initial, el) {initial + el});
// };
