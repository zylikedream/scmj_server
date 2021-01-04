package irule

// 计算牌型接口
type ScorePoint struct {
	Point int // 番数
	Base  int // 底数(自摸加底之类的)
}

type IScorePoint interface {
	ScorePoint(models []int, winMode int, kongCard map[int]struct{}) ScorePoint
}
