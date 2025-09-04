package main

import (
	"ergo.services/ergo/gen"
)

type mgr struct {
	ActorBase
	apid    gen.PID
	session int64
	ctxs    map[int64]*context
	curr    gen.PID
	wait    gen.PID
	waitRef gen.Ref
}

func factoryMgr() gen.ProcessBehavior {
	return &mgr{
		ctxs: make(map[int64]*context),
	}
}

func (a *mgr) ActorInit(args ...any) error {
	a.apid = args[0].(gen.PID)
	return nil
}

func (a *mgr) HandleCustomCall(from gen.PID, ref gen.Ref, request any) (any, error) {
	switch request.(type) {
	case loginMessage:
		return a.login(from, ref, request)
	case logoutMessage:
		return a.logout(from, ref, request)
	}
	return nil, ErrUnknown
}

func (a *mgr) requestLogin(from gen.PID, ref gen.Ref, request any) error {
	a.curr = from
	a.Log().Info("[mgr]request login from:%v", from)
	err := a.AsyncCall(a.apid, request, func(err error, msg any) {
		if err != nil {
			a.Log().Error("[mgr]response login from:%v, err:%v", from, err)
			a.SendResponseError(from, ref, err)
		} else {
			a.Log().Info("[mgr]response login from:%v", from)
			a.SendResponse(from, ref, msg)
		}
	})
	if err != nil {
		a.Log().Error("[mgr]request login err:%v from:%v", err, from)
	}
	return err
}

func (a *mgr) login(from gen.PID, ref gen.Ref, request any) (any, error) {
	ret := loginMessage{}
	if a.wait != emptyPID {
		a.Log().Warning("[mgr]login fail, from:%v, wait:%v", from, a.wait)
		ret.err = ErrLogin
		return ret, nil
	}

	if a.curr != emptyPID {
		a.wait = from
		a.waitRef = ref
		a.Log().Warning("[mgr]replace %v by %v", a.curr, from)
		a.Send(a.curr, kickMessage{})
		return nil, nil
	}

	err := a.requestLogin(from, ref, request)
	if err != nil {
		ret.err = err
		return ret, nil
	}
	return nil, nil
}

func (a *mgr) logout(from gen.PID, ref gen.Ref, request any) (any, error) {
	a.curr = gen.PID{}
	a.Log().Info("[mgr]request logout from:%v", from)
	ret := logoutMessage{}
	err := a.AsyncCall(a.apid, request, func(err error, msg any) {
		if err != nil {
			a.Log().Error("[mgr]response logout from:%v, err:%v", from, err)
			a.SendResponseError(from, ref, err)
		} else {
			a.Log().Info("[mgr]response logout from:%v", from)
			a.SendResponse(from, ref, msg)
		}

		if a.wait != emptyPID {
			a.Log().Info("[mgr]replace %v", a.wait)
			a.requestLogin(a.wait, a.waitRef, loginMessage{
				pid: a.wait,
			})
			a.wait = emptyPID
		}
	})
	if err != nil {
		ret.err = err
		a.Log().Error("[mgr]request logout err:%v from:%v", err, from)
		return ret, nil
	}
	return nil, nil
}
