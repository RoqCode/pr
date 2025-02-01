package main

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

func getBranchName() (string, error) {
	branchCmd := exec.Command("git", "rev-parse", "--abbrev-ref", "HEAD")
	branchName, err := branchCmd.Output()
	if err != nil {
		fmt.Println(err.Error())
		return "", err
	}

	return strings.TrimSpace(string(branchName)), nil
}

func getRemoteName() (string, error) {
	remoteCmd := exec.Command("git", "config", "--get", "remote.origin.url")
	remoteName, err := remoteCmd.Output()
	if err != nil {
		fmt.Println(err.Error())
		return "", err
	}

	return strings.TrimSpace(string(remoteName)), nil
}

func getRepoInfo(remoteName *string) (string, string, error) {
	// get GitHubUsername and repoName
	// todo handle bitbucket
	if strings.HasPrefix(*remoteName, "git@") {
		// todo handle https
		repoPath := strings.TrimSpace(strings.Split(*remoteName, ":")[1])
		fmt.Println(repoPath)
		path := strings.TrimSuffix(repoPath, ".git")

		// todo handle edge cases in string format
		parts := strings.Split(path, "/")

		username, repoName := parts[0], parts[1]

		return username, repoName, nil
	}

	err := errors.New("unhandled exception: remote is not SSH or GitHub")

	return "", "", err
}

func main() {
	branchName, err := getBranchName()
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	remoteName, err := getRemoteName()
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	username, repoName, err := getRepoInfo(&remoteName)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	// build url
	// https://github.com<GitHubUsername>/<repoName>/pull/new/<branchName>
	url := fmt.Sprintf("https://github.com/%s/%s/pull/new/%s", username, repoName, branchName)

	fmt.Println(url)
	exec.Command("open", url)
}
