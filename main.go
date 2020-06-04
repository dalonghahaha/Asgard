package main

import (
	"fmt"
	"runtime/debug"

	"Asgard/cmds"
)

func main() {
	defer func() {
		if err := recover(); err != nil {
			fmt.Println("panic:", err)
			fmt.Println("stack:", string(debug.Stack()))
			return
		}
	}()
	cmds.Execute()
}
