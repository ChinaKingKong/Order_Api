package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"order_api/config"
	"order_api/errors"
	"order_api/model"
	"sync"
	"time"

	"github.com/go-redis/redis/v8"
)

type Cache struct {
	localCache sync.Map
	redis      *redis.Client
	config     *config.Config
}

func NewCache(cfg *config.RedisConfig) (*Cache, error) {
	client := redis.NewClient(&redis.Options{
		Addr:         fmt.Sprintf("%s:%s", cfg.Host, cfg.Port),
		Password:     cfg.Password,
		DB:           cfg.DB,
		PoolSize:     cfg.PoolSize,
		MaxRetries:   cfg.MaxRetries,
		MinIdleConns: cfg.MaxIdleConns,
	})

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := client.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("Redis连接失败: %w", err)
	}

	return &Cache{
		localCache: sync.Map{},
		redis:      client,
	}, nil
}

// GetOrder 获取订单信息
func (c *Cache) GetOrder(ctx context.Context, orderID string) (*model.Order, error) {
	// 1. 先查本地缓存
	if value, ok := c.localCache.Load(orderID); ok {
		return value.(*model.Order), nil
	}

	// 2. 查Redis缓存
	data, err := c.redis.Get(ctx, c.getOrderKey(orderID)).Bytes()
	if err == nil {
		var order model.Order
		if err := json.Unmarshal(data, &order); err == nil {
			// 写入本地缓存
			c.localCache.Store(orderID, &order)
			return &order, nil
		}
	}

	return nil, errors.New("cache miss")
}

// SetOrder 将订单信息写入缓存
func (c *Cache) SetOrder(ctx context.Context, order *model.Order) error {
	data, err := json.Marshal(order)
	if err != nil {
		return errors.Wrap(err, "订单序列化失败")
	}

	// 使用管道批量执行Redis命令
	pipe := c.redis.Pipeline()
	pipe.Set(ctx, c.getOrderKey(order.ID), data, 30*time.Minute)
	pipe.SAdd(ctx, c.getUserOrdersKey(order.UserID), order.ID)

	if _, err := pipe.Exec(ctx); err != nil {
		return errors.Wrap(err, "缓存写入失败")
	}

	// 更新本地缓存
	c.localCache.Store(order.ID, order)
	return nil
}

// DeleteOrder 从缓存中删除订单信息
func (c *Cache) DeleteOrder(ctx context.Context, orderID string, userID string) error {
	pipe := c.redis.Pipeline()
	pipe.Del(ctx, c.getOrderKey(orderID))
	pipe.SRem(ctx, c.getUserOrdersKey(userID), orderID)

	if _, err := pipe.Exec(ctx); err != nil {
		return errors.Wrap(err, "缓存删除失败")
	}

	c.localCache.Delete(orderID)
	return nil
}

// Close 关闭缓存连接
func (c *Cache) Close() error {
	return c.redis.Close()
}

// getOrderKey 生成订单缓存键
func (c *Cache) getOrderKey(orderID string) string {
	return fmt.Sprintf("order:%s", orderID)
}

// getUserOrdersKey 生成用户订单列表缓存键
func (c *Cache) getUserOrdersKey(userID string) string {
	return fmt.Sprintf("user:%s:orders", userID)
}