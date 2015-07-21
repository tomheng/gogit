package main

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"path"

	"github.com/tomheng/gogit/git"
	"github.com/tomheng/gogit/internal/file"
)

var cloneFlag = flag.NewFlagSet("clone", flag.ExitOnError)

func newCloneCmd() *Command {
	return &Command{
		Run:       runClone,
		UsageLine: "clone specific git repo",
		Flag:      *cloneFlag,
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
	dir := repo.GetName()
	if len(args) > 1 {
		dir = args[1]
	}
	if file.IsExist(dir) {
		return errors.New("fatal: destination path '" + dir + "' already exists and is not an empty directory.")
	}
	file.MakeDir(dir)
	//Todo:may be we should Chdir
	//E.g. in native git this is something like .git/objects/pack/tmp_pack_6bo2La
	tmpPackFilePath := path.Join(dir, ".git/objects/pack/tmp_pack_incoming")
	repoFile, err := file.OpenFile(tmpPackFilePath, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0)
	if err != nil {
		return
	}
	defer repoFile.Close()
	fmt.Printf("Cloning into '%s'...\n", dir)
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
			fmt.Print("\r", "remote: "+progress)
		case '3': //had convert to error value in receiveWithSideband
		}
	})
	if err != nil {
		return
	}
	return
}
