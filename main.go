package main

import (
	"fmt"
	"runtime/debug"

	"Asgarde/cmd"
)

func main() {
	defer func() {
		if err := recover(); err != nil {
			fmt.Println("panic:", err)
			fmt.Println("stack:", string(debug.Stack()))
			return
		}
	}()
	cmd.Execute()
}
