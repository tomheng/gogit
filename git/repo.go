package git

//https://www.kernel.org/pub/software/scm/git/docs/v1.7.0.5/technical/pack-protocol.txt

import (
	"errors"
	"fmt"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/tomheng/gogit/internal/file"
)

//Repo git repo instance
type Repo struct {
	conn      *Conn
	url       *URL
	ClonePath string
	Name      string
}

//NewRepo create new Repo struct
func NewRepo(addr, dir string) (repo *Repo, err error) {
	gitURL := NewURL(addr)
	_, repoName := path.Split(gitURL.RepoPath)
	if len(dir) < 1 {
		dir = repoName
	}
	conn, err := NewConn(gitURL.Host, gitURL.Port)
	if err != nil {
		return
	}
	repo = &Repo{
		conn,
		gitURL,
		dir,
		repoName,
	}
	return
}

//GetTmpPackFile return clone temp pack file
//E.g. in native git this is something like .git/objects/pack/tmp_pack_6bo2La
func (repo *Repo) GetTmpPackFile() (*os.File, error) {
	filePath := path.Join(repo.ClonePath, ".git/objects/pack/tmp_pack_incoming")
	return file.OpenFile(filePath, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0)
}

//Distruct distruct repo
func (repo *Repo) Distruct() {
	repo.conn.Close()
}

//RefDiscover just get refs from server
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

//GetRepoFilePath get repo file path with current repo path
func (repo *Repo) GetRepoFilePath(dir, fname string) string {
	return filepath.Join(repo.ClonePath, ".git", dir, fname)
}

//CreateLocalRefs save ref to local disk
func (repo *Repo) CreateLocalRefs(name string, ref Ref) (err error) {
	var refPath string
	switch {
	case strings.HasPrefix(name, "refs/tags/"): //tags
		refPath = repo.GetRepoFilePath("refs/tags", strings.TrimLeft(name, "refs/tags/"))
	case strings.HasPrefix(name, "refs/heads/"): //origin
		refPath = repo.GetRepoFilePath("refs/remotes/origin", strings.TrimLeft(name, "refs/heads/"))
	case name == "HEAD":
		refPath = repo.GetRepoFilePath("refs/remotes/origin", "HEAD")
		//update local HEAD
		err = file.WriteFile(repo.GetRepoFilePath("", "HEAD"), []byte("ref: refs/remotes/origin/master\n"), 0644)
		if err != nil {
			return
		}
	}
	if len(refPath) > 0 {
		err = file.WriteFile(refPath, []byte(ref.ID+"\n"), 0644)
	}
	return
}

//FetchPack negotionate with remote server
//send want list and parse pack file to objects
//sideBandHandle func(dataType byte, data []byte)
func (repo *Repo) FetchPack(
	sideBandHandle func(dataType byte, data []byte),
	refHandle func(name string, ref Ref) (err error),
) (err error) {
	refs, _, err := repo.RefDiscover()
	if err != nil {
		return
	}
	var wantIDs []string //wantIDs := make([]string, 0)
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
			wantIDs = append(wantIDs, ref.ID)
		default:
			//fmt.Println(name, " skiped")
		}
		if refHandle != nil {
			err = refHandle(name, ref)
			if err != nil {
				return err
			}
		}
		//fmt.Println(ref.ID, "\t", name)
	}
	if sideBandHandle == nil {
		return
	}
	err = repo.SendWantList(wantIDs...)
	for {
		dataType, data, done, err := repo.conn.receiveWithSideband()
		if done {
			break
		}
		if err != nil {
			break
		}
		sideBandHandle(dataType, data)
	}
	//wait all return
	return
}

//SaveObject save object on disk
func (repo *Repo) SaveObject(obj *Object) (err error) {
	id := obj.GetID()
	if len(id) < 40 {
		return errors.New("wrong object id")
	}
	filePath := filepath.Join(repo.ClonePath, ".git/objects/", id[:2], id[2:])
	//fmt.Println(obj)
	fh, err := file.Create(filePath)
	if err != nil {
		return
	}
	bs, err := obj.DeflateZlib()
	if err != nil {
		return
	}
	_, err = fh.Write(bs)
	if err != nil {
		return
	}
	return
}

//SaveLooseObjects directly use pack file to recover objects
func (repo *Repo) SaveLooseObjects(f *os.File) (err error) {
	fi, err := f.Stat()
	if err != nil {
		return
	} //The final 20 bytes of the file are a SHA-1 checksum of all the previous data in the file.
	//Todo: check SHA-1 checksum
	//packReader, err := NewPackReader(io.NewSectionReader(f, 0, fi.Size()-20))
	packReader, err := NewPackReader(f, fi.Size())
	if err != nil {
		return
	}
	return packReader.ParseObjects((*repo).SaveObject)
}

//SendCmd send cmd to server
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

/*SendWantList telling the
server what objects it wants and what objects it has, so the server
can make a packfile that only contains the objects that the client needs.
The client will also send a list of the capabilities it wants to be in
effect, out of what the server said it could do with the first 'want' line.
----
  upload-request    =  want-list
		       have-list
		       compute-end

  want-list         =  first-want
		       *additional-want
		       flush-pkt

  first-want        =  PKT-LINE("want" SP obj-id SP capability-list LF)
  additional-want   =  PKT-LINE("want" SP obj-id LF)

  have-list         =  *have-line
  have-line         =  PKT-LINE("have" SP obj-id LF)
  compute-end       =  flush-pkt / PKT-LINE("done")
----
*/
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
	}
	lines = append(lines, nil, []byte("done\000")) //flush pktline
	err = repo.conn.WriteMultiPktLine(lines)       //flush pktline
	return
}
