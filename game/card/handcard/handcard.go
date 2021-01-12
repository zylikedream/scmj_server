package handcard

import (
	"fmt"
	"zinx-mj/game/gamedefine"
	"zinx-mj/util"

	"github.com/pkg/errors"
)

type KongInfo struct {
	Card  int
	KType int
}

const (
	KONG_TYPE_WIND      = iota // 刮风
	KONG_TYPE_CONCEALED        // 暗杠
	KONG_TYPE_EXPOSED          // 明杠
)

// 玩家手牌
type HandCard struct {
	CardMap      map[int]int      // 手牌map
	cardCount    int              // 手牌数组
	TingCard     map[int]struct{} // 以听的牌
	DiscardCards []int            // 玩家已打出的手牌
	DrawCards    []int            // 玩家已摸到的手牌
	KongCards    []KongInfo       // 玩家杠的牌
	PongCards    []int            // 玩家碰的牌
	MaxCardCount int              // 玩家最大手牌数量
}

func New(maxCardCount int) *HandCard {
	p := &HandCard{
		MaxCardCount: maxCardCount,
	}

	p.CardMap = make(map[int]int)
	return p
}

func (p *HandCard) SetHandCard(cards []int) error {
	if len(cards) > p.MaxCardCount {
		return errors.Errorf("more than max, cnt:%d, max:%d", len(cards), p.MaxCardCount)
	}
	for _, card := range cards {
		p.CardMap[card] += 1
	}
	p.cardCount = len(cards)
	return nil
}

/*
 * Descrp: 得到某张牌的数量
 * Create: zhangyi 2020-07-03 14:43:07
 */
func (p *HandCard) GetCardNum(c int) int {
	return p.CardMap[c]
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
		return errors.Errorf("dec failed, card not enough, card=%d, num=%d, dec=%d",
			c, p.GetCardNum(c), num)
	}
	p.CardMap[c] -= num
	if p.CardMap[c] == 0 {
		delete(p.CardMap, c)
	}
	p.cardCount -= num
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
	if p.cardCount >= p.MaxCardCount {
		return fmt.Errorf("card too much, cardCount=%d, maxCardCount=%d", p.cardCount, p.MaxCardCount)
	}
	p.CardMap[c] += 1
	p.cardCount++
	p.DrawCards = append(p.DrawCards, c)
	return nil
}

/*
 * Descrp: 得到某种花色的牌
 * Param: cardSuit 花色
 * Create: zhangyi 2020-07-03 16:06:34
 */
func (p *HandCard) GetCardBySuit(cardSuit int) []int {
	var cards []int
	for c, num := range p.CardMap {
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
	p.KongCards = append(p.KongCards, KongInfo{
		Card:  c,
		KType: KONG_TYPE_WIND,
	})
	return nil
}

func (p *HandCard) IsKong(c int) bool {
	for _, k := range p.KongCards {
		if k.Card == c {
			return true
		}
	}
	return false
}

/*
 * Descrp: 明杠牌（碰了以后杠)
 * Create: zhangyi 2020-07-03 16:49:31
 */
func (p *HandCard) ExposedKong(c int) error {
	if !p.IsPonged(c) {
		return fmt.Errorf("can't exposed kong, card not pong, card=%d", c)
	}
	if err := p.DecCard(c, 1); err != nil {
		return err
	}
	p.KongCards = append(p.KongCards, KongInfo{
		Card:  c,
		KType: KONG_TYPE_EXPOSED,
	})
	p.RemovePong(c)
	return nil
}

func (p *HandCard) IsPonged(c int) bool {
	for _, p := range p.PongCards {
		if c == p {
			return true
		}
	}
	return false
}

func (p *HandCard) RemovePong(c int) []int {
	var newPong []int
	for _, p := range p.PongCards {
		if p == c {
			continue
		}
		newPong = append(newPong, p)
	}
	p.PongCards = newPong
	return newPong
}

/*
 * Descrp: 暗杠牌
 * Create: zhangyi 2020-07-03 17:03:26
 */
func (p *HandCard) ConcealedKong(c int) error {
	if err := p.DecCard(c, 4); err != nil {
		return err
	}
	p.KongCards = append(p.KongCards, KongInfo{
		Card:  c,
		KType: KONG_TYPE_CONCEALED,
	})
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
	p.PongCards = append(p.PongCards, c)
	return nil
}

func (p *HandCard) IsTingCard(c int) bool {
	_, ok := p.TingCard[c]
	return ok
}

func (p *HandCard) GetCardTotalCount() int {
	return p.cardCount
}

func (p *HandCard) GetHandCard() []int {
	return util.IntMapToIntSlice(p.CardMap)
}

// 推荐打的牌
func (p *HandCard) GetRecommandCard() int {
	for c := range p.CardMap {
		return c
	}
	return gamedefine.CARD_MAX
}

// 生成哨兵手牌
func (p *HandCard) GetGuardHandCard() []int {
	var cards []int
	for _, count := range p.CardMap {
		for i := 0; i < count; i++ {
			cards = append(cards, -1)
		}
	}
	return cards
}
