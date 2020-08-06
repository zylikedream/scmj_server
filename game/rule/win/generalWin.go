package win

import (
	"zinx-mj/game/card"
	"zinx-mj/game/rule/irule"
)

// 通用麻将胡牌规则
type generalWin struct {
}

func NewGeneralWin() irule.IWin {
	return &generalWin{}
}

func (r *generalWin) CanWin(cards []int) bool {
	cardLen := len(cards)
	if cardLen == 0 {
		return false
	}
	cardMap := [card.CARD_MAX]int{} // 使用数组来作为Map, 原生的Map太慢了
	var pairPos []int               // 对子的位置
	// 构造map并得到对子的位置
	for i, c := range cards {
		cardMap[c] += 1
		if i == cardLen-1 || c != cards[i+1] {
			if cardMap[c] >= 2 {
				pairPos = append(pairPos, c)
			}
		}
	}
	// 小七对
	if len(pairPos)*2 == len(cards) {
		return true
	}
	// 遍历所有可能的情况
	for _, c := range pairPos {
		if isAllMeld(cards, &cardMap, c) {
			return true
		}
	}
	return false
}

/*
 * Descrp: 是否都是一组牌(一个顺子或者刻子）
 * Create: zhangyi 2020-08-05 23:11:32
 */
func isAllMeld(sortCards []int, cardMap *[card.CARD_MAX]int, pair int) bool {
	cardUsed := [card.CARD_MAX]int{}
	cardUsed[pair] = 2 // 用了2个作为了麻将
	for _, c := range sortCards {
		cardLeft := cardMap[c] - cardUsed[c]
		switch {
		case cardLeft >= 3: // 3个或3个以上只能作为刻子
			cardUsed[c] += 3
		case cardLeft == 0: //没有牌跳过
		default: // 只能做为顺子, 那么其他的牌数量必须大于等于它
			if cardMap[c+1]-cardUsed[c+1] < cardLeft ||
				cardMap[c+2]-cardUsed[c+2] < cardLeft {
				return false
			}
			cardUsed[c] += cardLeft
			cardUsed[c+1] += cardLeft
			cardUsed[c+2] += cardLeft
		}
	}
	return cardUsed == *cardMap
}

//func (g *generalWin) CanWin(cards []int) bool {
//	// 拷贝一个新的并排序
//	sortCards := make([]int, len(cards))
//	copy(sortCards, cards)
//	//sort.Ints(sortCards)
//
//	var pairPos []int // 对子的位置
//	// 得到对子的位置
//	var start int
//	for i := 0; i < len(sortCards); i++ {
//		if i == len(sortCards)-1 || sortCards[i] != sortCards[i+1] {
//			if i-start > 0 {
//				pairPos = append(pairPos, start)
//			}
//			start = i + 1
//		}
//	}
//	if len(pairPos)*2 == len(sortCards) {
//		return true
//	}
//	for _, pos := range pairPos {
//		cards := removePair(sortCards, pos)
//		if isAllMeld(cards) {
//			return true
//		}
//	}
//	return false
//}
//
//func removePair(sortedCards []int, pos int) []int {
//	return util.RemoveSlice(sortedCards, pos, 2)
//}
//
///*
// * Descrp: 是否都是一组牌(一个顺子或者刻子）
// * Create: zhangyi 2020-08-05 23:11:32
// */
//func isAllMeld(sortCards []int) bool {
//	var ok bool
//	for len(sortCards) > 0 {
//		if len(sortCards) < 3 { // 小于3个肯定不为顺子
//			return false
//		}
//		switch {
//		// 一刻或者一杠那么只能成为刻子
//		case sortCards[0] == sortCards[1] && sortCards[1] == sortCards[2]:
//			sortCards = sortCards[3:]
//			continue
//		case sortCards[0] == sortCards[1]:
//			// 一对，只能拆成顺子AABBCC，所以至少需要六个
//			if len(sortCards) < 6 {
//				return false
//			}
//			ok, sortCards = removeSequence(sortCards, sortCards[0], 2)
//			if !ok {
//				return false
//			}
//		default:
//			// 是否为ABC
//			// sortCards[0]只有1个
//			ok, sortCards = removeSequence(sortCards, sortCards[0], 1)
//			if !ok {
//				return false
//			}
//		}
//	}
//	return true
//}
//
///*
// * Descrp: 从牌组中移除多组刻子
// * Notice: count表示移除的组数, 比如2就表示移除两组ABC
// * Create: zhangyi 2020-08-05 23:25:07
// */
//func removeSequence(sortCards []int, startCard int, count int) (bool, []int) {
//	newCards := make([]int, len(sortCards)-count*3)
//	curCount := count
//	curCard := startCard
//	totalCount := count * 3
//	for i, cd := range sortCards {
//		if totalCount == 0 { // 都已经扣完了 剩余的不用扣除了, 直接保留
//			newCards = append(newCards, sortCards[i:]...)
//			return true, newCards
//		}
//		if cd == curCard {
//			if curCount > 0 {
//				curCount--
//				totalCount--
//			} else {
//				newCards = append(newCards, cd) // 当前牌扣除的数量足够了 多余的存起来
//			}
//			continue
//		}
//		if cd != curCard + 1 {  // 不是下一张牌
//			return false, newCards
//		}
//		// 更新要扣除的牌
//		curCard = cd
//		curCount = count
//		// 跟新牌的数量
//		curCount--
//		totalCount--
//	}
//	return true, newCards
//}
