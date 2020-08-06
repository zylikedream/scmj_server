package shuffle

import (
	"math/rand"
	"time"
	"zinx-mj/game/rule/irule"
)

// 随机洗牌
type randomShuffle struct {
}

func NewRandomShuffle() irule.IShuffle {
	rand.Seed(time.Now().UnixNano())
	return &randomShuffle{}
}

func (r *randomShuffle) Shuffle(cards []int) []int {
	rand.Shuffle(len(cards), func(i int, j int) {
		cards[i], cards[j] = cards[j], cards[i]
	})
	return cards
}
