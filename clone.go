package main

import (
	"errors"
	"flag"
	"fmt"

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
	dir := ""
	if len(args) > 1 {
		dir = args[1]
	}
	repo, err := git.NewRepo(args[0], dir)
	if err != nil {
		return
	}
	defer repo.Distruct()
	if file.IsExist(repo.ClonePath) {
		err = errors.New("fatal: destination path '" + repo.ClonePath + "' already exists and is not an empty directory.")
		return
	}
	file.MakeDir(repo.ClonePath)
	//Todo:may be we should Chdir
	repoFile, err := repo.GetTmpPackFile()
	if err != nil {
		return
	}
	defer func() {
		repoFile.Close()
		//os.Remove(tmpPackFilePath)
	}()
	fmt.Printf("Cloning into '%s'...\n", repo.Name)
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
		case git.ERROR_FRAME: //had convert to error value in receiveWithSideband
			fmt.Println(string(data))
		default:
		}
	})
	if err != nil {
		return
	}
	err = repo.SaveLooseObjects(repoFile)
	if err != nil {
		return
	}
	return
}
