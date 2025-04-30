package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/xxii22w/Social/internal/store"
)

type UserStorage struct {
	rdb *redis.Client
}

const UserExpTime = time.Minute

func (s *UserStorage) Get(ctx context.Context, userID int64) (*store.User, error) {
	cacheKey := fmt.Sprintf("user-%d", userID)

	data, err := s.rdb.Get(ctx, cacheKey).Result()
	if err == redis.Nil {
		return nil, nil
	} else if err != nil {
		return nil, err
	}

	var user store.User
	if data != "" {
		err := json.Unmarshal([]byte(data), &user)
		if err != nil {
			return nil, err
		}
	}
	return &user, nil
}

func(s *UserStorage) Set(ctx context.Context, user *store.User) error {
	cacheKey := fmt.Sprintf("user-%d",user.ID)

	json,err := json.Marshal(user)
	if err != nil {
		return err
	}

	return s.rdb.SetEX(ctx,cacheKey,json,UserExpTime).Err()
}

func(s *UserStorage) Delete(ctx context.Context,userID int64) {
	cacheKey := fmt.Sprintf("user-%d",userID)
	s.rdb.Del(ctx,cacheKey)
}