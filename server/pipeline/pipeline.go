package pipeline

import (
	"container/list"
	"github.com/panjf2000/gnet/v2"
)

type Handler struct {
	Func func(c gnet.Conn) (action gnet.Action)

	name string
}

type ChannelPipeline struct {
	pipe *list.List
}

func New() *ChannelPipeline {
	return &ChannelPipeline{
		pipe: list.New(),
	}
}

func (p *ChannelPipeline) AddFirst(name string, handler *Handler) {
	handler.name = name
	p.pipe.PushFront(handler)
}

func (p *ChannelPipeline) AddLast(name string, handler *Handler) {
	handler.name = name
	p.pipe.PushBack(handler)
}

func (p *ChannelPipeline) AddBefore(name string, handler *Handler) {
	p.findInvoke(name, func(e *list.Element) {
		p.pipe.InsertBefore(handler, e)
	})
}

func (p *ChannelPipeline) AddAfter(name string, handler *Handler) {
	p.findInvoke(name, func(e *list.Element) {
		p.pipe.InsertAfter(handler, e)
	})
}

func (p *ChannelPipeline) Fire(c gnet.Conn) gnet.Action {
	for el := p.pipe.Front(); el != nil; el = el.Next() {
		if action := el.Value.(*Handler).Func(c); action != gnet.None {
			return action
		}
	}
	return gnet.None
}

func (p *ChannelPipeline) findInvoke(name string, f func(e *list.Element)) {
	for el := p.pipe.Front(); el != nil; el = el.Next() {
		if el.Value.(*Handler).name == name {
			f(el)
		}
	}
}
