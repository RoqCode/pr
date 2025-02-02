package main

import (
	"errors"
	"fmt"
	"net/url"
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

func getRepoInfo(remoteName string) (string, string, string, string, error) {
	if strings.HasPrefix(remoteName, "git@") || strings.HasPrefix(remoteName, "ssh://git@") {
		repoUrl := strings.TrimSpace(strings.Split(remoteName, "@")[1])
		repoUrl = strings.Split(repoUrl, ":")[0]
		// todo handle https

		repoPath := strings.TrimSpace(strings.Split(remoteName, ":")[1])
		fmt.Println(repoPath)
		path := strings.TrimSuffix(repoPath, ".git")

		// todo handle edge cases in string format
		parts := strings.Split(path, "/")

		userName := parts[0]
		var projectName string
		var repoName string

		if len(parts) == 2 {
			repoName = parts[1]
		} else if len(parts) == 3 {
			projectName = parts[1]
			repoName = parts[2]
		} else if len(parts) < 2 || len(parts) > 3 {
			err := errors.New("invalid remote URL")
			return "", "", "", "", err
		}

		return repoUrl, userName, projectName, repoName, nil
	}

	err := errors.New("unhandled exception: remote is not SSH")
	return "", "", "", "", err
}

func buildUrl(repoUrl string, userName string, projectName string, repoName string, branchName string) (string, error) {
	var prUrl string
	var err error = nil

	if strings.Contains(repoUrl, "github") {
		// is github url
		if projectName == "" {
			prUrl = fmt.Sprintf("https://%s/%s/%s/pull/new/%s", repoUrl, userName, repoName, branchName)
		} else {
			prUrl = fmt.Sprintf("https://%s/%s/%s/%s/pull/new/%s", repoUrl, userName, projectName, repoName, branchName)
		}
	} else if strings.Contains(repoUrl, "bitbucket") {
		// is bitbucket url
		if projectName == "" {
			err = errors.New("no projectName found")
		} else {
			urlSafeBranchName := url.QueryEscape("refs/heads/" + branchName)
			prUrl = fmt.Sprintf("https://%s/projects/%s/repos/%s/pull-requests?create&sourceBranch=%s", repoUrl, projectName, repoName, urlSafeBranchName)
		}
	} else {
		err = errors.New("unrecognized URL format")
	}

	return prUrl, err
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

	repoUrl, userName, projectName, repoName, err := getRepoInfo(remoteName)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	url, err := buildUrl(repoUrl, userName, projectName, repoName, branchName)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	fmt.Println("PR URL:", url)
	fmt.Print("Soll die URL geöffnet werden? (y/n): ")

	// Read the user's input
	var response string
	fmt.Scanln(&response)

	if strings.ToLower(response) == "y" {
		fmt.Println("Öffne URL...")
		// Note: 'open' works on macOS. Use 'xdg-open' for Linux or 'start' for Windows.
		cmd := exec.Command("open", url)
		err = cmd.Run()
		if err != nil {
			fmt.Println("Failed to open URL:", err)
			os.Exit(1)
		}
	} else {
		fmt.Println("URL nicht geöffnet.")
	}
	// Todo build into binary and add to path?
}
