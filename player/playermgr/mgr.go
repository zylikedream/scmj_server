package playermgr

import (
	"time"
	"zinx-mj/database"
	"zinx-mj/player"

	"github.com/aceld/zinx/zlog"
	"github.com/zheng-ji/goSnowFlake"
)

type PID = player.PID

var pcache *playerCache      // 玩家缓存实例
var iw *goSnowFlake.IdWorker // 唯一id生成器实例

func init() {
	// 初始化唯一id生成器
	iw, _ = goSnowFlake.NewIdWorker(1)
	pcache = NewPlayerCache()
}

/*
 * Descrp: 生成一个唯一id
 * Create: zhangyi 2020-6-20 01:03:48
 */
func GenUuid() (int64, error) {
	return iw.NextId()
}

/*
 * Descrp: 生成唯一的pid
 * Create: zhangyi 2020-6-20 01:04:03
 */
func genPid() (PID, error) {
	uuid, err := GenUuid()
	if err != nil {
		zlog.Errorf("genPid failed, err=%s\n", err)
		return 0, err
	}
	return PID(uuid), nil
}

/*
 * Descrp: 创建一个与账号关联的player对象
 * Create: zhangyi 2020-6-20 01:04:34
 */
func CreatePlayer(account string) (*player.Player, error) {
	pid, err := genPid()
	if err != nil {
		zlog.Errorf("create player failed, account=%s, err=%s\n", account, err)
		return nil, err
	}
	p := player.New()
	p.Pid = pid
	p.Account = account
	p.Name = account
	ts := time.Now().Unix()
	p.CreateTime = ts
	p.LastLoginTime = ts

	if err = database.GetDB().SavePlayer(p); err != nil {
		zlog.Errorf("create player failed, account=%s, err=%s\n", account, err)
		return nil, err
	}
	if err = pcache.AddPlayer(p); err != nil {
		zlog.Errorf("create player failed, account=%s, err=%s\n", account, err)
		return nil, err
	}
	zlog.Infof("create player success, account=%s\n", account)
	return p, nil
}

/*
 * Descrp: 通过账号, 取得玩家数据，首先从缓存中取，如果没有就从数据库中拉取
 * Create: zhangyi 2020-06-20 02:06:35
 */
func GetPlayerByAccount(account string) (*player.Player, error) {
	ply := pcache.GetPlayerByAccount(account)
	if ply != nil {
		return ply, nil
	}
	ply, err := database.GetDB().LoadPlayer(account)
	if err != nil {
		zlog.Errorf("get player by account failed, account=%s", account)
		return nil, err
	}
	if err = pcache.AddPlayer(ply); err != nil {
		zlog.Errorf("get player by account failed, account=%s", account)
		return nil, err
	}
	return ply, nil
}

func GetPlayerByPid(pid player.PID) *player.Player {
	return pcache.GetPlayerByPid(pid)
}
