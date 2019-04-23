package game

import "game/protoMsg"

func (p *Player) newResChange() {
	p.resChange = &ResChange{
		items: make(map[uint32]*ItemChange),
	}
}

type ResChange struct {
	items map[uint32]*ItemChange
}

type ItemChange struct {
	id     uint32
	now    uint64
	change uint64
	source uint32
}

func (p *ResChange) reset() {
	p.items = map[uint32]*ItemChange{}
}
func (p *ResChange) addChange(itemId uint32, now, change uint64, source uint32) {
	info := p.items[itemId]
	if info == nil {
		info = &ItemChange{
			id: itemId,
		}
		p.items[itemId] = info
	}
	info.now = now
	info.change += change
	info.source = source
}
func (p *ResChange) Empty() bool {
	return len(p.items) == 0
}
func (p *Player) ResChange() {
	if p.resChange.Empty() {
		return
	}
	r := &protoMsg.S_ResChangeNotify{}
	for k, v := range p.resChange.items {
		if v.change == 0 {
			continue
		}
		r.Info = append(r.Info, &protoMsg.ResChange{
			Kind:   k,
			Now:    v.now,
			Change: v.change,
			Source: v.source,
		})
	}
	if len(r.Info) == 0 {
		return
	}
	// this.Debug("ResChange r:", r)
	p.SendToClient(protoMsg.S_CMD_S_RESCHANGENOTIFY, r)
	p.resChange.reset()
}
