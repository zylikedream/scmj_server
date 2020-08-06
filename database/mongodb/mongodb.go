/*
 * Copyright (c) zhangyi All rights reserved. 2020
 * Descrp: 封装mongo数据库操作
 * Create: zhangyi 2020-6-20 0:10:27
 */
package mongodb

import (
	"context"
	"errors"
	"fmt"
	"github.com/aceld/zinx/zlog"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/writeconcern"
	"time"
	"zinx-mj/player"
)

const (
	DB_NAME           = "scmj"
	COLLECTION_PLAYER = "player"
)

type mongoDB struct {
	client          *mongo.Client
	addr            string
	port            int
	options         *options.ClientOptions
	databaseCache   map[string]*mongo.Database
	collectionCache map[string]*mongo.Collection
}

/*
 * Descrp: 初始化数据库实例
 * Notice: 只允许初始化一个实例
 * Create: zhangyi 2020-6-20 0:20:32
 */
func New() *mongoDB {
	db := &mongoDB{
		options: options.Client(),
	}
	db.options.SetMaxPoolSize(10)
	db.options.SetMinPoolSize(2)
	db.options.SetConnectTimeout(10 * time.Second)         // 连接超时时间
	db.options.SetServerSelectionTimeout(10 * time.Second) // ping 超时时间
	db.options.SetWriteConcern(writeconcern.New(           // write concern
		writeconcern.W(1), writeconcern.J(false), writeconcern.WTimeout(5*time.Second)))
	db.initCache() // 初始化缓存
	return db
}

/*
 * Descrp: 初始化cache
 * Create: zhangyi 2020-6-20 0:18:42
 */
func (m *mongoDB) initCache() {
	m.databaseCache = make(map[string]*mongo.Database)
	m.collectionCache = make(map[string]*mongo.Collection)
}

/*
 * Descrp: 连接数据库
 * Param: dbAddr-数据库地址, port-数据库端口
 * Create: zhangyi 2020-6-20 0:13:9
 */
func (m *mongoDB) Connect(dbAddr string, port int) error {
	m.options.ApplyURI(fmt.Sprintf("mongodb://%s:%d", dbAddr, port))
	client, err := mongo.Connect(context.TODO(), m.options)
	if err != nil {
		zlog.Errorf("connect mongo error, err=%s\n", dbAddr, port, err)
		return err
	}
	m.client = client
	m.addr = dbAddr
	m.port = port
	return nil
}

/*
 * Descrp: 检测数据库是否连通
 * Create: zhangyi 2020-6-20 0:14:35
 */
func (m *mongoDB) Ping() error {
	if m.client == nil {
		return nil
	}
	if err := m.client.Ping(context.TODO(), nil); err != nil {
		zlog.Errorf("ping mongo error, addr=%s, port=%s, err=%s\n", m.addr, m.port, err)
		return err
	}
	return nil
}

/*
 * Descrp: 断开数据库连接
 * Create: zhangyi 2020-6-20 0:16:4
 */
func (m *mongoDB) Disconnect() {
	if m.client == nil {
		return
	}
	if err := m.client.Disconnect(context.TODO()); err != nil {
		zlog.Errorf("disconnect mongodb failed")
	}
	m.client = nil
	m.initCache() // 清空cache
}

/*
 * Descrp: 得到mongo数据库实例
 * Param: dbName-数据库名字
 * Create: zhangyi 2020-6-20 0:15:31
 */
func (m *mongoDB) getDataBase(dbName string) (*mongo.Database, error) {
	if m.client == nil {
		return nil, errors.New("mongo not connect")
	}
	//if m.databaseCache[dbName] == nil {
	//	m.databaseCache[dbName] = m.client.Database(DB_NAME)
	//}
	//return m.databaseCache[dbName], nil
	return m.client.Database(DB_NAME), nil
}

/*
 * Descrp: 得到mongo集合
 * Create: zhangyi 2020-6-20 0:21:1
 */
func (m *mongoDB) getCollection(dbName string, collectName string) (*mongo.Collection, error) {
	//if m.collectionCache[collectName] == nil {
	//	db, err := m.getDataBase(dbName)
	//	if err != nil {
	//		return nil, err
	//	}
	//	m.collectionCache[collectName] = db.Collection(collectName)
	//}
	//return m.collectionCache[collectName], nil
	db, err := m.getDataBase(dbName)
	if err != nil {
		return nil, err
	}
	return db.Collection(collectName), nil
}

/*
 * Descrp: 从数据库加载玩家数据
 * Param: account-账号名字 ply-保存玩家数据指针(出参)
 * Create: zhangyi 2020-6-20 0:21:6
 */
func (m *mongoDB) LoadPlayer(account string) (*player.Player, error) {
	filter := bson.M{"account": account}
	collection, err := m.getCollection(DB_NAME, COLLECTION_PLAYER)
	if err != nil {
		zlog.Errorf("load player failed, account=%s, err=%s", account, err)
		return nil, err
	}
	ply := player.New()
	if err = collection.FindOne(context.TODO(), filter).Decode(ply); err != nil {
		zlog.Errorf("load player failed, account=%s, err=%s", account, err)
		return nil, err
	}
	return ply, nil
}

/*
 * Descrp: 保存玩家数据
 * Param: ply-玩家数据
 * Create: zhangyi 2020-6-20 0:41:23
 */
func (m *mongoDB) SavePlayer(ply *player.Player) error {
	collection, err := m.getCollection(DB_NAME, COLLECTION_PLAYER)
	if err != nil {
		zlog.Errorf("save player failed, pid=%d, err=%s", ply.Pid, err)
		return err
	}
	if _, err = collection.InsertOne(context.TODO(), ply); err != nil {
		zlog.Errorf("save player failed, pid=%d, err=%s", ply.Pid, err)
		return err
	}
	return nil
}
