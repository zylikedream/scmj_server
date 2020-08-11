package handler

import (
	"github.com/aceld/zinx/ziface"
	"github.com/aceld/zinx/zlog"
	"github.com/aceld/zinx/znet"
	"github.com/golang/protobuf/proto"
	"zinx-mj/game/gamedefine"
	"zinx-mj/game/table/tablemgr"
	"zinx-mj/network/protocol"
	"zinx-mj/player"
	"zinx-mj/util"
)

type Table struct {
	znet.BaseRouter
}

func (c *Table) Handle(request ziface.IRequest) { //处理conn业务的方法
	data := request.GetData()
	req := &protocol.CsCreateScmjTable{}
	if err := proto.Unmarshal(data, req); err != nil {
		zlog.Errorf("unpack login proto failed\n")
		return
	}
	zlog.Debugf("create room: pid=%d, rule=%v\n", req.Pid, req.Rule)

	pid := util.GetConnPid(request.GetConnection())
	reply := c.doCreateTable(pid, req, request.GetConnection())
	data, err := proto.Marshal(reply)
	if err != nil {
		zlog.Error("Marshal packet failed")
		return
	}
	if err = request.GetConnection().SendMsg(uint32(protocol.PROTOID_SC_TABLE_INFO), data); err != nil {
		zlog.Errorf("send msg failed, err=%s", err)
		return
	}
}

func (c *Table) doCreateTable(pid player.PID, req *protocol.CsCreateScmjTable, conn ziface.IConnection) *protocol.ScScmjTableInfo {
	reply := &protocol.ScScmjTableInfo{}
	table, err := tablemgr.CreateTable(pid, gamedefine.TABLE_TYPE_SCMJ, req)

	if err != nil {
		zlog.Errorf("create table failed")
		return nil
	}
	// 得到玩家已有对象
	return reply
}
