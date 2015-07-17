package git

import (
	"fmt"
	"net"
	"net/url"
)

//git clone on addr
func Clone(addr, repoPath string) {
	/*conn, err := NewGitConn(addr)
	if err != nil {
		panic(err)
	}
	_ := conn*/
}

//git-ls-remote - List references in a remote repository
func LsRemote(addr string) {
	gitUrl := NewGitURL(addr)
	conn, err := NewGitConn(gitUrl.Host, gitUrl.Port)
	if err != nil {
		panic(err)
	}
	cmd := gitUrl.GetCmd("upload-pack")
	_, err = conn.WritePktLine(cmd)
	if err != nil {
		//fmt.Printf("panic here:", err)
		panic(err)
	}
	conn.WritePktLine(nil)
	pktLines, err := conn.ReadPktLine()
	if err != nil {
		panic(err)
	}
	//fmt.Printf("return lines:%d", len(pktLines))
	for _, line := range pktLines {
		fmt.Print(line)
	}
}

type GitURL struct {
	Host     string
	Port     string
	RepoPath string
}

//covert git url string to a gitURL Stuct
func NewGitURL(addr string) *GitURL {
	gurl, err := url.Parse(addr)
	if err != nil {
		panic(err)
	}
	host, port, err := net.SplitHostPort(gurl.Host)
	if err != nil {
		port = "9418"
		host = gurl.Host
	}
	return &GitURL{
		Host:     host,
		Port:     port,
		RepoPath: gurl.Path,
	}
}

func (gu *GitURL) GetCmd(git_cmd string) []byte {
	//0032git-upload-pack /git-bottom-up\0Host=localHost\0
	msg := fmt.Sprintf("git-%s %s\000Host=%s\000", git_cmd, gu.RepoPath, gu.Host)
	return []byte(msg)
}
