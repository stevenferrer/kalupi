package tx

// Tx wraps the commit and rollback method
type Tx interface {
	Commit() error
	Rollback() error
}
