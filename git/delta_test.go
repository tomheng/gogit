package git

import (
	"bytes"
	"io"
	"io/ioutil"
	"os"
	"testing"
)

func TestPatch(t *testing.T) {
	basef, err := os.Open("testdata/base.txt")
	if err != nil {
		t.Error(err)
	}
	bfi, err := basef.Stat()
	if err != nil {
		t.Error(err)
	}
	targetf, err := os.Open("testdata/target.txt")
	if err != nil {
		t.Error(err)
	}
	deltaf, err := os.Open("testdata/delta.txt")
	if err != nil {
		t.Error(err)
	}
	buf := bytes.NewBuffer(make([]byte, 0))
	err = PatchDelta(io.NewSectionReader(basef, 0, bfi.Size()), deltaf, buf)
	if err != nil {
		t.Error(err)
	}
	rbs, err := ioutil.ReadAll(targetf)
	if err != nil {
		t.Error(err)
	}
	if string(rbs) != string(buf.Bytes()) {
		t.Error("diff to target file")
	}
}

func TestParseCopyOrInsert(t *testing.T) {
	bs := []struct {
		data   []byte
		offset int64
		length int64
		stype  int
	}{
		{
			[]byte{
				186, //10111010
				215, //11010111
				75,  //01001011
				209, //11010001
				1,   //00000001
			},
			1258346240,
			465,
			copySection,
		},
		{
			[]byte{
				58,  //00111010
				215, //11010111
				75,  //01001011
				209, //11010001
				1,   //00000001
			},
			0,
			58,
			insertSection,
		},
	}
	for _, delta := range bs {
		buf := bytes.NewBuffer(delta.data)
		st, offset, length, err := ParseCopyOrInsert(buf)
		if err != nil {
			t.Error(err)
		}
		if st != delta.stype {
			t.Error("expect %s, go %s section", delta.stype, st)
		}
		if offset != delta.offset {
			t.Errorf("length expected %d. got %d", delta.offset, offset)
		}
		if length != delta.length {
			t.Errorf("offset should be %d, got %d", delta.length, length)
		}
	}
}
