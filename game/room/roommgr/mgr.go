package roommgr

import "zinx-mj/game/room/iroom"

var rooms map[int]iroom.IMjRoom

func init() {
	rooms = make(map[int]iroom.IMjRoom)
}

func CreateRoom() iroom.IMjRoom {
	return nil
}
