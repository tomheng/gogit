package git

import (
	"bytes"
	"errors"
	"net"

	"github.com/bargez/pktline"
)

//Conn handle git connection(ssh or git or http)
type Conn struct {
	net.Conn
	//*pktline.Encoder
	//*pktline.Decoder
}

//NewConn create git conn
func NewConn(host, port string) (*Conn, error) {
	conn, err := net.Dial("tcp", net.JoinHostPort(host, port))
	if err != nil {
		return nil, err
	}
	return &Conn{
		conn,
	}, err
}

//WritePktLine wite pktline to conn
func (gconn *Conn) WritePktLine(line []byte) (int, error) {
	pktLineEder := pktline.NewEncoderDecoder(gconn)
	err := pktLineEder.Encode(line)
	return len(line), err
}

//WriteMultiPktLine write multi pktline co conn
func (gconn *Conn) WriteMultiPktLine(lines [][]byte) error {
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

//WriteEndPktLine write end pkt line
func (gconn *Conn) WriteEndPktLine() (int, error) {
	return gconn.WritePktLine(nil)
}

//ReadPktLine read pktline from conn
func (gconn *Conn) ReadPktLine() ([]string, error) {
	var pktlBytes [][]byte //pktlBytes := make([][]byte, 0)
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

//receiveWithSideband read data from conn
func (gconn *Conn) receiveWithSideband() (dataType byte, data []byte, done bool, err error) {
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
