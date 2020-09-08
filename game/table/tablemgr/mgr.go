package tablemgr

import (
	"errors"
	"fmt"
	"sync"
	"zinx-mj/game/gamedefine"
	"zinx-mj/game/rule/gamerule"
	"zinx-mj/game/table/gametable"
	"zinx-mj/game/table/itable"
	"zinx-mj/game/table/tableplayer"
	"zinx-mj/network/protocol"

	"github.com/Pallinder/go-randomdata"

	"github.com/aceld/zinx/zlog"
	"github.com/golang/protobuf/proto"
)

var tables map[uint32]itable.ITable
var tableLock sync.RWMutex
var poolID []uint32
var poolLock sync.Mutex

const poolCap = 10000

func init() {
	tables = make(map[uint32]itable.ITable)
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

func CreateTable(master *tableplayer.TablePlayerData, tableType int, message proto.Message) (itable.ITable, error) {
	var mjtable itable.ITable
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

func GetTable(id uint32) itable.ITable {
	tableLock.RLock()
	defer tableLock.RUnlock()
	return tables[id]
}

func createScmjTable(master *tableplayer.TablePlayerData, req proto.Message) (itable.ITable, error) {
	msg, ok := req.(*protocol.CsCreateScmjTable)
	if !ok {
		zlog.Errorf("wrong message type %T", req)
		return nil, fmt.Errorf("wrong message type %T%w", req, ErrCreateTableFailed)
	}

	ruleData := &gamerule.ScmjRuleData{}
	ruleData.UnpackFromPBMsg(msg.GetRule())
	tableID, err := poolPop()
	if err != nil {
		return nil, fmt.Errorf("get table id failed%w", ErrCreateTableFailed)
	}
	table, err := gametable.NewTable(tableID, master)
	if err != nil {
		return nil, fmt.Errorf("new scmj table failed")
	}
	gameRule := gamerule.NewScmjRule(ruleData, table)
	table.SetRule(gameRule)
	return table, nil
}
