package auth

import (
	"github.com/gomodule/redigo/redis"
	"github.com/jmoiron/sqlx"
)

// Repository - auth store APIs.
type Repository interface {
	CreateUser(user *User) (*User, error)
	CreateAuth(auth *Auth) (*Auth, error)
	CreateTopic(topic *Topic) (*Topic, error)
	StoreSessionToken(token string, user *User) error
}

type repo struct {
	DB        *sqlx.DB
	RedisPool *redis.Pool
}

// NewRepository is a factory function of auth store.
func NewRepository(db *sqlx.DB, pool *redis.Pool) Repository {
	return &repo{
		DB:        db,
		RedisPool: pool,
	}
}

func (r repo) CreateUser(user *User) (*User, error) {
	query := `
		INSERT INTO streamline.users (email, plan_id) VALUES ($1, $2) RETURNING id
	`
	var lastInsertedID int64
	err := r.DB.QueryRowx(query, user.Email, user.PlanID).Scan(&lastInsertedID)
	if err != nil {
		return nil, err
	}
	user.ID = lastInsertedID
	return user, nil
}

func (r repo) CreateAuth(auth *Auth) (*Auth, error) {
	query := `
		INSERT INTO streamline.auths (user_id, type, secret) VALUES ($1, $2, $3)
	`
	_, err := r.DB.Exec(query, auth.UserID, auth.Type, auth.Secret)
	if err != nil {
		return nil, err
	}
	return auth, nil
}

func (r repo) CreateTopic(topic *Topic) (*Topic, error) {
	query := `
		INSERT INTO streamline.topics (user_id, title) VALUES ($1, $2)
	`
	_, err := r.DB.Exec(query, topic.UserID, topic.Title)
	if err != nil {
		return nil, err
	}
	return topic, nil
}

func (r repo) StoreSessionToken(token string, user *User) error {
	c := r.RedisPool.Get()
	_, err := c.Do("SET", "user--"+token, user)
	if err != nil {
		return err
	}
	return nil
}
