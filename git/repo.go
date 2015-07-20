package git

import (
	"fmt"
	"path"
	"strings"
)

type Repo struct {
	conn *GitConn
	url  *GitURL
}

func NewRepo(addr string) (repo *Repo, err error) {
	gitUrl := NewGitURL(addr)
	conn, err := NewGitConn(gitUrl.Host, gitUrl.Port)
	if err != nil {
		return
	}
	repo = &Repo{
		conn,
		gitUrl,
	}
	return
}

//distruct repo
func (repo *Repo) Distruct() {
	repo.conn.Close()
}

func (repo *Repo) GetName() string {
	_, repoName := path.Split(repo.url.RepoPath)
	return repoName
}

//https://www.kernel.org/pub/software/scm/git/docs/v1.7.0.5/technical/pack-protocol.txt
func (repo *Repo) RefDiscover() (refs Refs, capabilities []string, err error) {
	err = repo.SendCmd("upload-pack")
	if err != nil {
		//fmt.Printf("panic here:", err)
		return
	}
	pktLines, err := repo.conn.ReadPktLine()
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

//fetch pack
//sideBandHandle func(dataType byte, data []byte)
func (repo *Repo) FetchPack(sideBandHandle func(dataType byte, data []byte)) (err error) {
	refs, _, err := repo.RefDiscover()
	if err != nil {
		return
	}
	want_ids := make([]string, 0)
	//Todo:compare local commit and remote commit
	//create want list and have list
	for name, ref := range refs {
		//filter some unsupport
		if ref.IsPeeled() {
			continue
		}
		switch {
		case strings.HasPrefix(name, "HEAD"):
			fallthrough
		case strings.HasPrefix(name, "refs/tags/"):
			fallthrough
		case strings.HasPrefix(name, "refs/heads/"):
			want_ids = append(want_ids, ref.Id)
		default:
			//fmt.Println(name, " skiped")
		}
		//fmt.Println(ref.Id, "\t", name)
	}
	err = repo.SendWantList(want_ids...)
	for {
		dataType, data, done, err := repo.conn.receiveWithSideband()
		if done {
			break
		}
		if err != nil {
			return err
		}
		go sideBandHandle(dataType, data)
	}
	//wait all return
	return
}

//send cmd to server
func (repo *Repo) SendCmd(cmds ...string) (err error) {
	for _, cmd := range cmds {
		switch cmd {
		case "upload-pack":
			msg := fmt.Sprintf("git-%s %s\000Host=%s\000", cmd, repo.url.RepoPath, repo.url.Host)
			_, err = repo.conn.WritePktLine([]byte(msg))
			if err != nil {
				return
			}
		}
	}
	//_, err = repo.conn.WritePktLine(nil) //flush pktline
	return
}

func (repo *Repo) SendWantList(ids ...string) (err error) {
	var lines [][]byte
	for i, id := range ids {
		msg := "want" + PKT_SP + id
		if i == 0 {
			caps := getSupportCapabilities()
			msg += PKT_SP + strings.Join(caps, " ")
		}
		msg += PKT_LR
		lines = append(lines, []byte(msg))
		//fmt.Println("========:", string(msg))
	}
	lines = append(lines, nil, []byte("done\000")) //flush pktline
	err = repo.conn.WriteMultiPktLine(lines)       //flush pktline
	return
}
