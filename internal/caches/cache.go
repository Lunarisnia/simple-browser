package caches

type CacheMeta struct {
	Value  map[string]string
	MaxAge int
}

type Cache map[string]CacheMeta

type CacheBox struct {
	box Cache
}

func New() CacheBox {
	return CacheBox{
		box: make(Cache),
	}
}

func (c *CacheBox) Set(key string, value map[string]string, maxAge int) {
	c.box[key] = CacheMeta{
		MaxAge: maxAge,
		Value:  value,
	}
}
