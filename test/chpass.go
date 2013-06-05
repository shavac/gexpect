package main

import (
	"../../gexpect"
	"fmt"
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
	idx, _ := child.ExpectTimeout(5*time.Second, r)
	if idx >= 0 {
		child.SendLine("P@ssw0rd")
	}
	if idx, _ := child.ExpectTimeout(5*time.Second, regexp.MustCompile("password:")); idx >= 0 {
		child.SendLine("P@ssw0rd")
	}
	if idx, _ := child.ExpectTimeout(5*time.Second, regexp.MustCompile("password:")); idx >= 0 {
		child.SendLine("P@ssw0rd")
	}
	child.InteractTimeout(1 * time.Minute)
}
