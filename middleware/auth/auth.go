package auth

import (
	"encoding/json"
	"github.com/gomodule/redigo/redis"
	"github.com/google/uuid"
	"github.com/gorilla/sessions"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
	"time"
)

// Auth is auth middleware
func Auth(p *redis.Pool) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			sess, err := session.Get("session", c)
			if err != nil {
				c.Error(err)
			}
			token, ok := sess.Values["session_token"].(string)
			if !ok {
				sessionToken := uuid.New().String()
				sess.Options = &sessions.Options{
					Path:     "/",
					MaxAge:   240,
					HttpOnly: true,
				}
				sess.Values["session_token"] = sessionToken
				err = sess.Save(c.Request(), c.Response())
				if err != nil {
					c.Error(err)
				}
			}
			if token != "" {
				c.Logger().Printf("session token: %s", token)
			}
			u := getSessionUser(p, token)
			if u != nil {
				c.Set("user_id", u.ID)
			}
			return next(c)
		}
	}
}

func getSessionUser(pool *redis.Pool, token string) *User {
	conn := pool.Get()
	b, err := redis.Bytes(conn.Do("GET", "user--"+token))
	if err != nil {
		return nil
	}
	var user User
	err = json.Unmarshal(b, &user)
	if err != nil {
		return nil
	}
	return &user
}

type model struct {
	ID        int64      `json:"id"`
	UUID      uuid.UUID  `json:"uuid"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at"`
}

type User struct {
	model
	Email  string `json:"email"`
	PlanID int64  `json:"plan_id"`
}
