package git

//https://github.com/git/git/blob/master/Documentation/technical/pack-format.txt

func IsMsbSet(b byte) bool {
	return b >> 7 == '\x01'
}
