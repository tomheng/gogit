package git

import "testing"

func TestIsMsbSet(t *testing.T) {
	var byteList = []struct {
		b   byte
		set bool
	}{
		{'\x01', false},
		{'\x80', true},
		{'\x89', true},
	}
	for _, bd := range byteList {
		if IsMsbSet(bd.b) != bd.set {
			t.Errorf(" %d is %v", bd.b, bd.set)
		}
	}
}
