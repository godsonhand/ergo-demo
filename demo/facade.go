package main

import "ergo.services/ergo/gen"

type AsyncResponse func(any, error) error

type IActorBase interface {
	ActorInit(args ...any) error
	HandleCustomMessage(from gen.PID, message any) error
	HandleCustomCall(from gen.PID, ref gen.Ref, request any) (any, error)
	DoAsyncCall(from gen.PID, msg any, resp AsyncResponse)
}
