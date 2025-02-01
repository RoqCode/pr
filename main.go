package main

import (
	"fmt"
	"os"
	"os/exec"
)

func getBranchName() (string, error) {
	branchCmd := exec.Command("git", "rev-parse", "--abbrev-ref", "HEAD")
	branchName, err := branchCmd.Output()
	if err != nil {
		fmt.Println(err.Error())
		return "", err
	}

	return string(branchName), nil
}

func getRemoteName() (string, error) {
	remoteCmd := exec.Command("git", "config", "--get", "remote.origin.url")
	remoteName, err := remoteCmd.Output()
	if err != nil {
		fmt.Println(err.Error())
		return "", err
	}

	return string(remoteName), nil
}

func main() {
	branchName, err := getBranchName()
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	fmt.Println(string(branchName))

	remoteName, err := getRemoteName()
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	fmt.Println(string(remoteName))
}
