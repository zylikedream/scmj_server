package util

import (
	"github.com/aceld/zinx/ziface"
	"sort"
	"zinx-mj/game/table/tableplayer"
	"zinx-mj/player"
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

func PackTablePlayerDataFromPly(ply *player.Player) *tableplayer.TablePlayerData {
	return &tableplayer.TablePlayerData{
		Pid: ply.Pid,
	}
}
