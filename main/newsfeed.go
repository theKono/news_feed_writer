package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"math/rand"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"

	"github.com/theKono/news_feed_writer/config"
)

const (
	MillisecondBits = 42
	ShardBits       = 12
	RandomBits      = 10
	MaxRandom       = 1024 // 2**RandomBits
)

var (
	random = rand.New(rand.NewSource(time.Now().UnixNano()))
	shards = []*gorm.DB{
		connectDatabase(config.GetShardDbUri(0)),
		connectDatabase(config.GetShardDbUri(1)),
	}
)

type Newsfeed struct {
	NewsfeedID     int64 `gorm:"primary_key"`
	Kid            int32
	ObservableType string `gorm:"size:32"`
	ObservableID   string `gorm:"size:36"`
	EventType      string
	TargetType     string
	TargetID       string
	Summary        string `gorm:"size:16383"`
	State          int8
	CreatedAt      time.Time
}

type NewsfeedJson struct {
	NewsfeedID     int64  `json:"newsfeed_id"`
	Kid            int32  `json:"kid"`
	ObservableType string `json:"observable_type"`
	ObservableID   string `json:"observable_id"`
	EventType      string `json:"event_type"`
	TargetType     string `json:"target_type"`
	TargetID       string `json:"target_id"`
	Summary        string `json:"summary"`
}

func (nj *NewsfeedJson) NewNewsfeed() *Newsfeed {
	return &Newsfeed{
		NewsfeedID:     nj.NewsfeedID,
		Kid:            nj.Kid,
		ObservableType: nj.ObservableType,
		ObservableID:   nj.ObservableID,
		EventType:      nj.EventType,
		TargetType:     nj.TargetType,
		TargetID:       nj.TargetID,
		Summary:        nj.Summary,
	}
}

func (n *Newsfeed) BeforeSave() error {
	return n.Validate()
}

func (n *Newsfeed) GetShard() int {
	return int(n.Kid % 2)
}

func (n *Newsfeed) Validate() error {
	if n.Kid == 0 {
		return errors.New("kid is required")
	}

	if n.NewsfeedID == 0 {
		n.NewsfeedID = GenerateId(n.GetShard())
	}

	var summary map[string]interface{}
	if err := json.Unmarshal([]byte(n.Summary), &summary); err != nil {
		return errors.New("summary is a bad json")
	}

	if _, ok := summary["id"]; !ok {
		summary["id"] = fmt.Sprint(n.NewsfeedID)
		j, err := json.Marshal(summary)

		if err != nil {
			return errors.New("cannot marshal to json string")
		}

		n.Summary = string(j)
	}

	return nil
}

func (n *Newsfeed) Save() error {
	if err := shards[n.GetShard()].Create(n).Error; err != nil {
		log.Println("Error", "Cannot insert into database")
		log.Print(err)
		return err
	}

	return nil
}

func GenerateId(shard int) int64 {
	milliSecond := time.Now().UnixNano() / 1.0e6
	ret := milliSecond << (ShardBits + RandomBits)
	ret += int64(shard << RandomBits)
	ret += int64(random.Intn(MaxRandom))
	return ret
}

func connectDatabase(dataSourceName string) *gorm.DB {
	db, err := gorm.Open("mysql", dataSourceName)

	if err != nil {
		log.Fatalf("Cannot connect to %q", dataSourceName)
		log.Fatal(err.Error())
	}

	if config.IsDebug() {
		db.LogMode(true)
	}

	if !db.HasTable("newsfeeds") {
		log.Fatalf(`Table "newsfeeds" does not exist`)
	}

	return db
}
