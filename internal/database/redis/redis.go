package redis

import (
	"context"
	"time"

	"github.com/gofrs/uuid"
	"github.com/lunn06/smart-toy-backend/internal/config"
	"github.com/lunn06/smart-toy-backend/internal/models"
	"github.com/redis/go-redis/v9"
)

var (
	dragonflyCtx = context.Background()
	rdb          *redis.Client
)

func Init() {
	rdb = redis.NewClient(&redis.Options{
		Addr:     config.CFG.Redis.Address,
		Password: config.CFG.Redis.Password,
		DB:       config.CFG.Redis.Database,
	})
}

func InsertRefreshToken(userId int, fingerprint string, tokenLife time.Duration) (string, error) {
	tokenUuidStruct, err := uuid.NewV7()
	if err != nil {
		return "", err
	}
    tokenUuid := tokenUuidStruct.String()

    jwtToken := models.JwtToken {
        UserId: userId,
        FingerPrint: fingerprint,
        CreationTime: time.Now(),
    }

	_, err = rdb.Set(dragonflyCtx, tokenUuid, jwtToken, tokenLife).Result()
	if err != nil {
		return "", err
	}

	return tokenUuid, nil
}

func GetRefreshToken(tokenUuid string) (models.JwtToken, error) {
    var token models.JwtToken

    err := rdb.Get(dragonflyCtx, tokenUuid).Scan(&token)
    if err != nil {
        return token, err
    }

    return token, nil
}

func PopRefreshToken(tokenUuid string) (models.JwtToken, error) {
    token, err := GetRefreshToken(tokenUuid)

    if err != nil {
        return token, err
    }

    _, err = rdb.Del(dragonflyCtx, tokenUuid).Result()

    if err != nil {
        return token, err
    }
    return token, nil
}

func DelRefreshToken(tokenUuid string) error {
    _, err := rdb.Del(dragonflyCtx, tokenUuid).Result()

    if err != nil {
        return err
    }
    return nil
}