package git

import "errors"

type ObjectStore struct {
	List        []*Object
	offsetIndex map[int64]int
	hashIndex   map[string]int
	refDeltas   map[string]int
}

//AddObject add object to pack object list
func (objStore *ObjectStore) AddObject(obj *Object, offset int64) (err error) {
	index := objStore.addObject(obj, offset)
	switch obj.Type {
	case OBJ_OFS_DELTA:
		offset, ok := obj.Base.(int64)
		if !ok {
			return errors.New("base should be int64 offset")
		}
		base := objStore.FindByOffset(offset)
		if base == nil {
			return errors.New("ofs_delta base object can not be nil")
		}
		err = obj.Patch(base)
		if err != nil {
			return
		}
	case OBJ_REF_DELTA:
		id, ok := obj.Base.(string)
		if !ok {
			return errors.New("base should be string hash id")
		}
		base := objStore.FindByHash(id)
		if base == nil {
			objStore.refDeltas[id] = index
		} else {
			err = obj.Patch(base)
			if err != nil {
				return
			}
		}
	}
	return objStore.checkDepDelta(obj)
}

//addObject add a object internal
func (objStore *ObjectStore) addObject(obj *Object, offset int64) (index int) {
	id := obj.GetID()
	objStore.List = append(objStore.List, obj)
	index = len(objStore.List) - 1
	objStore.hashIndex[id] = index
	objStore.offsetIndex[offset] = index
	return
}

func (objStore *ObjectStore) updateObject(index int, obj *Object) {
	id := obj.FlushID()
	objStore.List[index] = obj
	objStore.hashIndex[id] = index
}

//checkDepDelta check delta object base on
func (objStore *ObjectStore) checkDepDelta(obj *Object) (err error) {
	if obj.Type.IsDelta() {
		return
	}
	index, ok := objStore.refDeltas[obj.GetID()]
	if !ok {
		return
	}
	delta := objStore.List[index]
	err = delta.Patch(obj)
	if err != nil {
		return
	}
	if delta.Type.IsDelta() {
		return
	}
	objStore.updateObject(index, obj)
	return objStore.checkDepDelta(delta)
}

//FindByOffset find object by offset
func (objStore *ObjectStore) FindByOffset(offset int64) (obj *Object) {
	index, ok := objStore.offsetIndex[offset]
	if !ok {
		return nil
	}
	return objStore.List[index]
}

//FindByHash find object by hash id
func (objStore *ObjectStore) FindByHash(hash string) (obj *Object) {
	index, ok := objStore.hashIndex[hash]
	if !ok {
		return nil
	}
	return objStore.List[index]
}

//NewObjectStore create it
func NewObjectStore(count uint32) *ObjectStore {
	//objList := make([]*Object, count)
	return &ObjectStore{
		make([]*Object, 0),
		make(map[int64]int),
		make(map[string]int),
		make(map[string]int),
	}
}
