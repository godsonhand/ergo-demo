package main

import (
	"errors"

	"ergo.services/ergo/gen"
)

type (
	AsyncMessage struct {
		Session int64
		Message any
	}

	AyncResult struct {
		Session int64
		Err     error
		Message any
	}

	doStartLogin struct{}
	doStartKick  struct{}
	loginMessage struct {
		err error
		pid gen.PID
	}
	logoutMessage struct {
		err error
		pid gen.PID
	}
	onlineMessage struct {
		err error
		pid gen.PID
	}
	offlineMessage struct {
		err error
		pid gen.PID
	}
	kickMessage struct{}
)

var (
	emptyPID = gen.PID{}
)

var (
	ErrUnknown     = errors.New("unknown")
	ErrCxtNotFound = errors.New("ctx not found")
	ErrTypeAssert  = errors.New("type assert fail")
	ErrLogin       = errors.New("login fail")
	ErrLogout      = errors.New("logout fail")
	ErrOnline      = errors.New("online fail")
	ErrOffline     = errors.New("offline fail")
	ErrDupResp     = errors.New("duplcate response")
)
