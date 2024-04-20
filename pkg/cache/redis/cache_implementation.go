package redis

import (
	"github.com/giustech/RavenCache/pkg/cache"
	_redis "github.com/go-redis/redis"
	"log"
	"time"
)

type (
	redis struct {
		client *_redis.Client
	}
)

func GetInstance(addr string, password string, database int) cache.Repository {
	client := _redis.NewClient(&_redis.Options{
		Addr:     addr,
		Password: password,
		DB:       database,
	})
	_, err := client.Ping().Result()
	if err != nil {
		log.Fatalf("Não foi possível conectar ao Redis: %v", err)
	}

	return &redis{
		client: client,
	}
}

func (r redis) Get(key string, ttl int64) string {
	val, err := r.client.Get(key).Result()
	if err != nil {
		//log.Fatalf("Erro ao obter chave %s: %v", key, err)
	} else {
		r.client.Expire(key, time.Duration(ttl)*time.Millisecond)
	}
	return val
}

func (r redis) Put(key string, content []byte, ttl int64) {
	err := r.client.Set(key, content, time.Duration(ttl)*time.Millisecond).Err()
	if err != nil {
		log.Fatalf("Erro ao definir chave %s: %v", key, err)
	}
}
