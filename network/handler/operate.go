package handler

import (
	"zinx-mj/game/table/tablemgr"
	"zinx-mj/game/table/tableoperate"
	"zinx-mj/network/protocol"
	"zinx-mj/player/playermgr"
	"zinx-mj/util"

	"github.com/aceld/zinx/ziface"
	"github.com/aceld/zinx/zlog"
	"github.com/aceld/zinx/znet"
	"google.golang.org/protobuf/proto"
)

type DoOperate struct {
	znet.BaseRouter
}

func (c *DoOperate) Handle(request ziface.IRequest) {
	data := request.GetData()
	req := &protocol.CsDoOperate{}
	if err := proto.Unmarshal(data, req); err != nil {
		zlog.Errorf("unpack join table proto failed, err=%s\n", err)
		return
	}

	pid := util.GetConnPid(request.GetConnection())
	ply := playermgr.GetPlayerByPid(pid)
	if ply == nil {
		return
	}

	tb := tablemgr.GetTable(ply.TableID)
	plyOperate := tableoperate.OperateCommand{
		OpType: int(req.OpType),
		OpData: tableoperate.OperateData{
			Card: int(req.Data.Card),
		},
	}
	if err := tb.OnPlyOperate(pid, plyOperate); err != nil {
		zlog.Errorf("operate failed, pid=%d, err=%s\n", pid, err)
		return
	}
}
