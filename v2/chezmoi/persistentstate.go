package chezmoi

// FIXME should not need both PersistentState and PersistentStateWriter
// FIXME PersistentState should be just PersistentStateReader + Delete + Set

// A PersistentStateReader reads a persistent state.
type PersistentStateReader interface {
	Get(bucket, key []byte) ([]byte, error)
}

// A PersistentStateWriter writes to a persistent state.
type PersistentStateWriter interface {
	Delete(bucket, key []byte) error
	Set(bucket, key, value []byte) error
}

// A PersistentState is a persistent state.
type PersistentState interface {
	PersistentStateReader
	PersistentStateWriter
}
