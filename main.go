package main

import "fmt"

func main() {
	for i := 1; i <= 5; i++ {
		fmt.Println("")
		for j := 1; j <= i; j++ {
			fmt.Print(" ")
		}
		for k := 5; k >= i; k-- {
			fmt.Print("*")
		}
	}
}
