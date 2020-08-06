package playercard

import (
	"fmt"
	"zinx-mj/game/card"
)

// 玩家手牌
type PlayerCard struct {
	HandCardMap  map[int]int      // 手牌map
	DiscardCards []int            // 玩家已打出的手牌
	DrawCards    []int            // 玩家已摸到的手牌
	KongCards    map[int]struct{} // 玩家杠的牌
	PongCards    map[int]struct{} // 玩家碰的牌
	CardCount    int              // 玩家手牌数量
	MaxCardCount int              // 玩家最大手牌数量
}

func NewPlayerCard(maxCardCount int) *PlayerCard {
	playerCard := &PlayerCard{
		MaxCardCount: maxCardCount,
	}
	playerCard.HandCardMap = make(map[int]int)
	playerCard.KongCards = make(map[int]struct{})
	playerCard.PongCards = make(map[int]struct{})

	return playerCard
}

/*
 * Descrp: 得到某张牌的数量
 * Create: zhangyi 2020-07-03 14:43:07
 */
func (p *PlayerCard) GetCardNum(card int) int {
	return p.HandCardMap[card]
}

/*
 * Descrp: 出某一张牌
 * Create: zhangyi 2020-07-03 14:42:46
 */
func (p *PlayerCard) Discard(card int) error {
	if err := p.DecCard(card, 1); err != nil {
		return err
	}
	p.DiscardCards = append(p.DiscardCards, card)
	return nil
}

/*
 * Descrp: 减少手牌
 * Create: zhangyi 2020-07-03 16:56:18
 */
func (p *PlayerCard) DecCard(card int, num int) error {
	if p.GetCardNum(card) < num {
		return fmt.Errorf("dec failed, card not enough, card=%d, num=%d, dec=%d",
			card, p.GetCardNum(card), num)
	}
	p.HandCardMap[card] -= num
	p.CardCount -= num
	return nil
}

/*
 * Descrp: 得到上次摸的牌
 * Create: zhangyi 2020-07-03 14:43:23
 */
func (p *PlayerCard) GetLastDraw() int {
	if len(p.DrawCards) == 0 {
		return 0
	}
	return p.DrawCards[len(p.DrawCards)-1]
}

/*
 * Descrp:  摸一张牌
 * Create: zhangyi 2020-07-03 15:02:36
 */
func (p *PlayerCard) Draw(card int) error {
	if p.CardCount >= p.MaxCardCount {
		return fmt.Errorf("card too much, cardCount=%d, maxCardCount=%d", p.CardCount, p.MaxCardCount)
	}
	p.HandCardMap[card] += 1
	p.CardCount++
	return nil
}

/*
 * Descrp: 得到某种花色的牌
 * Param: cardSuit 花色
 * Create: zhangyi 2020-07-03 16:06:34
 */
func (p *PlayerCard) GetCardBySuit(cardSuit int) []int {
	var cards []int
	for c, num := range p.HandCardMap {
		if card.GetCardSuit(c) != cardSuit {
			continue
		}
		for i := 0; i < num; i++ {
			cards = append(cards, c)
		}
	}
	return cards
}

/*
 * Descrp: 杠牌
 * Create: zhangyi 2020-07-03 16:34:57
 */
func (p *PlayerCard) Kong(card int) error {
	if err := p.DecCard(card, 3); err != nil {
		return err
	}
	p.KongCards[card] = struct{}{}
	return nil
}

/*
 * Descrp: 明杠牌（碰了以后杠)
 * Create: zhangyi 2020-07-03 16:49:31
 */
func (p *PlayerCard) ExposedKong(card int) error {
	if _, ok := p.PongCards[card]; !ok {
		return fmt.Errorf("can't exposed kong, card not pong, card=%d", card)
	}
	if err := p.DecCard(card, 1); err != nil {
		return err
	}
	p.KongCards[card] = struct{}{}
	delete(p.KongCards, card)
	return nil
}

/*
 * Descrp: 暗杠牌
 * Create: zhangyi 2020-07-03 17:03:26
 */
func (p *PlayerCard) ConcealedKong(card int) error {
	if err := p.DecCard(card, 4); err != nil {
		return err
	}
	p.KongCards[card] = struct{}{}
	return nil
}

/*
 * Descrp: 碰牌
 * Create: zhangyi 2020-07-03 17:10:07
 */
func (p *PlayerCard) Pong(card int) error {
	if err := p.DecCard(card, 2); err != nil {
		return err
	}
	p.PongCards[card] = struct{}{}
	return nil
}
