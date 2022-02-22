package pipeline

import (
	"container/list"
	"fmt"
	"github.com/panjf2000/gnet/v2"
)

const (
	ErrEOP = "pipeline: an exception occurred in the pipeline, but was not handled by any of the handlers"
)

type baseHandler interface {
	// Exception gets invoked when an exception occurs in the pipeline.
	// If the exception was properly handled, return true, otherwise return false.
	Exception(c gnet.Conn, err error) (bool, gnet.Action)
}

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

func (p *ChannelPipeline) Fire(c gnet.Conn) gnet.Action {
	for el := p.pipe.Front(); el != nil; el = el.Next() { // loop through all handlers
		eh := el.Value.(*handler)
		action, err := p.invokeHandler(eh, c) // invoke the handler
		if action == gnet.Close {
			return action // we're being told to close, no reason to continue
		}
		if err != nil { // an exception occurred
			// propagate exception down the pipeline from the next handler
			handled := false
			for cel := el.Next(); cel != nil; cel = cel.Next() {
				ceh := cel.Value.(*handler)
				handled, action = ceh.h.Exception(c, err)
				if handled {
					// handled successfully, stop propagating
					if action == gnet.Close {
						return action
					}
					break
				}
			}
			if !handled {
				// no handler handled the exception, so we have to close the connection
				fmt.Println(ErrEOP) // TODO
				return gnet.Close
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

func (ChannelPipeline) invokeHandler(h *handler, c gnet.Conn) (action gnet.Action, err error) {
	defer func() {
		if r := recover(); r != nil {
			if e, ok := r.(error); ok {
				err = e
			} else {
				err = fmt.Errorf("%v", r)
			}
		}
	}()

	switch h.h.(type) {
	case InboundHandler:
		action = h.h.(InboundHandler).Read(c)
	case OutboundHandler:
		h.h.(OutboundHandler).Write()
	}
	return
}
