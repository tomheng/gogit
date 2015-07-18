package main

import "github.com/tomheng/gogit/git"

func newLsRemoteCmd() *Command {
	return &Command{
		Run:       runLsRemote,
		UsageLine: "ls-remote list remote refs",
	}
}

func runLsRemote(cmd *Command, args []string) {
	for _, url := range args {
		git.LsRemote(url)
	}
}
