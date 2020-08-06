package handler

import (
	"github.com/aceld/zinx/ziface"
	"github.com/aceld/zinx/zlog"
	"github.com/aceld/zinx/znet"
	"github.com/golang/protobuf/proto"
	"go.mongodb.org/mongo-driver/mongo"
	"zinx-mj/network/protocol"
	"zinx-mj/player/playermgr"
)

type Login struct {
	znet.BaseRouter
}

func (l *Login) Handle(request ziface.IRequest) { //处理conn业务的方法
	data := request.GetData()
	login := &protocol.CsLogin{}
	if err := proto.Unmarshal(data, login); err != nil {
		zlog.Errorf("unpack login proto failed\n")
		return
	}
	zlog.Debugf("user login: username=%s, password=%s\n", login.Account, login.Password)

	reply := l.doLogin(login, request.GetConnection())
	data, err := proto.Marshal(reply)
	if err != nil {
		zlog.Error("Marshal packet failed")
		return
	}
	if err = request.GetConnection().SendMsg(uint32(protocol.PROTOID_SC_LOGIN), data); err != nil {
		zlog.Errorf("send msg failed, err=%s", err)
		return
	}
}

func (*Login) doLogin(login *protocol.CsLogin, conn ziface.IConnection) *protocol.ScLogin {
	reply := &protocol.ScLogin{
		Pid:     0,
		Account: login.Account,
	}
	account := login.Account
	// 得到玩家已有对象
	ply, err := playermgr.GetPlayerByAccount(account)
	if err == mongo.ErrNoDocuments { // 没有找到记录
		// 创建一个新号
		ply, err = playermgr.CreatePlayer(account)
	}
	// 如果没有就创建一个
	if err != nil {
		zlog.Errorf("get player account failed, account=%s\n", account)
		return reply
	}
	// 绑定连接
	ply.Conn = conn

	reply.Pid = ply.Pid
	reply.Name = ply.Name
	reply.Sex = int32(ply.Sex)
	reply.RoomCard = ply.RoomCard
	return reply
}
