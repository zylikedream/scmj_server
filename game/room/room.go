package room

const (
	ROOM_IDENTIY_MASTER  = 1 << iota // 房主
	ROOM_IDENTIY_PLAYER              // 玩家
	ROOM_IDENTIY_WATCHER             // 观看者
)

// 麻将房间类型
const (
	ROOM_MJ_SCMJ = iota // 四川麻将
)
