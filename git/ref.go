package git

type Refs map[string]Ref

//ref is a human name for a commit
type Ref struct {
	Obj   Object //i think it is a commit object
	Child Object //contain some info for this ref
}

type Object struct {
	Id       string //SHA-1 40 char
	SelfType string //commit, blob, tree
}
