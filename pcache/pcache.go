package pcache

import (
	"fmt"
	"log"
	"pcache/pcache/view"
	"sync"
)

type Getter interface {
	Get(key string) ([]byte, error)
}

type GetterFunc func(key string) ([]byte, error)

func (fn GetterFunc) Get(key string) ([]byte, error) {
	return fn(key)
}

type Group struct {
	name      string //命名空间
	getter    Getter //回调函数
	mainCache *cache
}

func (g *Group) Get(key string) (view.ByteView, error) {
	if key == "" {
		return view.ByteView{}, fmt.Errorf("key is needed")
	}
	if val, ok := g.mainCache.Get(key); ok {
		log.Println("cache hit")
		return val, nil
	}
	return g.load(key)
}

func (g *Group) load(key string) (view.ByteView, error) {
	return g.getLocally(key)
}

func (g *Group) getLocally(key string) (view.ByteView, error) {
	bytes, err := g.getter.Get(key)
	if err != nil {
		return view.ByteView{}, err

	}
	value := view.NewByteView(bytes)
	g.populateCache(key, value)
	return value, nil
}

func (g *Group) populateCache(key string, value view.ByteView) {
	g.mainCache.Add(key, value)
}

var (
	mu     sync.RWMutex
	groups = make(map[string]*Group)
)

func NewGroup(name string, cacheBytes uint64, getter Getter) *Group {
	if getter == nil {
		panic("getter is nil")
	}
	g := &Group{
		name:   name,
		getter: getter,
		mainCache: &cache{
			cacheBytes: cacheBytes,
		},
	}
	mu.Lock()
	defer mu.Unlock()
	groups[name] = g
	return g
}

func GetGroup(name string) *Group {
	mu.RLock()
	defer mu.RUnlock()
	return groups[name]
}
