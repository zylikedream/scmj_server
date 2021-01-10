package shuffle

import (
	"math/rand"
	"sort"
	"time"
	"zinx-mj/game/rule/irule"
)

// 随机洗牌
type sortShuffle struct {
}

func NewSortShuffle() irule.IShuffle {
	rand.Seed(time.Now().UnixNano())
	return &sortShuffle{}
}

func (r *sortShuffle) Shuffle(cards []int) []int {
	sort.Ints(cards)
	return cards
}
