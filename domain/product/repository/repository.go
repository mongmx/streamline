package repository

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"sync"

	"github.com/didiyudha/sse-redis/config"
	"github.com/didiyudha/sse-redis/domain/product/model"
	internalredis "github.com/didiyudha/sse-redis/internal/platform/redis"
	"github.com/gomodule/redigo/redis"
	"github.com/google/uuid"
)

var mu sync.Mutex

// ProductCache - product cache APIs.
type ProductCache interface {
	Store(product *model.Product) error
	Save(product *model.Product) error
	GetByID(id uuid.UUID) (model.Product, error)
	Streams(ctx context.Context, id uuid.UUID, prodChan chan model.Product)
}

type productCache struct {
	Conn *redis.Pool
}

// NewProductCache is a factory function of product cache.
func NewProductCache(conn *redis.Pool) ProductCache {
	return &productCache{
		Conn: conn,
	}
}

func (p *productCache) Store(product *model.Product) error {
	b, err := json.Marshal(product)
	if err != nil {
		return err
	}
	key := fmt.Sprintf("product-%s", product.ID)
	c := p.Conn.Get()
	_, err = c.Do("SET", key, string(b))
	return err
}

func (p *productCache) Save(product *model.Product) error {
	prod, err := p.GetByID(product.ID)
	if err != nil {
		return err
	}
	product.CreatedAt = prod.CreatedAt
	b, err := json.Marshal(product)
	if err != nil {
		return err
	}
	key := fmt.Sprintf("product-%s", product.ID)
	c := p.Conn.Get()
	_, err = c.Do("SET", key, string(b))
	return err
}

func (p *productCache) GetByID(id uuid.UUID) (model.Product, error) {
	mu.Lock()
	defer mu.Unlock()

	key := fmt.Sprintf("product-%s", id)
	c := p.Conn.Get()
	b, err := redis.Bytes(c.Do("GET", key))
	if err != nil {
		log.Println("error: ", err)
		return model.Product{}, err
	}
	log.Printf("bytes: %s\n", string(b))
	var product model.Product
	if err := json.Unmarshal(b, &product); err != nil {
		log.Println("error: ", err)
		return model.Product{}, err
	}
	log.Println("product: ", product)
	return product, nil
}

func (p *productCache) Streams(ctx context.Context, id uuid.UUID, prodChan chan model.Product) {
	conn, err := internalredis.NewRedis(config.Cfg.Redis)
	if err != nil {
		log.Println(err)
		close(prodChan)
	}
	c := conn.Get()
	if _, err = c.Do("CONFIG", "SET", "notify-keyspace-events", "KEA"); err != nil {
		close(prodChan)
	}
	psc := redis.PubSubConn{Conn: c}
	key := fmt.Sprintf("product-%s", id)
	keyspace := fmt.Sprintf("__keyspace@*__:%s", key)
	if err := psc.PSubscribe(keyspace, "set"); err != nil {
		log.Println(err)
		close(prodChan)
	}
	for {
		switch m := psc.Receive().(type) {
		case redis.Message:
			log.Printf("message: %v\n", m)
			prod, err := p.GetByID(id)
			if err != nil {
				close(prodChan)
				return
			}
			prodChan <- prod
		case error:
			close(prodChan)
		}
	}
}
