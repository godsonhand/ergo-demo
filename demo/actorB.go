package main

import (
	"math/rand"
	"sync"
	"time"

	"ergo.services/ergo/act"
	"ergo.services/ergo/gen"
)

type actorB struct {
	act.Actor
	pid gen.PID
	wd  *sync.WaitGroup
}

func factoryB() gen.ProcessBehavior {
	return &actorB{}
}

func (a *actorB) Init(args ...any) error {
	a.pid, _ = args[0].(gen.PID)
	a.wd, _ = args[1].(*sync.WaitGroup)
	return nil
}

func (a *actorB) startCall() {
	defer func() {
		if r := recover(); r != nil {
		}
	}()

	_, err := a.Call(a.pid, callMessage{})
	if err != nil {
		a.Log().Error("call %v err:%v", a.pid, err)
		a.wd.Done()
		return
	}

	randNum := rand.Intn(100) + 30
	time.Sleep(time.Millisecond * time.Duration(randNum))
	a.Send(a.PID(), doStartCall{})
}

func (a *actorB) HandleMessage(from gen.PID, message any) error {
	switch message.(type) {
	case doStartCall:
		a.startCall()
	}
	return nil
}
