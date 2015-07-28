package git

import (
	"bytes"
	"errors"
)

type ObjectStore struct {
	List        []*Object
	offsetIndex map[int64]int
	hashIndex   map[string]int
}

//AddObject add object to pack object list
func (objStore *ObjectStore) AddObject(obj *Object, offset int64) (err error) {
	if obj.Type.IsDelta() {
		//may be we should parse delta object here
		err := objStore.ParseDelta(obj)
		if err != nil {
			//break
		}
	}
	id := obj.GetID()
	index := len(objStore.List) - 1
	objStore.List[index] = obj
	objStore.hashIndex[id] = index
	objStore.offsetIndex[offset] = index
	return
}

func (objStore *ObjectStore) ParseDelta(delta *Object) (err error) {
	//find base object
	//patch it to real object
	var baseObject *Object
	switch delta.Type {
	case OBJ_OFS_DELTA:
		offset, ok := delta.Base.(int64)
		if !ok {
			err = errors.New("base should be int64 offset")
			return
		}
		baseObject = objStore.FindByOffset(offset)
	case OBJ_REF_DELTA:
		id, ok := delta.Base.(string)
		if !ok {
			err = errors.New("base should be string hash id")
			return
		}
		baseObject = objStore.FindByHash(id)
	}
	if baseObject == nil {
		err = errors.New("baseObject is nil")
		return
	}
	var brw bytes.Buffer
	err = PatchDelta(bytes.NewReader(baseObject.Content), bytes.NewReader(delta.Content), &brw)
	if err != nil {
		return
	}
	delta.Type = OBJ_BLOB
	delta.Content = brw.Bytes()
	return
}

func (objStore *ObjectStore) FindByOffset(offset int64) (obj *Object) {
	index, ok := objStore.offsetIndex[offset]
	if !ok {
		return nil
	}
	return objStore.List[index]
}

func (objStore *ObjectStore) FindByHash(hash string) (obj *Object) {
	index, ok := objStore.hashIndex[hash]
	if !ok {
		return nil
	}
	return objStore.List[index]
}

func NewObjectStore(count uint32) *ObjectStore {
	objList := make([]*Object, count)
	return &ObjectStore{
		objList,
		make(map[int64]int, count/2),
		make(map[string]int, count/2),
	}
}
