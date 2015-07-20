package main

import (
	"errors"
	"fmt"
	"os"

	"github.com/tomheng/gogit/git"
)

func newCloneCmd() *Command {
	return &Command{
		Run:       runClone,
		UsageLine: "clone specific git repo",
	}
}

//git clone on addr
func runClone(cmd *Command, args []string) (err error) {
	if len(args) < 1 {
		return errors.New("args not enogh")
	}
	repo, err := git.NewRepo(args[0])
	if err != nil {
		return
	}
	defer repo.Distruct()
	repoFile, err := os.Create("/tmp/abcd")
	if err != nil {
		return
	}
	defer repoFile.Close()
	err = repo.FetchPack(func(dataType byte, data []byte) {
		/*
			1 the remainder of the packet line is a chunk of the pack file - this is the payload channel
			2 this is progress information that the server sends - the client prints that on STDERR prefixed with remote: "
			3 this is error infomration that will cause the client to print out the message on STDERR and exit with an error code (not implemented in our example)
		*/
		switch dataType {
		case git.DATA_FRAME:
			_, err := repoFile.Write(data)
			if err != nil {
				panic(err)
				//returnData <- err
			}
		case git.PROGRESS_FRAME:
			progress := string(data)
			fmt.Printf("%c[2K\r", 27)
			fmt.Println(progress)
		case '3': //had convert to error value in receiveWithSideband
		}
	})
	if err != nil {
		return
	}
	return
}