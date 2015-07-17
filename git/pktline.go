package git

/*
import (
	"errors"
	"fmt"
	"strconv"
)

// Errors returned by methods in package pktline.
var (
	//ErrShortRead   = errors.New("input is too short")
	//ErrInputExcess = errors.New("input is too long")
	ErrTooLong    = errors.New("too long payload")
	ErrInvalidLen = errors.New("invalid length")
)

const (
	headLen = 4
	maxLen  = 65524 // 65520 bytes of data
)

// Encode returns payload encoded in pkt-line format.
func Encode(payload []byte) ([]byte, error) {
	if payload == nil {
		return []byte("0000"), nil
	}
	if len(payload)+headLen > maxLen {
		return nil, ErrTooLong
	}
	head := []byte(fmt.Sprintf("%04x", len(payload)+headLen))
	return append(head, payload...), nil
}

//it is simple just with some check
func Decode(line []byte) (payload []byte, err error) {
	head := line[:4]
	lineLen, err := strconv.ParseInt(string(head), 16, 16)
	if err != nil {
		return err
	}
	if lineLen == 0 { // flush-pkt
		return nil
	}
	if lineLen < headLen {
		return ErrInvalidLen
	}
	*payload = make([]byte, lineLen-headLen)
	payload := line[4 : lineLen+1]
	return
}
*/
