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
	offset            int64
	Version, ObjCount uint32
	ObjStore          *ObjectStore
}

//NewPackReader create a pack reader
func NewPackReader(r io.ReaderAt, size int64) (*PackReader, error) {
	sr := io.NewSectionReader(r, 0, size-20)
	version, objectCount, err := ParsePackHeader(sr)
	if err != nil {
		return nil, err
	}
	return &PackReader{
		SectionReader: sr,
		Version:       version,
		ObjCount:      objectCount,
		ObjStore:      NewObjectStore(objectCount),
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
	for i := uint32(0); i < pack.ObjCount; i++ {
		obj, offset, err := pack.ParseObjectEntry()
		//fmt.Println(obj.Type)
		/*if err == io.EOF {
			return nil
		}*/
		if err != nil {
			break
		}
		if obj == nil || f == nil {
			continue
		}
		err = pack.ObjStore.AddObject(obj, offset)
		if err != nil {
			break
		}
		err = f(obj)
		if err != nil {
			break
		}
	}
	//Todo: check SHA1
	//check if we reach the end of reader
	if pack.Tell() < pack.Size() {
		return errors.New("pack has junk at the end")
	}
	return
}

//Tell tell current cursor on reader
func (pack *PackReader) Tell() int64 {
	n, _ := pack.SectionReader.Seek(0, 1)
	return n
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
func (pack *PackReader) ParseObjectEntry() (obj *Object, offset int64, err error) {
	offset = pack.Tell()
	b, err := ReadOneByte(pack)
	if err != nil {
		return
	}
	objType := ObjType(b & 0x70 >> 4)
	var (
		objLen uint64 //unsupport big than uinit64
		//offset int64
		//shift uint = 4
		base interface{}
	)
	objLen |= uint64(b) & 0x0f
	if IsMsbSet(b) {
		objLen += readMSBEncodedSize(pack, 4)
	}
	switch objType {
	case OBJ_COMMIT:
	case OBJ_TREE:
	case OBJ_BLOB:
	case OBJ_TAG:

	case OBJ_OFS_DELTA:
		// read negative offset
		binary.Read(pack, binary.BigEndian, &b)
		noffset := int64(b & 0x7f)
		for (b & 0x80) != 0 {
			noffset += 1
			binary.Read(pack, binary.BigEndian, &b)
			noffset = (noffset << 7) + int64(b&0x7f)
		}
		base = offset - noffset
	case OBJ_REF_DELTA: //Todo:maybe we don`t support this
		tmpID := make([]byte, 20)
		n, err := pack.Read(tmpID)
		if err != nil {
			return nil, 0, err
		}
		if n != 20 {
			err = errors.New("read less than 20 bytes")
			return nil, 0, err
		}
		base = tmpID
	default:
		err = errors.New(fmt.Sprintf("unkown object type %d", objType))
		return
	}
	oc, err := InflateZlib(pack.SectionReader, int(objLen))
	if err != nil {
		return
	}
	obj, err = NewObject(objType, oc, base)
	return obj, offset, err
}
