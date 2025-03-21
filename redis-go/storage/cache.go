package storage

import (
	"context"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/redis/go-redis/v9"
	"net/http"
	"reflect"
	"strconv"
	"time"
)

type Card struct {
	ID   int    `json:"id" redis:"id"`
	Name string `json:"name" redis:"name"`
	Data string `json:"date" redis:"date"`
}

func GetCard(ctx context.Context, db *redis.Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(3 * time.Second)
		idStr := chi.URLParam(r, "id")
		if idStr == "" {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		id, err := strconv.Atoi(idStr)
		if err != nil {
			render.Status(r, http.StatusBadRequest)
			return
		}

		card := Card{
			ID:   id,
			Name: "Test",
			Data: "TestData",
		}

		if err := card.ToRedisSet(ctx, db, idStr); err != nil {
			render.Status(r, http.StatusInternalServerError)
			return
		}

		render.Status(r, http.StatusOK)
		render.JSON(w, r, card)
	}

}

func (c *Card) ToRedisSet(ctx context.Context, db *redis.Client, key string) error {
	val := reflect.ValueOf(c).Elem()

	settter := func(p redis.Pipeliner) error {
		for i := 0; i < val.NumField(); i++ {
			field := val.Type().Field(i)
			tag := field.Tag.Get("redis")
			if err := p.HSet(ctx, key, tag, val.Field(i).Interface()).Err(); err != nil {
				return err
			}
		}
		if err := p.Expire(ctx, key, time.Second*30).Err(); err != nil {
			return err
		}
		return nil
	}
	if _, err := db.Pipelined(ctx, settter); err != nil {
		return err
	}
	return nil
}

func CacheMiddleware(ctx context.Context, db *redis.Client) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			idStr := chi.URLParam(r, "id")
			if idStr == "" {
				w.WriteHeader(http.StatusBadRequest)
				return
			}
			data := new(Card)
			if err := db.HGetAll(ctx, idStr).Scan(data); err == nil && (*data != Card{}) {
				render.JSON(w, r, data)
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}
func NewCardHandler(ctx context.Context, db *redis.Client) func(r chi.Router) {
	return func(r chi.Router) {
		r.With(CacheMiddleware(ctx, db)).Get("/{id}", GetCard(ctx, db))
	}
}
