package main

import "fmt"

func main() {
	name := "Aaron"
	age := 28

	greeting := fmt.Sprintf("Hello, %s! You are %d years old", name, age)
	fmt.Println(greeting)
}
