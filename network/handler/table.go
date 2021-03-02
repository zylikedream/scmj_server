package handler

import (
	"fmt"
	"zinx-mj/game/gamedefine"
	"zinx-mj/game/table/tablemgr"
	"zinx-mj/game/table/tableplayer"
	"zinx-mj/mjerror"
	"zinx-mj/network/protocol"
	"zinx-mj/player"
	"zinx-mj/player/playermgr"
	"zinx-mj/util"

	"github.com/aceld/zinx/ziface"
	"github.com/aceld/zinx/zlog"
	"github.com/aceld/zinx/znet"
	"google.golang.org/protobuf/proto"
)

func PackTablePlayerDataFromPly(pid player.PID) (*tableplayer.TablePlayerData, error) {
	ply := playermgr.GetPlayerByPid(pid)
	if ply == nil {
		return nil, fmt.Errorf("pid=%d:%w", pid, mjerror.ErrPlyNotFound)
	}
	return &tableplayer.TablePlayerData{
		Pid:  ply.Pid,
		Name: ply.Name,
	}, nil
}

type CreateTable struct {
	znet.BaseRouter
}

func (c *CreateTable) Handle(request ziface.IRequest) { //处理conn业务的方法
	data := request.GetData()
	req := &protocol.CsCreateScmjTable{}
	if err := proto.Unmarshal(data, req); err != nil {
		zlog.Errorf("unpack create table proto failed, err:%s\n", err)
		return
	}
	pid := util.GetConnPid(request.GetConnection())

	if err := c.doCreateTable(pid, req, request.GetConnection()); err != nil {
		zlog.Errorf("create table failed, err:%s", err)
		return
	}
}

func (c *CreateTable) doCreateTable(pid player.PID, req *protocol.CsCreateScmjTable, conn ziface.IConnection) error {
	master, err := PackTablePlayerDataFromPly(pid)
	if err != nil {
		return err
	}
	table, err := tablemgr.CreateTable(master, gamedefine.TABLE_TYPE_SCMJ, req)
	if err != nil {
		return err
	}
	ply := playermgr.GetPlayerByPid(pid)
	if ply == nil {
		return mjerror.ErrPlyNotFound
	}
	ply.SetTableID(table.GetID())
	// 得到玩家已有对象
	return nil
}

type JoinTable struct {
	znet.BaseRouter
}

func (c *JoinTable) Handle(request ziface.IRequest) {
	data := request.GetData()
	req := &protocol.CsJoinTable{}
	if err := proto.Unmarshal(data, req); err != nil {
		zlog.Errorf("unpack join table proto failed, err:%s\n", err)
		return
	}
	tb := tablemgr.GetTable(uint32(req.TableId))
	if tb == nil {
		zlog.Errorf("get table failed, id:%d", req.TableId)
		return
	}

	pid := util.GetConnPid(request.GetConnection())
	plyData, err := PackTablePlayerDataFromPly(pid)
	if err != nil {
		return
	}
	_, err = tb.PlayerJoin(plyData, gamedefine.TABLE_IDENTIY_PLAYER)
	if err != nil {
		zlog.Errorf("player join failed, err:%s", err)
		return
	}
	ply := playermgr.GetPlayerByPid(pid)
	if ply == nil {
		return
	}
	ply.SetTableID(tb.GetID())

}

type PlayerReady struct {
	znet.BaseRouter
}

func (c *PlayerReady) Handle(request ziface.IRequest) {
	data := request.GetData()
	req := &protocol.CsPlayerReady{}
	if err := proto.Unmarshal(data, req); err != nil {
		zlog.Errorf("unpack join table proto failed, err:%s\n", err)
		return
	}
	pid := util.GetConnPid(request.GetConnection())
	ply := playermgr.GetPlayerByPid(pid)
	if ply == nil {
		return
	}
	tb := tablemgr.GetTable(ply.TableID)
	if tb == nil {
		zlog.Errorf("get table failed, id:%d", ply.TableID)
		return
	}
	tb.SetReady(pid, req.Ready)
}
