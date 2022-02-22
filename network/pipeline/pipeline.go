package pipeline

import (
	"container/list"
	"fmt"
	"github.com/panjf2000/gnet/v2"
	"github.com/sculas/moonlight/network/serde"
	"github.com/sculas/moonlight/server/client"
)

const (
	ErrEOP = "pipeline: an exception occurred in the pipeline, but was not handled by any of the handlers"
)

type baseHandler interface {
	// Exception gets invoked when an exception occurs in the pipeline.
	// If the exception was properly handled, return nil.
	// You may propagate the error to the next handler in the pipeline by returning the error or any non-nil error value.
	// If you try to handle the error but fail to do so, return your own error.
	Exception(c *client.Client, err error) (error, gnet.Action)
}

type InboundHandler interface {
	baseHandler
	Read(c *client.Client, buf *serde.ByteBuf) (action gnet.Action) // TODO
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
	c    *client.Client
}

func New(c *client.Client) *ChannelPipeline {
	return &ChannelPipeline{
		pipe: list.New(),
		c:    c,
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
	b, err := c.Next(-1)
	if err != nil {
		return gnet.Close // just give up
	}
	buf := serde.With(b)
	for el := p.pipe.Front(); el != nil; el = el.Next() { // loop through all handlers
		eh := el.Value.(*handler)
		action, err := p.invokeHandler(eh, p.c, buf) // invoke the handler
		if action == gnet.Close {
			return action // we're being told to close, no reason to continue
		}
		if err != nil { // an exception occurred
			// propagate exception down the pipeline from the next handler
			handled := false
			for cel := el.Next(); cel != nil; cel = cel.Next() {
				ceh := cel.Value.(*handler)
				err, action = ceh.h.Exception(p.c, err)
				if err == nil {
					// handled successfully, stop propagating
					if action == gnet.Close {
						return action
					}
					handled = true
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

func (ChannelPipeline) invokeHandler(h *handler, c *client.Client, buf *serde.ByteBuf) (action gnet.Action, err error) {
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
		action = h.h.(InboundHandler).Read(c, buf)
	case OutboundHandler:
		h.h.(OutboundHandler).Write()
	}
	return
}
