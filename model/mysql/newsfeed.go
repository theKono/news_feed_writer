package mysql

import (
	"time"

	"github.com/theKono/orchid/model/messagejson"
)

// UnseenAndUnread is the default value for NewsFeed.State.
const UnseenAndUnread = 0

// NewsFeed is an ORM for newsfeeds table.
type NewsFeed struct {
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

// TableName returns the name of table.
func (NewsFeed) TableName() string {
	return "newsfeeds"
}

// Shard returns the index of mysql shard instance into which the
// record will be inserted.
func (nf *NewsFeed) Shard() int {
	return int(nf.Kid) % 2
}

// NewNewsFeed creates a NewsFeed instance from a messagejson.NewsFeed.
var NewNewsFeed = func(nf *messagejson.NewsFeed) (n *NewsFeed, err error) {
	n = &NewsFeed{
		NewsfeedID:     nf.ID,
		Kid:            nf.UserID,
		ObservableType: nf.ObservableType,
		ObservableID:   nf.ObservableID,
		EventType:      nf.EventType,
		TargetType:     nf.TargetType,
		TargetID:       nf.TargetID,
		Summary:        nf.Summary,
		State:          UnseenAndUnread,
	}
	return
}
