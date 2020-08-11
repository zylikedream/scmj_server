package irule

import (
	"zinx-mj/game/table/tableplayer"
	"zinx-mj/player"
)

type IMjRule interface {
	Chow(pid player.PID, c int) error
	Discard(pid player.PID, c int) error
	Draw(pid player.PID, c int) error
	Kong(pid player.PID, c int) error
	Pong(pid player.PID, c int) error
	Win(pid player.PID, c int) error

	GetCurPlayer() *tableplayer.TablePlayer
	IsPlayerTurn(pid player.PID) bool
}
