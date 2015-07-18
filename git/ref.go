package git

import (
	"fmt"
	"strings"
)

type Refs map[string]Ref

func NewRefs() Refs {
	return make(Refs, 0)
}

func (ref Refs) AddRef(name, oid string) {
	ref[name] = Ref{
		Object{Id: oid},
		Object{},
	}
}

const (
	FLUSH_PKT = "0000"
	PKT_SP    = " "
)

//ref is a human name for a commit
type Ref struct {
	Object        //i think it is a commit object
	Child  Object //contain some info for this ref
}

type Object struct {
	Id       string //SHA-1 40 char
	SelfType string //commit, blob, tree
}

//https://www.kernel.org/pub/software/scm/git/docs/v1.7.0.5/technical/pack-protocol.txt
func RefDiscover(addr string) (refs Refs, capabilities []string, err error) {
	gitUrl := NewGitURL(addr)
	conn, err := NewGitConn(gitUrl.Host, gitUrl.Port)
	if err != nil {
		return
	}
	cmd := gitUrl.GetCmd("upload-pack")
	_, err = conn.WritePktLine(cmd)
	if err != nil {
		//fmt.Printf("panic here:", err)
		return
	}
	conn.WritePktLine(nil)
	pktLines, err := conn.ReadPktLine()
	if err != nil {
		return
	}
	refs = make(Refs, len(pktLines))
	//fmt.Printf("return lines:%d", len(pktLines))
	for i, line := range pktLines {
		//first line with Capabilities
		if i == 0 {
			index := strings.Index(line, "\000")
			if index > -1 {
				capabilities = strings.Split(line[index+1:], PKT_SP)
				line = line[:index]
			}
		}
		line := strings.SplitN(line, PKT_SP, 2)
		if len(line) < 2 {
			fmt.Println("coninut one")
			continue
		}
		refs.AddRef(line[1], line[0])
	}
	return
}
