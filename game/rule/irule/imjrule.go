package irule

import "zinx-mj/game/table/player"

type IMjRule interface {
	Chow(pid int, c int) error
	Discard(pid int, c int) error
	Draw(pid int, c int) error
	Kong(pid int, c int) error
	Pong(pid int, c int) error
	Win(pid int, c int) error

	GetCurPlayer() *player.TablePlayer
	IsPlayerTurn(pid int) bool
}
