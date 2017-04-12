package base

type Entity struct {
	id     uint64
	kind   uint32
	nodeID uint32
}

func (e *Entity) GetID() uint64 {
	return e.id
}
func (e *Entity) GetType() uint32 {
	return e.kind
}
func (e *Entity) GetNode() uint32 {
	return e.nodeID
}
