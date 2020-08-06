package server

import (
	"github.com/aceld/zinx/ziface"
	"github.com/aceld/zinx/zlog"
	"github.com/aceld/zinx/znet"
	"zinx-mj/database"
	"zinx-mj/network/handler"
	"zinx-mj/network/protocol"
)

type Server struct {
	core ziface.IServer
}

func NewServer() *Server {
	return &Server{}
}

func (s *Server) initRoute() error {
	s.core.AddRouter(uint32(protocol.PROTOID_CS_LOGIN), &handler.Login{})
	s.core.AddRouter(uint32(protocol.PROTOID_CS_CREATE_ROOM), &handler.CreateRoom{})
	return nil
}

func (s *Server) Init() error {
	s.core = znet.NewServer()
	// 注册路由
	if err := s.initRoute(); err != nil {
		zlog.Errorf("init route failed")
	}

	// 连接db
	if err := database.Connect(); err != nil {
		zlog.Errorf("init database failed")
		return err
	}

	return nil
}

func (s *Server) Run() {
	s.core.Serve()
}
