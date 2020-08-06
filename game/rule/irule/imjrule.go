package irule

type IMjRule interface {
	GetBoardRule() IBoard
	GetChowRule() IChow
	GetDiscardRule() IDiscard
	GetDrawRule() IDraw
	GetKongRule() IKong
	GetPongRule() IPong
	GetShuffleRule() IShuffle
	GetTingRule() ITing
	GetWinRule() IWin
	GetDealRule() IDeal
}
