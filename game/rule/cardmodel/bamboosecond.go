package cardmodel

import (
	"zinx-mj/game/gamedefine"
	"zinx-mj/game/rule/irule"
	"zinx-mj/util"
)

type BambooSecond struct {
}

func (s *BambooSecond) IsModel(data *irule.CardModel) bool { // 是否是某种牌型
	bamboo2 := gamedefine.GetCardNumber(gamedefine.CARD_SUIT_BAMBOO, 2)
	if data.WinCard != bamboo2 {
		return false
	}
	// 卡二条必须要有1条和3条
	bamboo1 := gamedefine.GetCardNumber(gamedefine.CARD_SUIT_BAMBOO, 1)
	bamboo3 := gamedefine.GetCardNumber(gamedefine.CARD_SUIT_BAMBOO, 3)
	if data.HandCard[bamboo1] == 0 || data.HandCard[bamboo3] == 0 {
		return false
	}
	// 去掉1，2，3条是否还能胡
	cardCopy := util.CopyIntMap(data.HandCard)
	cardCopy[bamboo1] -= 1
	cardCopy[bamboo2] -= 1
	cardCopy[bamboo3] -= 1

	return data.WiRule.CanWin(util.IntMapToIntSlice(cardCopy))
}

func (s *BambooSecond) GetModelType() int {
	return CARD_MODEL_BAMBOO_SECOND
}

func NewBambooSecond() *BambooSecond {
	return &BambooSecond{}
}
