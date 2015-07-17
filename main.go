package main

import (
	"github.com/tomheng/go-git-clone/git"
)

func main() {
	/*conn, err := net.Dial("tcp", "127.0.0.1:9418")
	if err != nil {
		fmt.Println("one")
		panic(err)
	}
	_, err = conn.Write([]byte("0032git-upload-pack /git-bottom-up\000host=localhost\000"))
	if err != nil {
		fmt.Println("2")
		panic(err)
	}
	resp := make([]byte, 1024)
	for{
		_, err = conn.Read(resp)
		if err != nil {
			fmt.Println("3")
			panic(err)
		}
		fmt.Println(string(resp))
	}*/
	url := "git://127.0.0.1/flybird"
	/*remote := git.NewRemote(url)
	remote.LsRemote()*/
	//remote.Clone()
	//git.Clone("git://127.0.0.1/flybird", "")
	//git.LsRemote("git://127.0.0.1/flybird")
	git.LsRemote(url)
}
