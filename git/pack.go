package git

import (
	"encoding/binary"
	"errors"
	"fmt"
	"io"
)

//https://github.com/git/git/blob/master/Documentation/technical/pack-format.txt

const (
	packSignature = "PACK"
)

//IsMsbSet check whether the most significant bit is set
func IsMsbSet(b byte) bool {
	return b>>7 == '\x01'
}

//PackReader a struct
type PackReader struct {
	*io.SectionReader
	offset               int64
	Version, ObjectCount uint32
}

//NewPackReader create a pack reader
func NewPackReader(pack *io.SectionReader) (*PackReader, error) {
	version, objectCount, err := ParsePackHeader(pack)
	if err != nil {
		return nil, err
	}
	return &PackReader{
		SectionReader: pack,
		Version:       version,
		ObjectCount:   objectCount,
	}, nil
}

/*ParsePackHeader A header appears at the beginning and consists of the following:

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
	if signature := string(buf[:4]); signature != packSignature {
		err = errors.New("pack header has wrong signature: " + signature)
		return
	}
	version = binary.BigEndian.Uint32(buf[4:8])
	if version != 2 {
		err = fmt.Errorf("version unsupport: %d ", version)
		return
	}
	objectCount = binary.BigEndian.Uint32(buf[8:])
	return
}

//ParseObjects translate all object in pack reader
func (pack *PackReader) ParseObjects(f func(object *Object) error) (err error) {
	for {
		object, err := pack.ParseObjectEntry()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}
		if object == nil {
			continue
		}
		if f == nil {
			continue
		}
		err = f(object)
		if err != nil {
			break
		}
	}
	return nil
}

//Read record the offset
func (pack *PackReader) Read(p []byte) (n int, err error) {
	n, err = pack.SectionReader.Read(p)
	pack.offset += int64(n)
	return
}

//Tell tell current cursor on reader
func (pack *PackReader) Tell() int64 {
	return pack.offset
}

/*ParseObjectEntry parse object from pack reader
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
func (pack *PackReader) ParseObjectEntry() (object *Object, err error) {
	b, err := ReadOneByte(pack)
	if err != nil {
		return
	}
	//end
	if b == 0 {
		return nil, io.EOF
	}
	objType := ObjType(b & '\x70' >> 4)
	var (
		objLen uint64 //unsupport big than uinit64
		offset int64
		shift  uint = 4
		base   []byte
	)
	objLen |= uint64(b) & '\x1f'
	for IsMsbSet(b) {
		b, err = ReadOneByte(pack)
		if err != nil {
			return
		}
		objLen |= (uint64(b) & '\x7f') << shift
		shift += 7
	}
	switch objType {
	case OBJ_COMMIT:
	case OBJ_TREE:
	case OBJ_BLOB:
	case OBJ_TAG:

	case OBJ_OFS_DELTA:
		offset, err = ParseVarLen(pack)
		if err != nil {
			break
		}
		tmpSectionReader := io.NewSectionReader(pack, pack.Tell()-offset, offset)
		base, err = InflateZlib(tmpSectionReader)
		if err != nil {
			break
		}
	case OBJ_REF_DELTA: //Todo:maybe we don`t support this
		tmpID := make([]byte, 20)
		n, err := pack.Read(tmpID)
		if err != nil {
			break
		}
		if n != 20 {
			err = errors.New("read less than 20 bytes")
			break
		}
		base = tmpID
	default:
		//err = errors.New(fmt.Sprintf("unkown object type %d", objType))
		return
	}
	oc, err := InflateZlib(pack.SectionReader)
	if err != nil {
		return
	}
	object = NewObject(objType, oc, base)
	return
}
