package git

import (
	"bufio"
	"bytes"
	"compress/zlib"
	"encoding/binary"
	"fmt"
	"io"
	"net"
	"net/url"
)

//URL wrap url.URL
type URL struct {
	Host     string
	Port     string
	RepoPath string
}

//NewURL covert git url string to a gitURL Stuct
func NewURL(addr string) *URL {
	gurl, err := url.Parse(addr)
	if err != nil {
		panic(err)
	}
	host, port, err := net.SplitHostPort(gurl.Host)
	if err != nil {
		port = "9418"
		host = gurl.Host
	}
	return &URL{
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
		"ofs-delta",
		"agent=git/2.3.2",
	}
}

func readMSBEncodedSize(reader io.Reader, initialOffset uint) uint64 {
	var b byte
	var sz uint64
	shift := initialOffset
	sz = 0
	for {
		binary.Read(reader, binary.BigEndian, &b)
		sz += (uint64(b) &^ 0x80) << shift
		shift += 7
		if (b & 0x80) == 0 {
			break
		}
	}
	return sz
}

//ParseVarLen parse variable-length integers
func ParseVarLen(r io.Reader) (len int64, err error) {
	var shift uint //0
	for {
		b, err := ReadOneByte(r)
		if err != nil {
			break
		}
		len |= (int64(b) & '\x7f') << shift
		if !IsMsbSet(b) {
			break
		}
		shift += 7
	}
	return
}

//ReadOneByte read only one byte from the Reader
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

//reference
func inflate(reader io.Reader, sz int) ([]byte, error) {
	zr, err := zlib.NewReader(reader)
	if err != nil {
		return nil, fmt.Errorf("error opening packfile's object zlib: %v", err)
	}
	buf := make([]byte, sz)

	n, err := zr.Read(buf)
	if err != nil {
		return nil, err
	}

	if n != sz {
		return nil, fmt.Errorf("inflated size mismatch, expected %d, got %d", sz, n)
	}

	zr.Close()
	return buf, nil
}

//InflateZlib unbuffered io
func InflateZlib(r *io.SectionReader, len int) (bs []byte, err error) {
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
	if out.Len() != len {
		return nil, fmt.Errorf("inflated size mismatch, expected %d, got %d", len, out.Len())
	}
	bs = out.Bytes()
	_, err = r.Seek(0-int64(br.Buffered()), 1)
	if err != nil {
		return
	}
	return
}
