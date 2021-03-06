package git

import "io"

const (
	copySection = iota
	insertSection
)

type Delta struct {
	Type    ObjType     //commit, blob, tree, tag
	Content []byte      //content
	Base    interface{} //hashid or offset(string or int)
}

//ParseCopyOrInsert parse copy or insert section info from delta reader
func ParseCopyOrInsert(r io.Reader) (stype int, offset, length int64, err error) {
	b, err := ReadOneByte(r)
	if err != nil {
		return
	}
	var _b byte
	switch IsMsbSet(b) {
	case true: //copy section
		stype = copySection
		//check last 4 byte
		for i := uint(0); i < 7; i++ {
			//we should read 1 byte from reader
			_b = 0
			if b&(1<<i) != 0 {
				_b, err = ReadOneByte(r)
				if err != nil {
					break
				}
			}
			//fmt.Printf("i:%d, _b:%b, :%b\n", i, _b, b&(1<<i))
			if i < 4 {
				offset += int64(_b) << (i * 8)
			} else {
				length += int64(_b) << ((i - 4) * 8)
			}
		}
	case false: //insert section
		stype = insertSection
		length = int64(b)
	}
	return
}

//PatchDelta recover real object data
func PatchDelta(base io.ReaderAt, delta io.Reader, target io.ReadWriter) (err error) {
	baseLen, err := ParseVarLen(delta)
	if err != nil {
		return
	}
	targetLen, err := ParseVarLen(delta)
	_ = baseLen
	_ = targetLen
	if err != nil {
		return
	}
	for {
		st, offset, length, err := ParseCopyOrInsert(delta)
		if length < 1 {
			break
		}
		if err != nil {
			break
		}
		bs := make([]byte, length)
		switch st {
		case copySection:
			_, err = base.ReadAt(bs, offset)
			if err != nil {
				break
			}
		case insertSection:
			_, err = delta.Read(bs)
			if err != nil {
				break
			}
		}
		_, err = target.Write(bs)
		if err != nil {
			break
		}
	}
	//fmt.Printf("bl:%d,tl:%d", baseLen, targetLen)
	return
}
