package cache

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/allegro/bigcache"
	"sync"
	"time"
)

var once sync.Once
var lc *legisCache

func GetLegisCache() *legisCache {
	once.Do(func() {
		c, err := bigcache.NewBigCache(bigcache.DefaultConfig(12 * time.Hour))
		if err != nil {
			panic("WE NEED DA CACHE" + err.Error())
		}
		lc = &legisCache{c: c}
	})
	return lc
}

type legisCache struct {
	c *bigcache.BigCache
}

func (l *legisCache) Delete(key string) (err error) {
	if err = l.c.Delete(key); err != nil {
		err = fmt.Errorf("legisCache delete error: %w", err)
	}
	return
}

func (l *legisCache) AddToCache(key string, object interface{}) {
	objectJson, err := json.Marshal(object)
	if err != nil {
		fmt.Printf("legisCache marshal error %s\nerror: %s ", key, err.Error())
		return
	}
	if err = l.c.Set(key, objectJson); err != nil {
		fmt.Printf("legisCache set error %s\nerror: %s", key, err.Error())
	}
}

func (l *legisCache) GetFromCache(key string, objectToUnmarshal interface{}) error {
	retrievedItem, err := l.c.Get(key)
	if err != nil {
		return err
	}
	if err := json.Unmarshal(retrievedItem, objectToUnmarshal); err != nil {
		return fmt.Errorf("legisCache unmmarshal error: %w")
	}
	fmt.Printf("hit legisCache for %s, returning unmarshaled object \n", key)
	return nil
}

func (l *legisCache) NotFound(err error) bool {
	return errors.Is(err, bigcache.ErrEntryNotFound)
}