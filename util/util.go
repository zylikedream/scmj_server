package util

import (
	"fmt"
	"sort"
	"zinx-mj/game/table/tableplayer"
	"zinx-mj/mjerror"
	"zinx-mj/network/protocol"
	"zinx-mj/player"
	"zinx-mj/player/playermgr"

	"google.golang.org/protobuf/proto"

	"github.com/aceld/zinx/ziface"
)

func RemoveSlice(sli []int, startPos int, length int) []int {
	remainCards := make([]int, 0, len(sli)-length)
	remainCards = append(remainCards, sli[:startPos]...)
	remainCards = append(remainCards, sli[startPos+length:]...)
	return remainCards
}

func CopyIntMap(src map[int]int) map[int]int {
	copyMap := make(map[int]int, len(src))
	for k, v := range src {
		copyMap[k] = v
	}
	return copyMap
}

func IntMapToIntSlice(src map[int]int) []int {
	var slice []int
	for k, v := range src {
		for i := 0; i < v; i++ {
			slice = append(slice, k)
		}
	}
	sort.Ints(slice)
	return slice
}

/*
 * Descrp: 通过连接得到pid
 * Create: zhangyi 2020-08-12 00:00:03
 */
func GetConnPid(conn ziface.IConnection) player.PID {
	data, err := conn.GetProperty("pid")
	if err != nil {
		return 0
	}
	return data.(player.PID)
}

func GetPidConn(pid player.PID) ziface.IConnection {
	ply := playermgr.GetPlayerByPid(pid)
	if ply == nil {
		return nil
	}
	return ply.Conn
}

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

func SendMsg(pid player.PID, protoID protocol.PROTOID, msg proto.Message) error {
	data, err := proto.Marshal(msg)
	if err != nil {
		return fmt.Errorf("marshal packet failed, %w", err)
	}
	conn := GetPidConn(pid)
	if err = conn.SendMsg(uint32(protoID), data); err != nil {
		return fmt.Errorf("send msg failed, %w", err)
	}
	return nil
}
