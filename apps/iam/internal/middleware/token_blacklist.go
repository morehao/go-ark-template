package middleware

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"time"

	"github.com/morehao/goark/pkg/dbclient"
	"github.com/redis/go-redis/v9"
)

const tokenBlacklistPrefix = "iam:token:blacklist:"

func tokenHash(token string) string {
	h := sha256.Sum256([]byte(token))
	return hex.EncodeToString(h[:])
}

func AddToBlacklist(ctx context.Context, token string, expiry time.Time) error {
	ttl := time.Until(expiry)
	if ttl <= 0 {
		return nil
	}
	key := fmt.Sprintf("%s%s", tokenBlacklistPrefix, tokenHash(token))
	return dbclient.RedisCli.Set(ctx, key, "1", ttl).Err()
}

func IsBlacklisted(ctx context.Context, token string) bool {
	key := fmt.Sprintf("%s%s", tokenBlacklistPrefix, tokenHash(token))
	val, err := dbclient.RedisCli.Get(ctx, key).Result()
	if err == redis.Nil {
		return false
	}
	return val == "1"
}
