package main

import (
	"math/rand"
	"sync"
	"time"

	"ergo.services/ergo"
	"ergo.services/ergo/gen"
	"ergo.services/logger/colored"
	"ergo.services/logger/rotate"
)

func main() {
	var options gen.NodeOptions
	// apps := []gen.ApplicationBehavior{
	// 	observer.CreateApp(observer.Options{}),
	// }

	// options.Applications = apps

	// disable default logger to get rid of multiple logging to the os.Stdout
	options.Log.DefaultLogger.Disable = true

	// add logger "colored".
	optionColored := colored.Options{TimeFormat: time.DateTime}
	loggerColored, err := colored.CreateLogger(optionColored)
	if err != nil {
		panic(err)
	}
	options.Log.Loggers = append(options.Log.Loggers, gen.Logger{Name: "cl", Logger: loggerColored})

	optionRotate := rotate.Options{
		TimeFormat: time.DateTime,
		Period:     time.Minute * time.Duration(10),
		Depth:      2,
	}
	loggerRotate, err := rotate.CreateLogger(optionRotate)
	if err != nil {
		panic(err)
	}
	options.Log.Loggers = append(options.Log.Loggers, gen.Logger{Name: "rl", Logger: loggerRotate})
	options.Log.Level = gen.LogLevelTrace

	// set network cookie
	options.Network.Cookie = "123"

	// starting node
	node, err := ergo.StartNode("node@localhost", options)
	if err != nil {
		panic(err)
	}

	apid, err := node.SpawnRegister("user", factoryUser, gen.ProcessOptions{})
	if err != nil {
		panic(err)
	}

	cpid, err := node.SpawnRegister("mgr", factoryMgr, gen.ProcessOptions{}, apid)
	if err != nil {
		panic(err)
	}

	rand.Seed(time.Now().UnixNano())
	var wd sync.WaitGroup
	wd.Add(1)
	callPIDs := make([]gen.PID, 3)
	for range 3 {
		p, err := node.Spawn(factoryAgent, gen.ProcessOptions{}, cpid, &wd)
		if err != nil {
			panic(err)
		}
		callPIDs = append(callPIDs, p)
	}

	for _, p := range callPIDs {
		node.Send(p, doStartLogin{})
	}

	wd.Wait()
	node.Stop()
}
