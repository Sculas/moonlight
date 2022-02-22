package pipeline

import (
	"container/list"
	"github.com/panjf2000/gnet/v2"
)

type baseHandler interface{}

type InboundHandler interface {
	baseHandler
	Read(c gnet.Conn) (action gnet.Action) // TODO
}

type OutboundHandler interface {
	baseHandler
	Write() // TODO
}

type handler struct {
	h    baseHandler
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

func (p *ChannelPipeline) AddFirst(name string, h baseHandler) {
	p.pipe.PushFront(&handler{
		h:    h,
		name: name,
	})
}

func (p *ChannelPipeline) AddLast(name string, h baseHandler) {
	p.pipe.PushBack(&handler{
		h:    h,
		name: name,
	})
}

func (p *ChannelPipeline) AddBefore(name string, handler baseHandler) {
	p.findInvoke(name, func(e *list.Element) {
		p.pipe.InsertBefore(handler, e)
	})
}

func (p *ChannelPipeline) AddAfter(name string, handler baseHandler) {
	p.findInvoke(name, func(e *list.Element) {
		p.pipe.InsertAfter(handler, e)
	})
}

func (p *ChannelPipeline) Read(c gnet.Conn) gnet.Action {
	for el := p.pipe.Front(); el != nil; el = el.Next() {
		eh := el.Value.(*handler)
		if h, ok := eh.h.(InboundHandler); ok {
			if action := h.Read(c); action != gnet.None {
				return action
			}
		}
	}
	return gnet.None
}

func (p *ChannelPipeline) findInvoke(name string, f func(e *list.Element)) {
	for el := p.pipe.Front(); el != nil; el = el.Next() {
		if el.Value.(*handler).name == name {
			f(el)
		}
	}
}
