package main

import (
	"math/rand"
	"sync"
	"time"

	"ergo.services/ergo/gen"
)

type agent struct {
	ActorBase
	cpid     gen.PID
	wd       *sync.WaitGroup
	isLogin  bool
	isOnline bool
	apid     gen.PID
}

func factoryAgent() gen.ProcessBehavior {
	return &agent{}
}

func (a *agent) ActorInit(args ...any) error {
	a.cpid, _ = args[0].(gen.PID)
	a.wd, _ = args[1].(*sync.WaitGroup)
	return nil
}

func (a *agent) HandleCustomMessage(from gen.PID, message any) error {
	switch message.(type) {
	case doStartLogin:
		return a.startLogin()
	case kickMessage:
		return a.kick()
	case doStartKick:
		return a.startKick()
	}
	return ErrUnknown
}

func (a *agent) login() error {
	a.Log().Info("[agent]login")
	randNum := rand.Intn(200) + 100
	_, err := a.SendAfter(a.PID(), doStartLogin{}, time.Millisecond*time.Duration(randNum))
	return err
}

func (a *agent) kick() error {
	a.Log().Info("[agent]kick")
	randNum := rand.Intn(200) + 100
	_, err := a.SendAfter(a.PID(), doStartKick{}, time.Millisecond*time.Duration(randNum))
	return err
}

func (a *agent) handleError(err error) error {
	defer func() {
		if r := recover(); r != nil {
		}
	}()

	a.Log().Error("[agent]error:%v", err)
	if err != gen.ErrTimeout {
		return a.login()
	}
	a.wd.Done()
	return err
}

func (a *agent) startLogin() error {
	a.Log().Info("[agent]startLogin")
	ret, err := a.Call(a.cpid, loginMessage{
		pid: a.PID(),
	})

	if err != nil {
		return a.handleError(err)
	}

	msg, _ := ret.(loginMessage)
	if msg.err != nil {
		return a.handleError(msg.err)
	}

	a.isLogin = true
	a.Log().Info("[agent]login succ")

	ret, err = a.Call(msg.pid, onlineMessage{
		pid: a.PID(),
	})

	if err != nil {
		return a.handleError(err)
	}

	msg2, _ := ret.(onlineMessage)
	if msg2.err != nil {
		return a.handleError(msg2.err)
	}

	a.apid = msg.pid
	a.isOnline = true
	a.Log().Info("[agent]online succ")
	return nil
}

func (a *agent) startKick() error {
	a.Log().Info("[agent]startKick")
	if a.isOnline {
		a.isOnline = false
		a.Log().Info("[agent]offline")
		ret, err := a.Call(a.apid, offlineMessage{
			pid: a.PID(),
		})
		if err != nil {
			return a.handleError(err)
		}

		msg, _ := ret.(offlineMessage)
		if msg.err != nil {
			return a.handleError(msg.err)
		}
		a.Log().Info("[agent]offline succ")
	}

	if a.isLogin {
		a.isLogin = false
		a.Log().Info("[agent]logout")
		ret, err := a.Call(a.cpid, logoutMessage{
			pid: a.PID(),
		})
		if err != nil {
			return a.handleError(err)
		}

		msg, _ := ret.(logoutMessage)
		if msg.err != nil {
			return a.handleError(msg.err)
		}
		a.Log().Info("[agent]logout succ")
	}

	return a.login()
}
