package gamedefine

const (
	CARD_SUIT_EMPTY     = iota // 空类型
	CARD_SUIT_CHARACTER        // 万
	CARD_SUIT_BAMBOO           // 条
	CARD_SUIT_DOT              // 筒
	CARD_SUIT_MAX
)

const CARD_BASE = 10
const CARD_MAX = CARD_SUIT_MAX * CARD_BASE
const CARD_EMPTY = 0

/*
 * Descrp: 得到牌的花色
 * Create: zhangyi 2020-07-02 18:18:01
 */
func GetCardSuit(cardNum int) int {
	return cardNum / CARD_BASE
}

/*
 * Descrp: 得到牌的数字1-9
 * Create: zhangyi 2020-07-02 18:18:09
 */
func GetCardRank(cardNum int) int {
	return cardNum % CARD_BASE
}

func GetCardNumber(suit int, cardIndex int) int {
	return suit*CARD_BASE + cardIndex
}

/*
 * Descrp: 得到一个牌的相邻牌
 * Create: zhangyi 2020-08-03 22:51:07
 */
func GetNeighborCards(card int) []int {
	neighbors := []int{card}
	if IsValidCard(card + 1) {
		neighbors = append(neighbors, card+1)
	}
	if IsValidCard(card - 1) {
		neighbors = append(neighbors, card-1)
	}
	return neighbors
}

func IsValidCard(card int) bool {
	if card >= CARD_SUIT_MAX*CARD_BASE || card <= CARD_SUIT_EMPTY {
		return false
	}
	if card%CARD_BASE == 0 {
		return false
	}
	return true
}
