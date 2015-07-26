package git

import "strings"

type Refs map[string]Ref

const (
	PEELED_FLAG = "^{}"
)

func NewRefs() Refs {
	return make(Refs, 0)
}

func (ref Refs) AddRef(name, oid string) {
	ref[name] = Ref{
		name,
		Object{Id: oid},
		Object{},
	}
}

//ref is a human name for a commit
type Ref struct {
	Name   string
	Object        //i think it is a commit object
	Child  Object //contain some info for this ref
}

func (ref *Ref) IsPeeled() bool {
	return strings.HasSuffix(ref.Name, PEELED_FLAG)
}
