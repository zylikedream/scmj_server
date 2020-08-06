package dbface

type IDBConnect interface {
	Connect(dbAddr string, port int) error
	Ping() error
	Disconnect()
}
