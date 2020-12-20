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
	frameTick *time.Ticker
}

func NewServer() *Server {
	return &Server{}
}

func (s *Server) initRoute() error {
	s.core.AddRouter(uint32(protocol.PROTOID_CS_LOGIN), &handler.Login{})
	s.core.AddRouter(uint32(protocol.PROTOID_CS_CREATE_TABLE), &handler.CreateTable{})
	s.core.AddRouter(uint32(protocol.PROTOID_CS_JOIN_TABLE), &handler.JoinTable{})
	s.core.AddRouter(uint32(protocol.PROTOID_CS_PLAYER_OPERATE), &handler.PlayerOperate{})
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

	s.frameTick = time.NewTicker(1000 * time.Millisecond / serverFrame)
	return nil
}

func (s *Server) Run() {
	go s.FixUpdate()
	s.core.Serve()
}

func (s *Server) FixUpdate() {
	cur := time.Now()
	for range s.frameTick.C {
		delta := time.Since(cur)
		cur = time.Now()

		tablemgr.Update(delta)
	}
}
