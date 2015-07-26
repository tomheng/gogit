package git

import "strings"

//Refs list
type Refs map[string]Ref

const (
	peeledSuffix = "^{}"
)

//NewRefs create ref list
func NewRefs() Refs {
	return make(Refs, 0)
}

//AddRef add ref to the list
func (ref Refs) AddRef(name, oid string) {
	ref[name] = Ref{
		name,
		Object{ID: oid},
		Object{},
	}
}

//Ref is a human name for a commit
type Ref struct {
	Name   string
	Object        //i think it is a commit object
	Child  Object //contain some info for this ref
}

//IsPeeled check if it is peeled
func (ref *Ref) IsPeeled() bool {
	return strings.HasSuffix(ref.Name, peeledSuffix)
}
