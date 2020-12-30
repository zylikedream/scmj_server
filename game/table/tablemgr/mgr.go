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
}

var tables map[uint32]ITable
var tableLock sync.RWMutex
var poolID []uint32
var poolLock sync.Mutex

const poolCap = 10000

func init() {
	tables = make(map[uint32]ITable)
	poolID = make([]uint32, poolCap)
	startID := randomdata.Number(111111, 678910)
	for i := 0; i < len(poolID); i++ {
		poolID[i] = uint32(startID + i)
	}
}

var ErrPoolIDEmpty = errors.New("")
var ErrPoolIDFull = errors.New("")
var ErrCreateTableFailed = errors.New("")

func poolPop() (uint32, error) {
	poolLock.Lock()
	defer poolLock.Unlock()
	poolSize := len(poolID)
	if poolSize == 0 {
		return 0, fmt.Errorf("no id valid%w", ErrPoolIDEmpty)
	}
	id := poolID[poolSize-1]
	poolID = poolID[:poolSize-1]
	return id, nil
}

func poolPush(id uint32) error {
	if len(poolID) == poolCap {
		return fmt.Errorf("id pool is full%w", ErrPoolIDFull)
	}
	poolID = append(poolID, id)
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
	tableID, err := poolPop()
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
		table.Update(delta)
	}
}
