package service

import (
	"testing"
	"time"

	"github.com/vicanso/tiny-site/util"
)

func TestGetRedisClient(t *testing.T) {
	client := GetRedisClient()
	if client == nil {
		t.Fatalf("get redis client fail")
	}
}
func TestLock(t *testing.T) {
	key := util.RandomString(8)
	ttl := time.Second
	success, err := Lock(key, ttl)
	if err != nil {
		t.Fatalf("redis lock fail, %v", err)
	}
	if !success {
		t.Fatalf("redis lock fail, it should success")
	}
	success, err = Lock(key, ttl)
	if err != nil {
		t.Fatalf("redis lock fail, %v", err)
	}
	// 第二次由锁失败
	if success {
		t.Fatalf("redis lock twice fail, it should fail")
	}
	time.Sleep(2 * ttl)
	success, err = Lock(key, ttl)
	if err != nil {
		t.Fatalf("redis lock fail, %v", err)
	}
	// 在等待ttl之后，又可以获取锁
	if !success {
		t.Fatalf("redis lock fail(after ttl), it should success")
	}
}

func TestRedisGetSet(t *testing.T) {
	m := map[string]string{
		"a": "1",
	}
	key := util.RandomString(8)
	t.Run("set success", func(t *testing.T) {

		ok, err := RedisSet(key, &m, time.Second)
		if err != nil {
			t.Fatalf("redis set fail, %v", err)
		}
		if !ok {
			t.Fatalf("redis set fail")
		}
	})

	t.Run("get success", func(t *testing.T) {
		tmp := make(map[string]string)
		ok, err := RedisGet(key, &tmp)
		if err != nil {
			t.Fatalf("redis get fail, %v", err)
		}
		if !ok {
			t.Fatalf("redis get fail")
		}
		if tmp["a"] != m["a"] {
			t.Fatalf("redis get data fail")
		}
	})

	t.Run("set fail", func(t *testing.T) {
		_, err := RedisSet("a", m, time.Nanosecond)
		if err == nil {
			t.Fatalf("redis set data 1 ns should be fail")
		}
	})

	t.Run("get fail", func(t *testing.T) {
		client := GetRedisClient()
		data := []byte(`{
			"a": 1,
		}`)
		_, err := client.Set(key, data, time.Second).Result()
		if err != nil {
			t.Fatalf("redis set fail, %v", err)
		}
		tmp := make(map[string]string)
		_, err = RedisGet(key, &tmp)
		if err == nil {
			t.Fatalf("get not json data shoul be fail")
		}

	})

}
