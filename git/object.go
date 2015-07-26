package git

import (
	"bytes"
	"compress/zlib"
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"io"
)

/*
OBJ_COMMIT = 1
OBJ_TREE = 2
OBJ_BLOB = 3
OBJ_TAG = 4
OBJ_OFS_DELTA = 6
OBJ_REF_DELTA = 7
*/
const (
	OBJ_COMMIT ObjType = 1 + iota //start 1.
	OBJ_TREE
	OBJ_BLOB
	OBJ_TAG
	_
	OBJ_OFS_DELTA
	OBJ_REF_DELTA
)

//ObjType represent a git object type
type ObjType int8

//String make it a stringer
func (ot ObjType) String() (ts string) {
	switch ot {
	case OBJ_COMMIT:
		ts = "commit"
	case OBJ_TREE:
		ts = "tree"
	case OBJ_BLOB:
		ts = "blob"
	case OBJ_TAG:
		ts = "tag"
	case OBJ_OFS_DELTA:
		ts = "ofs_delta"
	case OBJ_REF_DELTA:
		ts = "ref_delta"
	default:
		ts = "blob"
	}
	return
}

//Object git internal man
type Object struct {
	io.Reader
	ID      string  //SHA-1 40 char
	Type    ObjType //commit, blob, tree, tag
	Content []byte  //content
	//Base    []byte  //delta oject based object
}

//NewObject create object
//if it is a delta object, we will delating it to real content
func NewObject(objType ObjType, content, base []byte) *Object {
	if objType == OBJ_OFS_DELTA || objType == OBJ_REF_DELTA {
		brw := bytes.NewBuffer(make([]byte, 1024))
		base := bytes.NewReader(base)
		PatchDelta(io.NewSectionReader(base, 0, int64(base.Len())), bytes.NewReader(content), brw)
		objType = OBJ_BLOB
		content = brw.Bytes()
	}
	return &Object{
		//Reader: bytes.NewReader(c),
		Type: objType,
		//Base:    b,
		Content: content,
	}
}

//GetStoreHeader create a object store header
func (obj *Object) GetStoreHeader() []byte {
	return []byte(fmt.Sprintf("%s %d\000", obj.Type, len(obj.Content)))
}

//GetID generate the object id by content
func (obj *Object) GetID() (ids string) {
	if len(obj.ID) > 0 {
		return obj.ID
	}
	sw := sha1.New()
	sw.Write(obj.GetStoreHeader())
	sw.Write(obj.Content)
	obj.ID = hex.EncodeToString(sw.Sum(nil))
	return obj.ID
}

//String make object to be a stringer
func (obj *Object) String() string {
	return fmt.Sprintf("%s %s", obj.Type, obj.ID)
}

//DeflateZlib deflat the object to bytes
func (obj *Object) DeflateZlib() (bs []byte, err error) {
	var br bytes.Buffer
	zw := zlib.NewWriter(&br)
	zw.Write(obj.GetStoreHeader())
	zw.Write(obj.Content)
	zw.Close()
	bs = br.Bytes()
	return
}
