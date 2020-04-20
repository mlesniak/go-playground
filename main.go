package main

import "fmt"

func main() {
	fmt.Println(Message())
}

// Message returns a greeting string.
func Message() string {
	return "foo"
}