package git

import (
	"fmt"
	"net"
	"net/url"
	
	"github.com/bargez/pktline"
)

const (
	FLUSH_PACKET = "0000"
	PACKET_SP = " "
)

type Remote struct {
	Host string
	Port string
	Repo string
}

func (r *Remote) LsRemote() {
	conn := r.connect()
	defer conn.Close()
	pktLineEder := pktline.NewEncoderDecoder(conn)
	cmd := r.getCmd("upload-pack")
	err := pktLineEder.Encode(cmd)
	if err != nil {
		panic(err)
	}
	pktLines := make([][]byte, 0)
	err = pktLineEder.DecodeUntilFlush(&pktLines)
	pktLineEder.Encode(nil)
	if err != nil {
		panic(err)
	}
	for _, v := range pktLines{
		fmt.Println(string(v))
	}
}

func (r *Remote) Clone(){
	conn := r.connect()
	defer conn.Close()
	pktLineEder := pktline.NewEncoderDecoder(conn)
	cmd := r.getCmd("upload-pack")
	err := pktLineEder.Encode(cmd)
	if err != nil {
		panic(err)
	}
	pktLines := make([][]byte, 0)
	err = pktLineEder.DecodeUntilFlush(&pktLines)
	pktLineEder.Encode(nil)
	if err != nil {
		panic(err)
	}
	for _, v := range pktLines{
		fmt.Println(string(v))
		//first line with capabilities
	}
}

func (r *Remote ) getCmd(git_cmd string) []byte {
	//0032git-upload-pack /git-bottom-up\0Host=localHost\0
	msg := fmt.Sprintf("git-%s %s\000Host=%s\000", git_cmd, r.Repo, r.Host)
	return []byte(msg)
}

func (r *Remote) connect() (conn net.Conn) {
	conn, err := net.Dial("tcp", net.JoinHostPort(r.Host, r.Port))
	if err != nil {
		panic(err)
	}
	return
}

func NewRemote(git_url string) *Remote {
	gurl, err := url.Parse(git_url)
	if err != nil {
		panic(err)
	}
	host, port, err := net.SplitHostPort(gurl.Host)
	if err != nil {
		port = "9418"
		host = gurl.Host
	}
	return &Remote{
		host,
		port,
		gurl.Path,
	}
}
