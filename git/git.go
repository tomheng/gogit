package git

import (
	"net"
	"net/url"
)

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

func getSupportCapabilities() []string {
	return []string{
		"multi_ack_detailed",
		"side-band-64k",
		"thin-pack",
		"ofs-delta",
		"agent=git/1.8.2",
	}
}
