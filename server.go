package main

import (
	"github.com/aceld/zinx/znet"
	"zinx-mj/database"
	"zinx-mj/network/handler"
	"zinx-mj/network/protocol"
)

func main() {
	if err := database.Init(); err != nil {
		return
	}
	svr := znet.NewServer()
	svr.AddRouter(uint32(protocol.PROTOID_CS_LOGIN_ID), &handler.Login{})
	svr.Serve()
}
