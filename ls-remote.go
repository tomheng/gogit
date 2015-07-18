package main

import (
	"fmt"

	"github.com/tomheng/gogit/git"
)

func newLsRemoteCmd() *Command {
	return &Command{
		Run:       runLsRemote,
		UsageLine: "ls-remote list remote refs",
	}
}

func runLsRemote(cmd *Command, args []string) {
	for _, url := range args {
		refs, _, err := git.RefDiscover(url)
		if err != nil {
			panic(err)
		}
		for name, object := range refs {
			fmt.Println(object.Id, "\t", name)
		}
	}
}
