package mysql

// UnseenAndUnread is the default value for NewsFeed.State.
const UnseenAndUnread = 0

// Sharder is an interface for sharddable record.
type Sharder interface {
	Shard() int
}
