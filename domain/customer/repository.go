package customer

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/gomodule/redigo/redis"
	"github.com/google/uuid"
	"github.com/mongmx/streamline/config"
	internalRedis "github.com/mongmx/streamline/config/redis"
	"log"
	"sync"
)

var mu sync.Mutex

// Repository - customer store APIs.
type Repository interface {
	AllCustomer() ([]Customer, error)
	GetCustomer(customerUUID uuid.UUID) (*Customer, error)
	Streams(ctx context.Context, id uuid.UUID, chatChan chan ChatMessage)
	CreateCustomer(customer *Customer) (int64, error)
}

type repo struct {
	DB        *sql.DB
	RedisPool *redis.Pool
}

// NewRepository is a factory function of customer store.
func NewRepository(db *sql.DB, pool *redis.Pool) Repository {
	return &repo{
		DB:        db,
		RedisPool: pool,
	}
}

func (r *repo) AllCustomer() ([]Customer, error) {
	customers := make([]Customer, 0)
	cUUID, err := uuid.Parse("564f1966-0e2d-40e5-b25a-545a4f4ae7ff")
	if err != nil {
		return nil, err
	}
	customer := Customer{
		ID:   1,
		UUID: cUUID,
		Name: "John Doe",
	}
	customers = append(customers, customer)
	return customers, nil
}

func (r *repo) GetCustomer(_ uuid.UUID) (*Customer, error) {
	cUUID, err := uuid.Parse("564f1966-0e2d-40e5-b25a-545a4f4ae7ff")
	if err != nil {
		return nil, err
	}
	customer := &Customer{
		ID:   1,
		UUID: cUUID,
		Name: "John Doe",
	}
	return customer, nil
}

func (r *repo) getMessageByCustomerUUID(id uuid.UUID) (*ChatMessage, error) {
	mu.Lock()
	defer mu.Unlock()

	key := fmt.Sprintf("customer-%s", id)
	rp := r.RedisPool.Get()
	b, err := redis.Bytes(rp.Do("GET", key))
	if err != nil {
		log.Println("error: ", err)
		return nil, err
	}
	//log.Printf("bytes: %s\n", string(b))
	var chatMessage *ChatMessage
	if err := json.Unmarshal(b, &chatMessage); err != nil {
		log.Println("error: ", err)
		return nil, err
	}
	//log.Println("product: ", chatMessage)
	return chatMessage, nil
}

func (r *repo) Streams(_ context.Context, id uuid.UUID, chatChan chan ChatMessage) {
	conn, err := internalRedis.NewRedis(config.Cfg.Redis)
	if err != nil {
		log.Println(err)
		close(chatChan)
	}
	c := conn.Get()
	if _, err = c.Do("CONFIG", "SET", "notify-keyspace-events", "KEA"); err != nil {
		close(chatChan)
	}
	psc := redis.PubSubConn{Conn: c}
	key := fmt.Sprintf("customer-%s", id)
	keyspace := fmt.Sprintf("__keyspace@*__:%s", key)
	if err := psc.PSubscribe(keyspace, "set"); err != nil {
		log.Println(err)
		close(chatChan)
	}
	for {
		switch psc.Receive().(type) {
		case redis.Message:
			//log.Printf("message: %v\n", m)
			chatMessage, err := r.getMessageByCustomerUUID(id)
			if err != nil {
				close(chatChan)
				return
			}
			chatChan <- *chatMessage
		case error:
			close(chatChan)
		}
	}
}

func (r *repo) CreateCustomer(customer *Customer) (int64, error) {
	cm := ChatMessage{
		Message: "hello customer",
	}
	b, err := json.Marshal(cm)
	if err != nil {
		return 0, err
	}
	key := fmt.Sprintf("customer-%s", customer.UUID)
	rp := r.RedisPool.Get()
	_, err = rp.Do("SET", key, string(b))
	if err != nil {
		return 0, err
	}
	return customer.ID, nil
}
