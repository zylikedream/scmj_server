package table

const (
	TABLE_IDENTIY_MASTER  = 1 << iota // 房主
	TABLE_IDENTIY_PLAYER              // 玩家
	TABLE_IDENTIY_WATCHER             // 观看者
)

// 麻将房间类型
const (
	TABLE_TYPE_SCMJ = iota // 四川麻将
)
