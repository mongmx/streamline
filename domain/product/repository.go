package product

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"sync"

	"github.com/gomodule/redigo/redis"
	"github.com/google/uuid"
	"github.com/mongmx/streamline/config"
	internalredis "github.com/mongmx/streamline/config/redis"
)

var mu sync.Mutex

// Repository - product store APIs.
type Repository interface {
	Store(product *Product) error
	Save(product *Product) error
	GetByID(id uuid.UUID) (Product, error)
	Streams(ctx context.Context, id uuid.UUID, prodChan chan Product)
}

type repo struct {
	Conn *redis.Pool
	DB   *sql.DB
}

// NewRepository is a factory function of product store.
func NewRepository(conn *redis.Pool, db *sql.DB) Repository {
	return &repo{
		Conn: conn,
	}
}

func (r *repo) Store(product *Product) error {
	b, err := json.Marshal(product)
	if err != nil {
		return err
	}
	key := fmt.Sprintf("product-%s", product.ID)
	c := r.Conn.Get()
	_, err = c.Do("SET", key, string(b))
	return err
}

func (r *repo) Save(product *Product) error {
	prod, err := r.GetByID(product.ID)
	if err != nil {
		return err
	}
	product.CreatedAt = prod.CreatedAt
	b, err := json.Marshal(product)
	if err != nil {
		return err
	}
	key := fmt.Sprintf("product-%s", product.ID)
	c := r.Conn.Get()
	_, err = c.Do("SET", key, string(b))
	return err
}

func (r *repo) GetByID(id uuid.UUID) (Product, error) {
	mu.Lock()
	defer mu.Unlock()

	key := fmt.Sprintf("product-%s", id)
	c := r.Conn.Get()
	b, err := redis.Bytes(c.Do("GET", key))
	if err != nil {
		log.Println("error: ", err)
		return Product{}, err
	}
	log.Printf("bytes: %s\n", string(b))
	var product Product
	if err := json.Unmarshal(b, &product); err != nil {
		log.Println("error: ", err)
		return Product{}, err
	}
	log.Println("product: ", product)
	return product, nil
}

func (r *repo) Streams(ctx context.Context, id uuid.UUID, prodChan chan Product) {
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
			prod, err := r.GetByID(id)
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
