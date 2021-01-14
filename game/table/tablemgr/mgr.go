package tablemgr

import (
	"errors"
	"fmt"
	"sync"
	"time"
	"zinx-mj/game/gamedefine"
	"zinx-mj/game/table/gametable/sccardtable"
	"zinx-mj/game/table/tableoperate"
	"zinx-mj/game/table/tableplayer"
	"zinx-mj/network/protocol"
	"zinx-mj/player"

	"github.com/Pallinder/go-randomdata"
	"google.golang.org/protobuf/proto"

	"github.com/aceld/zinx/zlog"
)

type ITable interface {
	GetID() uint32
	GetPlayerByPid(pid player.PID) *tableplayer.TablePlayer
	GetPlayerBySeat(seat int) *tableplayer.TablePlayer
	// 加入间坐姿
	PlayerJoin(plyData *tableplayer.TablePlayerData, identity uint32) (*tableplayer.TablePlayer, error)
	// 退出桌子
	Quit(pid player.PID) error

	// 桌子编号
	GetTableNumber() uint32

	Update(delta time.Duration)

	OnPlyOperate(pid uint64, data tableoperate.OperateCommand) error
	SetReady(pid uint64, ready bool)
	IsEnd() bool
}

var tables map[uint32]ITable
var tableLock sync.RWMutex
var IDPool []uint32
var IDLock sync.Mutex

const IDPoolSize = 10000

func init() {
	tables = make(map[uint32]ITable)
	IDPool = make([]uint32, IDPoolSize)
	startID := randomdata.Number(111111, 678910)
	for i := 0; i < len(IDPool); i++ {
		IDPool[i] = uint32(startID + i)
	}
}

var ErrPoolIDEmpty = errors.New("")
var ErrPoolIDFull = errors.New("")
var ErrCreateTableFailed = errors.New("")

func popID() (uint32, error) {
	IDLock.Lock()
	defer IDLock.Unlock()
	poolSize := len(IDPool)
	if poolSize == 0 {
		return 0, fmt.Errorf("no id valid%w", ErrPoolIDEmpty)
	}
	id := IDPool[poolSize-1]
	IDPool = IDPool[:poolSize-1]
	return id, nil
}

func pushID(id uint32) error {
	if len(IDPool) == IDPoolSize {
		return fmt.Errorf("id pool is full%w", ErrPoolIDFull)
	}
	IDPool = append(IDPool, id)
	return nil
}

func CreateTable(master *tableplayer.TablePlayerData, tableType int, message proto.Message) (ITable, error) {
	var mjtable ITable
	var err error
	switch tableType {
	case gamedefine.TABLE_TYPE_SCMJ:
		mjtable, err = createScmjTable(master, message)
	default:
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	tableLock.Lock()
	tables[mjtable.GetID()] = mjtable
	tableLock.Unlock()
	return mjtable, nil
}

func RemoveTable(id uint32) {
	tableLock.Lock()
	defer tableLock.Unlock()
	delete(tables, id)
}

func GetTable(id uint32) ITable {
	tableLock.RLock()
	defer tableLock.RUnlock()
	return tables[id]
}

func createScmjTable(master *tableplayer.TablePlayerData, req proto.Message) (*sccardtable.ScCardTable, error) {
	msg, ok := req.(*protocol.CsCreateScmjTable)
	if !ok {
		zlog.Errorf("wrong message type %T", req)
		return nil, fmt.Errorf("wrong message type %T%w", req, ErrCreateTableFailed)
	}

	tableData := &sccardtable.ScTableData{}
	if err := tableData.UnpackFromPBMsg(msg.GetData()); err != nil {
		return nil, fmt.Errorf("unpcak pb msg failed: %s", err)
	}
	tableID, err := popID()
	if err != nil {
		return nil, fmt.Errorf("get table id failed%w", ErrCreateTableFailed)
	}
	table, err := sccardtable.NewTable(tableID, master, tableData)
	if err != nil {
		return nil, fmt.Errorf("new scmj table failed")
	}
	return table, nil
}

func Update(delta time.Duration) {
	tableLock.RLock()
	defer tableLock.RUnlock()
	for _, table := range tables {
		if table.IsEnd() { // 桌子结束后，删除改桌子
			RemoveTable(table.GetID())
			return
		}
		table.Update(delta)
	}
}
