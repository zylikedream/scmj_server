package tablemgr

import "zinx-mj/game/table/itable"

var tables map[int]itable.IMjTable

func init() {
	tables = make(map[int]itable.IMjTable)
}

func CreateRoom() itable.IMjTable {
	return nil
}
