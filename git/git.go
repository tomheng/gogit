package git

import (
	"bufio"
	"bytes"
	"compress/zlib"
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
		//"thin-pack",
		//"ofs-delta",
		"agent=git/2.3.2",
	}
}

//parse variable-length integers
func ParseVarLen(r io.Reader) (len int64, err error) {
	b, err := ReadOneByte(r)
	if err != nil {
		return
	}
	var shift uint = 7
	len |= int64(b) & '\x7f'
	for IsMsbSet(b) {
		b, err = ReadOneByte(r)
		if err != nil {
			return
		}
		len |= (int64(b) & '\x7f') << shift
		shift += 7
	}
	return
}

//read only one byte from the Reader
func ReadOneByte(r io.Reader) (b byte, err error) {
	buf := make([]byte, 1)
	n, err := r.Read(buf)
	if err != nil {
		return
	}
	if n == 1 {
		b = buf[0]
	}
	return
}

func InflateZlib(r *io.SectionReader) (bs []byte, err error) {
	var out bytes.Buffer
	br := bufio.NewReader(r)
	zr, err := zlib.NewReader(br)
	if err != nil {
		return
	}
	defer zr.Close()
	_, err = io.Copy(&out, zr)
	if err != nil {
		return
	}
	bs = out.Bytes()
	_, err = r.Seek(0-int64(br.Buffered()), 1)
	if err != nil {
		return
	}
	return
}
