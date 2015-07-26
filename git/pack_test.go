package git

import (
	"io"
	"os"
	"testing"
)

func TestIsMsbSet(t *testing.T) {
	var byteList = []struct {
		b   byte
		set bool
	}{
		{'\x01', false},
		{'\x80', true},
		{'\x89', true},
		{'\x02', false},
	}
	for _, bd := range byteList {
		if IsMsbSet(bd.b) != bd.set {
			t.Errorf(" %d is %v", bd.b, bd.set)
		}
	}
}

func TestParsePackHeader(t *testing.T) {
	repoFile, err := os.Open("testdata/data.pack")
	if err != nil {
		return
	}
	fi, err := repoFile.Stat()
	if err != nil {
		return
	}
	packReader, err := NewPackReader(io.NewSectionReader(repoFile, 0, fi.Size()))
	if err != nil {
		t.Error(err)
	}
	if packReader.Version != 2 {
		t.Errorf("version expected 2, got ", packReader.Version)
	}
	if packReader.ObjectCount != 926 {
		t.Errorf("objectCount expected 2, got ", packReader.ObjectCount)
	}
}

func TestParseObjectEntry(t *testing.T) {
	repoFile, err := os.Open("testdata/data.pack")
	if err != nil {
		return
	}
	fi, err := repoFile.Stat()
	if err != nil {
		return
	}
	packReader, err := NewPackReader(io.NewSectionReader(repoFile, 0, fi.Size()))
	if err != nil {
		t.Error(err)
	}
	_, err = packReader.ParseObjectEntry()
	if err != nil {
		t.Error(err)
	}
}
