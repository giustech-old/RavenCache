package cache

type (
	Repository interface {
		Get(key string, ttl int64) string
		Put(key string, content []byte, ttl int64)
	}
)
