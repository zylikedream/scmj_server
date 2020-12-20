package handler

import (
	"zinx-mj/game/gamedefine"
	"zinx-mj/game/table/tablemgr"
	"zinx-mj/network/protocol"
	"zinx-mj/player"
	"zinx-mj/player/playermgr"
	"zinx-mj/util"

	"github.com/aceld/zinx/ziface"
	"github.com/aceld/zinx/zlog"
	"github.com/aceld/zinx/znet"
	"google.golang.org/protobuf/proto"
)

type CreateTable struct {
	znet.BaseRouter
}

func (c *CreateTable) Handle(request ziface.IRequest) { //处理conn业务的方法
	data := request.GetData()
	req := &protocol.CsCreateScmjTable{}
	if err := proto.Unmarshal(data, req); err != nil {
		zlog.Errorf("unpack create table proto failed, err=%s\n", err)
		return
	}
	pid := util.GetConnPid(request.GetConnection())
	//zlog.Debugf("create table: pid=%d, rule=%v\n", pid, req.Rule)

	reply := c.doCreateTable(pid, req, request.GetConnection())
	if err := util.SendMsg(pid, protocol.PROTOID_SC_TABLE_INFO, reply); err != nil {
		zlog.Errorf("send msg failed, err=%s", err)
		return
	}
}

func (c *CreateTable) doCreateTable(pid player.PID, req *protocol.CsCreateScmjTable, conn ziface.IConnection) *protocol.ScScmjTableInfo {

	master, err := util.PackTablePlayerDataFromPly(pid)
	if err != nil {
		return nil
	}
	table, err := tablemgr.CreateTable(master, gamedefine.TABLE_TYPE_SCMJ, req)
	if err != nil {
		return nil
	}
	ply := playermgr.GetPlayerByPid(pid)
	if ply == nil {
		return nil
	}
	ply.SetTableID(table.GetID())
	reply := table.PackToPBMsg().(*protocol.ScScmjTableInfo)
	// 得到玩家已有对象
	return reply
}

type JoinTable struct {
	znet.BaseRouter
}

func (c *JoinTable) Handle(request ziface.IRequest) {
	data := request.GetData()
	req := &protocol.CsJoinTable{}
	if err := proto.Unmarshal(data, req); err != nil {
		zlog.Errorf("unpack join table proto failed, err=%s\n", err)
		return
	}
	tb := tablemgr.GetTable(req.TableId)
	if tb == nil {
		zlog.Errorf("get table failed, id=%d", req.TableId)
		return
	}

	pid := util.GetConnPid(request.GetConnection())
	plyData, err := util.PackTablePlayerDataFromPly(pid)
	if err != nil {
		return
	}
	_, err = tb.Join(plyData, gamedefine.TABLE_IDENTIY_PLAYER)
	if err != nil {
		return
	}
	ply := playermgr.GetPlayerByPid(pid)
	if ply == nil {
		return
	}
	ply.SetTableID(tb.GetID())

	reply := tb.PackToPBMsg().(*protocol.ScScmjTableInfo)

	if err = util.SendMsg(pid, protocol.PROTOID_SC_TABLE_INFO, reply); err != nil {
		zlog.Errorf("send msg failed, err=%s", err)
		return
	}

}
