package mongodb

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"sync"
	"testing"
	"time"
	"zinx-mj/database"
	"zinx-mj/player"
	"zinx-mj/player/playermgr"
)

func Test_mongo(t *testing.T) {
	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")
	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		log.Fatal(err)
	}

	err = client.Ping(context.TODO(), nil)

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Connected to mongo")
	collect := client.Database("test").Collection("player")
	ply := player.New()
	filter := bson.M{
		"Name": "zhangyi",
	}
	err = collect.FindOne(context.TODO(), filter).Decode(ply)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Found a doc: %+v\n", *ply)
}

func TestMongoDB_LoadPlayer(t *testing.T) {
	database.Init()
	if _, err := database.GetDB().LoadPlayer("11111111"); err != nil {
		fmt.Printf("load failed, err=%s\n", err)
	}
	fmt.Printf("The End\n")
}

func TestMongoDB_SavePlayer(t *testing.T) {
	db := database.GetDB()
	account := "zhangyi"
	ply, err := playermgr.CreatePlayer(account)
	dbAddr := "localhost"
	dbPort := 27017
	if err = db.Connect(dbAddr, dbPort); err != nil {
		fmt.Printf("connect mongo failed, addr=%s, port=%d, err=%s\n", dbAddr, dbPort, err)
		return
	}
	if err := db.SavePlayer(ply); err != nil {
		fmt.Printf("save ply failed, pid=%d, err=%s\n", ply.Pid, err)
		return
	}

	if ply, err := db.LoadPlayer(account); err != nil {
		fmt.Printf("load ply failed, account=%s, err=%s\n", account, err)
		return
	} else {
		fmt.Printf("ply: %+v\n", ply)
	}
}

func TestMongoDB_sync(t *testing.T) {
	db := database.GetDB()
	dbAddr := "localhost"
	dbPort := 27017
	if err := db.Connect(dbAddr, dbPort); err != nil {
		fmt.Printf("connect mongo failed, addr=%s, port=%d, err=%s\n", dbAddr, dbPort, err)
		return
	}

	tm := time.Now()
	wg := sync.WaitGroup{}
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() {
			account := "zhangyi"
			ply := player.New()
			fmt.Printf("ply: %+v\n", ply)
			loadstar := time.Now()
			if _, err := db.LoadPlayer(account); err != nil {
				fmt.Printf("load ply failed, account=%s, err=%s\n", account, err)
				return
			}
			loadcost := time.Since(loadstar)
			fmt.Printf("load cost %d ms\n", loadcost.Milliseconds())
			fmt.Printf("i ply: %+v\n", ply)
			wg.Done()
		}()
	}
	wg.Wait()
	cost := time.Since(tm).Milliseconds()
	fmt.Printf("total cost %d ms", cost)
}
