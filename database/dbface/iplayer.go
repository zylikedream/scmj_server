package dbface

import "zinx-mj/player"

type IDBPlayer interface {
	LoadPlayer(account string) (*player.Player, error)
	SavePlayer(ply *player.Player) error
}
