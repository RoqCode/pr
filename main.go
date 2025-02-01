package main

import (
	"fmt"
	"os/exec"
)

func main() {
	// get branch name
	branchCmd := exec.Command("git", "rev-parse", "--abbrev-ref", "HEAD")
	branchName, err := branchCmd.Output()
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	fmt.Println(string(branchName))

	// get remote url
	remoteCmd := exec.Command("git", "config", "--get", "remote.origin.url")
	remoteName, err := remoteCmd.Output()
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	fmt.Println(string(remoteName))
}
