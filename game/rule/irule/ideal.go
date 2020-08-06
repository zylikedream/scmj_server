package irule

type IDeal interface {
	Deal(cards []int, count int) []int
}
