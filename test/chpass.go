package main

import (
	"fmt"
	"../../gexpect"
//	"os"
//	"io"
	"regexp"
	"time"
)

func main() {
	child, _ := gexpect.NewSubProcess("/usr/bin/passwd")
	child.Echo()
	if err := child.Start(); err != nil {
		fmt.Println(err)
	}
	r := regexp.MustCompile("\\(current\\) UNIX password:")
	idx, _ := child.Expect(5 * time.Second, r)
	if idx >= 0 {
		child.SendLine("shava123")
	}
	if idx , _ := child.Expect(5* time.Second, regexp.MustCompile("password:")); idx >=0 {
		child.SendLine("ngcmlfu")
	}
	if idx, _ := child.Expect(5 * time.Second, regexp.MustCompile("password:")); idx >=0 {
		child.SendLine("ngcmlfu")
	}
	child.InteractTimeout(1 * time.Minute)
}
