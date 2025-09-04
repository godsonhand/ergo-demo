package main

import (
	"ergo.services/ergo/act"
	"ergo.services/ergo/gen"
)

type AsyncCallback func(error, any)

type context struct {
	ref gen.Ref
	pid gen.PID
	cb  AsyncCallback
}

type ActorBase struct {
	act.Actor
	session int64
	ctxs    map[int64]*context
}

func (a *ActorBase) Init(args ...any) error {
	a.ctxs = make(map[int64]*context)
	base, ok := a.Behavior().(IActorBase)
	if !ok {
		return ErrTypeAssert
	}
	return base.ActorInit(args...)
}

func (a *ActorBase) HandleMessage(from gen.PID, message any) error {
	switch m := message.(type) {
	case AsyncMessage:
		a.onAsyncCall(from, m.Session, m.Message)
		return nil
	case AyncResult:
		a.onAsyncResult(m.Session, m.Message, m.Err)
		return nil
	}
	base := a.Behavior().(IActorBase)
	return base.HandleCustomMessage(from, message)
}

func (a *ActorBase) HandleCall(from gen.PID, ref gen.Ref, request any) (any, error) {
	base := a.Behavior().(IActorBase)
	return base.HandleCustomCall(from, ref, request)
}

func (a *ActorBase) HandleCustomMessage(from gen.PID, message any) error {
	return ErrUnknown
}

func (a *ActorBase) HandleCustomCall(from gen.PID, ref gen.Ref, request any) (any, error) {
	return nil, ErrUnknown
}

func (a *ActorBase) AsyncCall(to any, msg any, cb AsyncCallback) error {
	a.session++
	ctx := &context{
		cb: cb,
	}

	err := a.Send(to, AsyncMessage{
		Session: a.session,
		Message: msg,
	})

	if err == nil {
		a.ctxs[a.session] = ctx
	}
	return err
}

func (a *ActorBase) DoAsyncCall(from gen.PID, msg any, resp AsyncResponse) {
	resp(nil, ErrUnknown)
}

func (a *ActorBase) onAsyncCall(from gen.PID, session int64, msg any) {
	call := false
	resp := func(ret any, err error) error {
		if call {
			return ErrDupResp
		}
		call = true
		return a.Send(from, AyncResult{
			Session: session,
			Err:     err,
			Message: ret,
		})
	}
	base := a.Behavior().(IActorBase)
	base.DoAsyncCall(from, msg, resp)
}

func (a *ActorBase) onAsyncResult(session int64, msg any, err error) {
	ctx, ok := a.ctxs[session]
	if !ok {
		return
	}
	delete(a.ctxs, session)
	ctx.cb(err, msg)
}
