package git

import (
	"errors"
	"io"
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

//parse variable-length integers
func ParseVarLen(r io.Reader) (len int64, err error) {
	buf := make([]byte, 1)
	n, err := r.Read(buf)
	if err != nil {
		return
	}
	if n < 1 {
		return 0, errors.New("less than 1 byte")
	}
	var shift uint = 7
	len |= int64(buf[0]) & '\x7f'
	for IsMsbSet(buf[0]) {
		n, err = r.Read(buf)
		if err != nil {
			return
		}
		if n < 1 {
			break
		}
		len |= (int64(buf[0]) & '\x7f') << shift
		shift += 7
	}
	return
}
