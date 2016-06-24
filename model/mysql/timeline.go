// +build timeline

package mysql

import (
	"time"

	"github.com/theKono/orchid/model/messagejson"
)

// Timeline is an ORM for newsfeeds table.
type Timeline struct {
	TimelineID  int64 `gorm:"primary_key"`
	Kid         int32
	Summary     string `gorm:"size:16383"`
	InScrapbook bool
	CreatedAt   time.Time
}

// Shard returns the index of mysql shard instance into which the
// record will be inserted.
func (nf *Timeline) Shard() int {
	return int(nf.Kid) % 2
}

// NewTimeline creates a Timeline instance from a messagejson.Timeline.
var NewTimeline = func(input *messagejson.Timeline) (t *Timeline, err error) {
	t = &Timeline{
		TimelineID:  input.ID,
		Kid:         input.UserID,
		Summary:     input.Summary,
		InScrapbook: false,
	}
	return
}
