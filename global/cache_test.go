package global

import (
	"testing"

	"github.com/vicanso/tiny-site/util"
)

func TestCache(t *testing.T) {
	t.Run("sync map", func(t *testing.T) {
		key := util.RandomString(8)
		value := "a"
		Store(key, value)
		v, ok := Load(key)
		if !ok || v.(string) != value {
			t.Fatalf("store and load cache fail")
		}
		_, loaded := LoadOrStore(key, "b")
		if !loaded {
			t.Fatalf("load or store should loaded while data exists")
		}

		key = util.RandomString(8)
		v, loaded = LoadOrStore(key, "b")
		if loaded {
			t.Fatalf("load or store should not be loaded while data not exists")
		}
		if v.(string) != "b" {
			t.Fatalf("load or store fail")
		}
	})
}

func TestLRUCache(t *testing.T) {
	key := "vicanso"
	value := "tree.xie"
	t.Run("add", func(t *testing.T) {
		evicted := Add(key, value)
		if evicted {
			t.Fatalf("the first cache should not trigger evict")
		}
	})

	t.Run("get", func(t *testing.T) {
		v, found := Get(key)
		if !found || v.(string) != value {
			t.Fatalf("get from lru cache fail")
		}
	})

	t.Run("remove", func(t *testing.T) {
		Remove(key)
		_, found := Get(key)
		if found {
			t.Fatalf("remove from lru cache fail")
		}
	})

	t.Run("new lru cache", func(t *testing.T) {
		c, err := NewLRU(128)
		if err != nil {
			t.Fatalf("new lru cache fail, %v", err)
		}
		if c == nil {
			t.Fatalf("new lru cache fail")
		}
	})
}

func TestConnectingCount(t *testing.T) {
	SaveConnectingCount(100)
	if GetConnectingCount() != 100 {
		t.Fatalf("save and get connecting count fail")
	}
}

func TestRouteInfos(t *testing.T) {
	routeInfos := make([]map[string]string, 0)
	routeInfos = append(routeInfos, map[string]string{
		"method": "GET",
		"path":   "/users/v1/me",
	})
	SaveRouteInfos(routeInfos)
	data := GetRouteInfos()
	if len(routeInfos) != len(data) {
		t.Fatalf("get route infos fail")
	}
}
