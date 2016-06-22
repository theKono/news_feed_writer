package mysql

import (
	"time"

	"github.com/theKono/orchid/model/messagejson"
)

// Notification is an ORM for notifications table.
type Notification struct {
	NotificationID int64 `gorm:"primary_key"`
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

// Shard returns the index of mysql shard instance into which the
// record will be inserted.
func (n *Notification) Shard() int {
	return int(n.Kid) % 2
}

// NewNotification creates a Notification instance from a
// messagejson.Notification.
var NewNotification = func(input *messagejson.Notification) (n *Notification, err error) {
	n = &Notification{
		NotificationID: input.ID,
		Kid:            input.UserID,
		ObservableType: input.ObservableType,
		ObservableID:   input.ObservableID,
		EventType:      input.EventType,
		TargetType:     input.TargetType,
		TargetID:       input.TargetID,
		Summary:        input.Summary,
		State:          UnseenAndUnread,
	}
	return
}
