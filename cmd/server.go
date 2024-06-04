package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/dukivt/rate-limiter/configs"
	"github.com/dukivt/rate-limiter/internal/infra/web/handlers"
	"github.com/dukivt/rate-limiter/internal/infra/web/middlewares"
	"github.com/dukivt/rate-limiter/internal/infra/web/webserver"
	"github.com/dukivt/rate-limiter/pkg/logger"
	"github.com/dukivt/rate-limiter/pkg/ratelimit"
	"github.com/go-redis/redis/v8"
	redisEventStorage "github.com/verdade/go-expert-ratelimit/pkg/ratelimit/redis"
)

func main() {
	configs, err := configs.LoadConfig(".")
	if err != nil {
		panic(err)
	}

	rdb := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", configs.RedisHost, configs.RedisPort),
		Password: configs.RedisPassword,
		DB:       configs.RedisDB,
	})
	fmt.Printf("%s:%s", configs.RedisHost, configs.RedisPort)

	redisEventStorage := redisEventStorage.NewRedisEventStorage(rdb)

	rlIp, err := ratelimit.New(redisEventStorage, "ip", configs.IPConfigLimit.MaxRequests, time.Duration(configs.IPConfigLimit.BlockTimeSecond)*time.Second)
	if err != nil {
		logger.Error("error por IP", err)
		return
	}

	rlToken, err := ratelimit.New(redisEventStorage, "token", 0, 0*time.Second)
	if err != nil {
		logger.Error("error por Token", err)
		return
	}

	m := middlewares.NewLimiter(rlToken, rlIp, configs.TokensConfigLimit)

	ws := webserver.New(configs.WebServerPort)
	h := handlers.NewHealthHandler()

	ws.AddHandler("/health", m.RateLimiter(http.HandlerFunc(h.HealthHandler)))
	ws.Start()
}
