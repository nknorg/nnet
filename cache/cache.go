package cache

// Cache is an anstract cache layer
type Cache interface {
	Add([]byte, interface{}) error
	Get([]byte) (interface{}, bool)
	Set([]byte, interface{}) error
}
