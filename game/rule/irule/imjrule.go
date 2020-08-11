package irule

import "zinx-mj/game/room/player"

type IMjRule interface {
	Chow(pid int, c int) error
	Discard(pid int, c int) error
	Draw(pid int, c int) error
	Kong(pid int, c int) error
	Pong(pid int, c int) error
	Win(pid int, c int) error

	GetCurPlayer() *player.RoomPlayer
	IsPlayerTurn(pid int) bool
}
