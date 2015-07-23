package git

import (
	"bytes"
	"testing"
)

func TestParseVarLen(t *testing.T) {
	var dataList = []struct {
		bytes []byte
		len   int64
	}{
		{[]byte{145, 46, 100}, 5905},
		{[]byte{1, 128, 100}, 1},
	}
	for _, data := range dataList {
		buf := bytes.NewBuffer(data.bytes)
		len, err := ParseVarLen(buf)
		if err != nil || len != data.len {
			t.Errorf("exptected %d, got %d", data.len, len)
		}
	}
}
