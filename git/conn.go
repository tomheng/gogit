package git

import (
	"bytes"
	"net"

	"github.com/bargez/pktline"
)

type GitConn struct {
	conn net.Conn
	//*pktline.Encoder
	//*pktline.Decoder
}

//create git conn
func NewGitConn(host, port string) (*GitConn, error) {
	conn, err := net.Dial("tcp", net.JoinHostPort(host, port))
	if err != nil {
		return nil, err
	}
	return &GitConn{
		conn: conn,
	}, err
}

//wite pktline to conn
func (gconn *GitConn) WritePktLine(line []byte) (int, error) {
	pktLineEder := pktline.NewEncoderDecoder(gconn.conn)
	err := pktLineEder.Encode(line)
	return len(line), err
}

//read pktline from conn
func (gconn *GitConn) ReadPktLine() ([]string, error) {
	pktlBytes := make([][]byte, 0)
	pktLineEder := pktline.NewEncoderDecoder(gconn.conn)
	err := pktLineEder.DecodeUntilFlush(&pktlBytes)
	if err != nil {
		return nil, err
	}
	pktLines := make([]string, len(pktlBytes))
	for i, v := range pktlBytes {
		pktLines[i] = string(bytes.TrimRight(v, "\r\n"))
	}
	return pktLines, err
}

func ParsePktLine(pktLine []byte) (id, name, other string) {
	id = string(pktLine[:41])
	i := bytes.Index(pktLine, []byte{'\000'})
	if i > 0 {
		other = string(pktLine[i:])
		//continue //bad format
	} else {
		i = len(pktLine) - 1
	}
	name = string(pktLine[41:i])
	return
}
