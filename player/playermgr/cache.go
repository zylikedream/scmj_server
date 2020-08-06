/*
 * Copyright (c) zhangyi All rights reserved. 2020
 * Descrp: 玩家缓存
 * Create: zhangyi 2020-06-20 01:33:35
 */
package playermgr

import (
	"fmt"
	"sync"
	"zinx-mj/player"
)

type playerCache struct {
	sync.RWMutex                               // 玩家管理器锁
	pidCache     map[player.PID]*player.Player // pid为key的缓存
	accountCache map[string]*player.Player     // account为key的管理器
}

func NewPlayerCache() *playerCache {
	return &playerCache{
		pidCache:     make(map[player.PID]*player.Player),
		accountCache: make(map[string]*player.Player),
	}
}

/*
 * Descrp: 向玩家缓存中添加玩家
 * Param: ply-玩家对象
 * Create: zhangyi 2020-06-20 01:06:54
 */
func (pc *playerCache) AddPlayer(ply *player.Player) error {
	pc.Lock()
	defer pc.Unlock()
	if pc.pidCache[ply.Pid] != nil {
		return fmt.Errorf("ply already in mgr, pid=%d", ply.Pid)
	}
	pc.pidCache[ply.Pid] = ply
	pc.accountCache[ply.Account] = ply
	return nil
}

/*
 * Descrp: 通过pid从缓存中获取玩家
 * Create: zhangyi 2020-06-20 01:33:52
 */
func (pc *playerCache) GetPlayerByPid(pid player.PID) *player.Player {
	pc.RLock()
	defer pc.RUnlock()
	return pc.pidCache[pid]
}

/*
 * Descrp: 通过账号从缓存中获取玩家
 * Create: zhangyi 2020-06-20 18:14:03
 */
func (pc *playerCache) GetPlayerByAccount(account string) *player.Player {
	pc.RLock()
	defer pc.RUnlock()
	return pc.accountCache[account]
}

/*
 * Descrp: 从缓存中删除玩家
 * Create: zhangyi 2020-06-20 01:34:24
 */
func (pc *playerCache) RemovePlayer(pid player.PID) {
	pc.Lock()
	defer pc.Unlock()
	ply := pc.GetPlayerByPid(pid)
	if ply == nil {
		return
	}
	pc.pidCache[ply.Pid] = nil
	pc.accountCache[ply.Account] = nil
}
