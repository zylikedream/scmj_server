package handcard

import (
	"fmt"
	"zinx-mj/game/gamedefine"
)

// 玩家手牌
type HandCard struct {
	HandCardMap  map[int]int      // 手牌map
	TingCard     map[int]struct{} // 以听的牌
	DiscardCards []int            // 玩家已打出的手牌
	DrawCards    []int            // 玩家已摸到的手牌
	KongCards    map[int]struct{} // 玩家杠的牌
	PongCards    map[int]struct{} // 玩家碰的牌
	CardCount    int              // 玩家手牌数量
	MaxCardCount int              // 玩家最大手牌数量
}

func New(handCards []int, maxCardCount int) *HandCard {
	playerCard := &HandCard{
		MaxCardCount: maxCardCount,
	}
	playerCard.HandCardMap = make(map[int]int, maxCardCount)
	for _, card := range handCards {
		playerCard.HandCardMap[card] += 1
	}
	playerCard.KongCards = make(map[int]struct{})
	playerCard.PongCards = make(map[int]struct{})

	return playerCard
}

/*
 * Descrp: 得到某张牌的数量
 * Create: zhangyi 2020-07-03 14:43:07
 */
func (p *HandCard) GetCardNum(c int) int {
	return p.HandCardMap[c]
}

/*
 * Descrp: 出某一张牌
 * Create: zhangyi 2020-07-03 14:42:46
 */
func (p *HandCard) Discard(c int) error {
	if err := p.DecCard(c, 1); err != nil {
		return err
	}
	p.DiscardCards = append(p.DiscardCards, c)
	return nil
}

/*
 * Descrp: 减少手牌
 * Create: zhangyi 2020-07-03 16:56:18
 */
func (p *HandCard) DecCard(c int, num int) error {
	if p.GetCardNum(c) < num {
		return fmt.Errorf("dec failed, card not enough, card=%d, num=%d, dec=%d",
			c, p.GetCardNum(c), num)
	}
	p.HandCardMap[c] -= num
	p.CardCount -= num
	return nil
}

/*
 * Descrp: 得到上次摸的牌
 * Create: zhangyi 2020-07-03 14:43:23
 */
func (p *HandCard) GetLastDraw() int {
	if len(p.DrawCards) == 0 {
		return 0
	}
	return p.DrawCards[len(p.DrawCards)-1]
}

// 得到上次打的牌
func (p *HandCard) GetLastDiscard() int {
	if len(p.DiscardCards) == 0 {
		return 0
	}
	return p.DiscardCards[len(p.DiscardCards)-1]
}

/*
 * Descrp:  摸一张牌
 * Create: zhangyi 2020-07-03 15:02:36
 */
func (p *HandCard) Draw(c int) error {
	if p.CardCount >= p.MaxCardCount {
		return fmt.Errorf("card too much, cardCount=%d, maxCardCount=%d", p.CardCount, p.MaxCardCount)
	}
	p.HandCardMap[c] += 1
	p.CardCount++
	return nil
}

/*
 * Descrp: 得到某种花色的牌
 * Param: cardSuit 花色
 * Create: zhangyi 2020-07-03 16:06:34
 */
func (p *HandCard) GetCardBySuit(cardSuit int) []int {
	var cards []int
	for c, num := range p.HandCardMap {
		if gamedefine.GetCardSuit(c) != cardSuit {
			continue
		}
		for i := 0; i < num; i++ {
			cards = append(cards, c)
		}
	}
	return cards
}

/*
 * Descrp: 刮风
 * Create: zhangyi 2020-07-03 16:34:57
 */
func (p *HandCard) Kong(c int) error {
	if err := p.DecCard(c, 3); err != nil {
		return err
	}
	p.KongCards[c] = struct{}{}
	return nil
}

/*
 * Descrp: 明杠牌（碰了以后杠)
 * Create: zhangyi 2020-07-03 16:49:31
 */
func (p *HandCard) ExposedKong(c int) error {
	if _, ok := p.PongCards[c]; !ok {
		return fmt.Errorf("can't exposed kong, card not pong, card=%d", c)
	}
	if err := p.DecCard(c, 1); err != nil {
		return err
	}
	p.KongCards[c] = struct{}{}
	delete(p.PongCards, c)
	return nil
}

func (p *HandCard) IsPonged(c int) bool {
	if _, ok := p.PongCards[c]; ok {
		return true
	}
	return false
}

/*
 * Descrp: 暗杠牌
 * Create: zhangyi 2020-07-03 17:03:26
 */
func (p *HandCard) ConcealedKong(c int) error {
	if err := p.DecCard(c, 4); err != nil {
		return err
	}
	p.KongCards[c] = struct{}{}
	return nil
}

/*
 * Descrp: 碰牌
 * Create: zhangyi 2020-07-03 17:10:07
 */
func (p *HandCard) Pong(c int) error {
	if err := p.DecCard(c, 2); err != nil {
		return err
	}
	p.PongCards[c] = struct{}{}
	return nil
}

func (p *HandCard) IsTingCard(c int) bool {
	_, ok := p.TingCard[c]
	return ok
}

func (p *HandCard) GetCardTotalCount() int {
	return p.CardCount
}

func (p *HandCard) GetHandCard() []int {
	var cards []int
	for card, count := range p.HandCardMap {
		for i := 0; i < count; i++ {
			cards = append(cards, card)
		}
	}
	return cards
}

// 推荐打的牌
func (p *HandCard) GetRecommandCard() int {
	for c := range p.HandCardMap {
		return c
	}
	return gamedefine.CARD_MAX
}

// 生成哨兵手牌
func (p *HandCard) GetGuardHandCard() []int {
	var cards []int
	for _, count := range p.HandCardMap {
		for i := 0; i < count; i++ {
			cards = append(cards, -1)
		}
	}
	return cards
}
