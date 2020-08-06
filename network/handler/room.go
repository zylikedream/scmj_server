package handler

import (
	"github.com/aceld/zinx/ziface"
	"github.com/aceld/zinx/zlog"
	"github.com/aceld/zinx/znet"
	"github.com/golang/protobuf/proto"
	"zinx-mj/network/protocol"
)

type CreateRoom struct {
	znet.BaseRouter
}

func (c *CreateRoom) Handle(request ziface.IRequest) { //处理conn业务的方法
	data := request.GetData()
	createRoom := &protocol.CsCreateRoom{}
	if err := proto.Unmarshal(data, createRoom); err != nil {
		zlog.Errorf("unpack login proto failed\n")
		return
	}
	zlog.Debugf("create room: pid=%d, rule=%v\n", createRoom.Pid, createRoom.Rule)

	reply := c.doCreateRoom(createRoom, request.GetConnection())
	data, err := proto.Marshal(reply)
	if err != nil {
		zlog.Error("Marshal packet failed")
		return
	}
	if err = request.GetConnection().SendMsg(uint32(protocol.PROTOID_SC_ROOM_INFO), data); err != nil {
		zlog.Errorf("send msg failed, err=%s", err)
		return
	}
}

func (c *CreateRoom) doCreateRoom(creatRoom *protocol.CsCreateRoom, conn ziface.IConnection) *protocol.ScRoomInfo {
	reply := &protocol.ScRoomInfo{}
	// 得到玩家已有对象
	return reply
}
