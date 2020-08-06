package main

import (
	"github.com/aceld/zinx/zlog"
	"zinx-mj/server"
)

func main() {
	svr := server.NewServer()
	if svr.Init() != nil {
		zlog.Errorf("svr init failed")
		return
	}
	svr.Run()
}
