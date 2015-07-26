package main

import (
	"errors"
	"fmt"

	"github.com/tomheng/gogit/git"
)

func newLsRemoteCmd() *Command {
	return &Command{
		Run:       runLsRemote,
		UsageLine: "ls-remote list remote refs",
	}
}

func runLsRemote(cmd *Command, args []string) error {
	if len(args) < 1 {
		return errors.New("args not enogh")
	}
	repo, err := git.NewRepo(args[0], "")
	if err != nil {
		return err
	}
	defer repo.Distruct()
	refs, _, err := repo.RefDiscover()
	if err != nil {
		return err
	}
	for name, ref := range refs {
		fmt.Println(ref.ID, "\t", name)
	}
	return nil
}
