package main

import (
	"ergo.services/ergo/gen"
)

type user struct {
	ActorBase
	pid gen.PID
}

func factoryUser() gen.ProcessBehavior {
	return &user{}
}

func (a *user) ActorInit(args ...any) error {
	return nil
}

func (a *user) HandleCustomCall(from gen.PID, ref gen.Ref, request any) (any, error) {
	switch request.(type) {
	case onlineMessage:
		return a.online(from)
	case offlineMessage:
		return a.offline(from)
	}
	return nil, ErrUnknown
}

func (a *user) DoAsyncCall(from gen.PID, message any, resp AsyncResponse) {
	switch m := message.(type) {
	case loginMessage:
		resp(a.login(m.pid))
		return
	case logoutMessage:
		resp(a.logout(m.pid))
		return
	}
	resp(nil, ErrUnknown)
}

func (a *user) login(from gen.PID) (any, error) {
	ret := loginMessage{}
	if a.pid != emptyPID {
		a.Log().Warning("[user]login fail, curr:%v, from:%v", a.pid, from)
		ret.err = ErrLogin
	}

	a.pid = from
	ret.pid = a.PID()
	a.Log().Info("[user]login succ, from:%v", from)
	return ret, nil
}

func (a *user) logout(from gen.PID) (any, error) {
	ret := logoutMessage{}
	if a.pid != from {
		a.Log().Warning("[user]logout fail, curr:%v, from:%v", a.pid, from)
		ret.err = ErrLogout
	} else {
		a.pid = emptyPID
		a.Log().Info("[user]logout succ, from:%v", from)
	}
	return ret, nil
}

func (a *user) online(from gen.PID) (any, error) {
	ret := onlineMessage{}
	if a.pid != from {
		a.Log().Warning("[user]online fail, curr:%v, from:%v", a.pid, from)
		ret.err = ErrOnline
	} else {
		a.Log().Info("[user]online, from:%v", from)
	}
	return ret, nil
}

func (a *user) offline(from gen.PID) (any, error) {
	ret := offlineMessage{}
	if a.pid != from {
		a.Log().Warning("[user]offline fail, curr:%v, from:%v", a.pid, from)
		ret.err = ErrOffline
	} else {
		a.Log().Info("[user]offline, from:%v", from)
	}
	return ret, nil
}
