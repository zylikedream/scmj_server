package tablemgr

import (
	"fmt"
	"github.com/aceld/zinx/zlog"
	"github.com/golang/protobuf/proto"
	"zinx-mj/game/gamedefine"
	"zinx-mj/game/rule/scmjrule"
	"zinx-mj/game/table/itable"
	"zinx-mj/game/table/scmjtable"
	"zinx-mj/network/protocol"
	"zinx-mj/player"
	"zinx-mj/player/playermgr"
	"zinx-mj/util"
)

var tables map[int]itable.IMjTable

func init() {
	tables = make(map[int]itable.IMjTable)
}

func CreateTable(pid player.PID, tableType int, message proto.Message) (itable.IMjTable, error) {
	var mjtable itable.IMjTable
	var err error
	switch tableType {
	case gamedefine.TABLE_TYPE_SCMJ:
		mjtable, err = createScmjTable(pid, message)
	default:
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return mjtable, nil
}

func createScmjTable(pid player.PID, req proto.Message) (itable.IMjTable, error) {
	msg, ok := req.(*protocol.CsCreateScmjTable)
	if !ok {
		zlog.Error("wrong message type %T", req)
		return nil, fmt.Errorf("wrong message type %T", req)
	}
	ply := playermgr.GetPlayerByPid(pid)
	if ply == nil {
		zlog.Errorf("can't find ply, pid=%d", pid)
		return nil, fmt.Errorf("can't find ply, pid=%d", pid)
	}
	master := util.PackTablePlayerDataFromPly(ply)
	ruleData := &scmjrule.ScmjRuleData{}
	ruleData.UnpackFromPBMsg(msg.GetRule())
	table, err := scmjtable.NewScmjTable(master, ruleData)
	if err != nil {
		return nil, fmt.Errorf("new scmj table failed")
	}
	return table, nil
}
