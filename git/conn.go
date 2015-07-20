package git

import (
	"bytes"
	"errors"
	"net"

	"github.com/bargez/pktline"
)

type GitConn struct {
	net.Conn
	//*pktline.Encoder
	//*pktline.Decoder
}

//global git conn
var giConn GitConn

//create git conn
func NewGitConn(host, port string) (*GitConn, error) {
	conn, err := net.Dial("tcp", net.JoinHostPort(host, port))
	if err != nil {
		return nil, err
	}
	return &GitConn{
		conn,
	}, err
}

//wite pktline to conn
func (gconn *GitConn) WritePktLine(line []byte) (int, error) {
	pktLineEder := pktline.NewEncoderDecoder(gconn)
	err := pktLineEder.Encode(line)
	return len(line), err
}

//write multi pktline co conn
func (gconn *GitConn) WriteMultiPktLine(lines [][]byte) error {
	var data []byte
	for _, line := range lines {
		line, err := pktline.Encode(line)
		if err != nil {
			return err
		}
		data = append(data, line...)
		//flush pktline
		if string(line) == FLUSH_PKT && len(data) > 0 {
			_, err := gconn.Write(data)
			if err != nil {
				return err
			}
			data = nil
		}
	}
	_, err := gconn.Write(data)
	if err != nil {
		return err
	}
	return nil
}

func (gconn *GitConn) WriteEndPktLine() (int, error) {
	return gconn.WritePktLine(nil)
}

//read pktline from conn
func (gconn *GitConn) ReadPktLine() ([]string, error) {
	pktlBytes := make([][]byte, 0)
	pktLineEder := pktline.NewEncoderDecoder(gconn)
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

//read data from conn
func (gconn *GitConn) receiveWithSideband() (dataType byte, data []byte, done bool, err error) {
	pktLineEder := pktline.NewEncoderDecoder(gconn)
	var line []byte
	err = pktLineEder.Decode(&line)
	if err != nil {
		return
	}
	if line == nil {
		done = true
		return
	}
	dataType = line[0] //first byte
	data = line[1:]    //remain bytes
	//error
	if dataType == ERROR_FRAME {
		err = errors.New(string(data))
		return
	}
	return
}
