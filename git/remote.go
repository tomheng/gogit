package git

import (
	"bytes"
	"fmt"
	"net"
	"net/url"
	"strings"
	"text/template"

	"github.com/bargez/pktline"
)

//https://www.kernel.org/pub/software/scm/git/docs/v1.7.0.5/technical/pack-protocol.txt
const (
	FLUSH_PKT = "0000"
	PKT_SP    = " "
)

type Remote struct {
	Host         string
	Port         string
	Repo         string
	Capabilities string
	AllRefs      Refs
}

func (r *Remote) Clone() {
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
	for i, line := range pktLines {
		id, name, other := ParsePktLine(line)
		if len(id) < 1 {
			continue //bad format
		}
		object := Object{
			id,
			"commit", //guess value
		}
		if i == 0 {
			r.Capabilities = other
		}
		name = strings.TrimSuffix(name, "^{}")
		ref, ok := r.AllRefs[name]
		if ok {
			ref.Child = object
		} else {
			ref = Ref{
				object,
				Object{},
			}
		}
		r.AllRefs[name] = ref
	}
	fmt.Println(r)
}

func (r *Remote) String() string {
	infoTemp := `
Host: {{.Host}}
Port: {{.Port}}
Repo: {{.Repo}}
Capabilities: {{.Capabilities}}
Refs:{{range $k, $v := .AllRefs}}
	{{$k}}:{{$v.Obj.Id|printf "%s"}},{{$v.Child.Id|printf "%s"}}
{{end}}
`
	t := template.Must(template.New("info").Parse(infoTemp))
	bfio := bytes.NewBufferString("")
	err := t.Execute(bfio, r)
	if err != nil {
		fmt.Println(err)
	}
	return bfio.String()
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
	for _, v := range pktLines {
		fmt.Println(string(v))
		//first line with capabilities
	}
}

func (r *Remote) getCmd(git_cmd string) []byte {
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
		Host:    host,
		Port:    port,
		Repo:    gurl.Path,
		AllRefs: make(Refs),
	}
}
