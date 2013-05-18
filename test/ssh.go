package main

import (
	"fmt"
	"../../gexpect"
	"time"
)

func main() {
	child, _ := gexpect.NewSubProcess("/usr/bin/ssh", "-p 2222", "knightmare@rackol.com")
	if err := child.Start(); err != nil {
		fmt.Println(err)
	}
	defer child.Close()
	child.InteractTimeout(5 * time.Second)
}
