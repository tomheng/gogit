package delta

import (
	"bytes"
	"testing"
)

func TestPatch(t *testing.T) {
	/*basef, err := os.Open("testdata/base.delta")
	if err != nil {
		t.Error(err)
	}
	/*targetf, err := os.Open("testdata/target.delta")
	if err != nil {
		t.Error(err)
	}*/
	/*deltaf, err := os.Open("testdata/delta.delta")
	if err != nil {
		t.Error(err)
	}
	w, err := Patch(io.NewSectionReader(basef), io.NewSectionReader(deltaf))
	_ = w
	if err != nil {
		t.Fatal(err)
	}*/
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
			CopySection,
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
			InsertSection,
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
