package server

import (
	"time"
	"zinx-mj/database"
	"zinx-mj/game/table/tablemgr"
	"zinx-mj/network/handler"
	"zinx-mj/network/protocol"

	"github.com/aceld/zinx/ziface"
	"github.com/aceld/zinx/zlog"
	"github.com/aceld/zinx/znet"
)

const serverFrame = 20

type Server struct {
	core      ziface.IServer
	frameTick <-chan time.Time
}

func NewServer() *Server {
	return &Server{}
}

func (s *Server) initRoute() error {
	s.core.AddRouter(uint32(protocol.PROTOID_CS_LOGIN), &handler.Login{})
	s.core.AddRouter(uint32(protocol.PROTOID_CS_CREATE_TABLE), &handler.CreateTable{})
	s.core.AddRouter(uint32(protocol.PROTOID_CS_JOIN_TABLE), &handler.JoinTable{})
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

	s.frameTick = time.Tick(1000 * time.Millisecond / serverFrame)
	return nil
}

func (s *Server) Run() {
	go s.core.Serve()
	go s.FixUpdate()
}

func (s *Server) FixUpdate() {
	cur := time.Now()
	for {
		<-s.frameTick
		delta := time.Since(cur)
		cur = time.Now()

		tablemgr.Update(delta)
	}
}
