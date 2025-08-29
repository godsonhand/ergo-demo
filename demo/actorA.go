package main

import (
	"ergo.services/ergo/act"
	"ergo.services/ergo/gen"
)

type actorA struct {
	act.Actor
}

func factoryA() gen.ProcessBehavior {
	return &actorA{}
}

func (a *actorA) HandleCall(from gen.PID, ref gen.Ref, request any) (any, error) {
	return retMessage{}, nil
}
