package irule

// 计算牌型接口
type ScorePoint struct {
	Point int // 番数
	Base  int // 底数(自摸加底之类的)
}

func (s *ScorePoint) GetFinalPoint() int {
	return 1<<s.Point + s.Base
}

type IScorePoint interface {
	ScorePoint(models []int, winMode int, kongs int, handcard map[int]int) ScorePoint
}
