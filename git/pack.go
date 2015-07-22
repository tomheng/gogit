package git

import (
	"bytes"
	"compress/zlib"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
)

//https://github.com/git/git/blob/master/Documentation/technical/pack-format.txt

const (
	PACK_SIGNATURE = "PACK"
	OBJ_COMMIT     = iota
	OBJ_TREE
	OBJ_BLOB
	OBJ_TAG
	_
	OBJ_OFS_DELTA
	OBJ_REF_DELTA
)

//check whether the most significant bit is set
func IsMsbSet(b byte) bool {
	return b>>7 == '\x01'
}

type PackReader struct {
	reader               *io.SectionReader
	Version, ObjectCount uint32
}

func NewPackReader(pack *io.SectionReader) (*PackReader, error) {
	version, objectCount, err := ParsePackHeader(pack)
	if err != nil {
		return nil, err
	}
	return &PackReader{
		reader:      pack,
		Version:     version,
		ObjectCount: objectCount,
	}, nil
}

/*
A header appears at the beginning and consists of the following:

     4-byte signature:
         The signature is: {'P', 'A', 'C', 'K'}

     4-byte version number (network byte order):
	 Git currently accepts version number 2 or 3 but
         generates version 2 only.

     4-byte number of objects contained in the pack (network byte order)

     Observation: we cannot have more than 4G versions ;-) and
     more than 4G objects in a pack.

*/
func ParsePackHeader(pack *io.SectionReader) (version, objectCount uint32, err error) {
	buf := make([]byte, 12)
	_, err = pack.Read(buf)
	if err != nil {
		return
	}
	if signature := string(buf[:4]); signature != PACK_SIGNATURE {
		err = errors.New("pack header has wrong signature: " + signature)
		return
	}
	version = binary.BigEndian.Uint32(buf[4:8])
	if version != 2 {
		err = errors.New(fmt.Sprintf("version unsupport: %d ", version))
		return
	}
	objectCount = binary.BigEndian.Uint32(buf[8:])
	return
}

func (pack *PackReader) ParseObjects() {

}

/*
	(undeltified representation)
     n-byte type and length (3-bit type, (n-1)*7+4-bit length)
     compressed data

     (deltified representation)
     n-byte type and length (3-bit type, (n-1)*7+4-bit length)
     20-byte base object name if OBJ_REF_DELTA or a negative relative
	 offset from the delta object's position in the pack if this
	 is an OBJ_OFS_DELTA object
     compressed delta data
*/
func (pack *PackReader) ParseObjectEntry() (err error) {
	buf := make([]byte, 1)
	n, err := pack.reader.Read(buf)
	if err != nil {
		return err
	}
	if n < 1 {
		return errors.New("less than 1 byte")
	}
	b := uint64(buf[0])
	objType := int(buf[0] & '\x70' >> 4)
	var objLen uint64 = 0 //unsupport big than uinit64
	objLen |= b & '\x1f'
	var shift uint = 4
	fmt.Printf("fb:%b\n", buf[0])
	for IsMsbSet(buf[0]) {
		_, err = pack.reader.Read(buf)
		if err != nil {
			return err
		}
		if n < 1 {
			break
		}
		objLen |= (uint64(buf[0]) & '\x7f') << shift
		shift += 7
	}
	switch objType {
	case OBJ_COMMIT:
	case OBJ_TREE:
	case OBJ_BLOB:
	case OBJ_TAG:

	case OBJ_OFS_DELTA:
	case OBJ_REF_DELTA:
	default:
		return errors.New("unkown object type")
	}
	var out bytes.Buffer
	zr, err := zlib.NewReader(pack.reader)
	if err != nil {
		return
	}
	nn, err := io.Copy(&out, zr)
	fmt.Println("expectedLen:", objLen, ",readLne:", nn)
	fmt.Println(string(out.Bytes()))
	zr.Close()
	return
}
