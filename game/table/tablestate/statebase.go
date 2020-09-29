package tablestate

type state struct {
	onEnter func(...interface{}) error
	onExit  func(...interface{}) error
	name    string
}
