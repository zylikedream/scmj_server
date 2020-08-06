package define

const (
	SCORING_DRAW_WIN      = iota // 自摸
	SCORING_KONG_WIN             // 杠上花
	SCORING_SEVEN_PAIRS          // 七对
	SCORING_ROBBING_KONG         // 抢杠
	SOCRING_KONG_DISCARD         // 杠上炮
	SCORING_HEVENLY_WIN          // 天胡
	SCORING_EARTHLY_WIN          // 地胡
	SCORING_CONCEALED_WIN        // 门清
	SCORING_MIDDLE_TILE          // 中张
	SCORING_ONE_SUIT             // 清一色
	SCORING_BIG_PAIR             // 大对子
)
