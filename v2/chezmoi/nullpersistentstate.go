package chezmoi

type nullPersistentState struct{}

func (nullPersistentState) Delete(buckey, key []byte) error        { return nil }
func (nullPersistentState) Get(bucket, key []byte) ([]byte, error) { return nil, nil }
func (nullPersistentState) Set(bucket, key, value []byte) error    { return nil }
