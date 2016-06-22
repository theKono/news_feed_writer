package mysql

// Sharder is an interface for sharddable record.
type Sharder interface {
	Shard() int
}
