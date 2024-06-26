package core

import (
	"sync"

	"github.com/zncdatadev/dolphinscheduler-operator/pkg/util"
)

var MergedCache = NewMapCache()

type MapCache struct {
	data map[string]interface{}
	lock sync.Mutex
}

func NewMapCache() *MapCache {
	return &MapCache{
		data: make(map[string]interface{}),
	}
}

func (c *MapCache) Set(key string, value interface{}) {
	c.lock.Lock()
	defer c.lock.Unlock()

	c.data[key] = value
}

func (c *MapCache) Get(key string) (interface{}, bool) {
	c.lock.Lock()
	defer c.lock.Unlock()

	value, ok := c.data[key]
	return value, ok
}

func (c *MapCache) Del(key string) {
	c.lock.Lock()
	defer c.lock.Unlock()

	delete(c.data, key)
}

func ReleaseCache() {
	MergedCache.lock.Lock()
	defer MergedCache.lock.Unlock()

	MergedCache.data = make(map[string]interface{})
}

func CreateRoleCfgCacheKey(instanceName string, role Role, groupName string) string {
	return util.NewResourceNameGenerator(instanceName, string(role), groupName).GenerateResourceName("cache")
}

func StoreSingleGroupConfig(instanceName string, role Role, groupName string, cfg any) {
	key := CreateRoleCfgCacheKey(instanceName, role, groupName)
	MergedCache.Set(key, cfg)
}

func GetRoleGroup(instanceName string, role Role, groupName string) any {
	key := CreateRoleCfgCacheKey(instanceName, role, groupName)
	value, _ := MergedCache.Get(key)
	return value
}
