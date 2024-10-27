package repository

import (
	"context"
	"order_api/errors"
	"order_api/model"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type OrderRepository struct {
	db    *gorm.DB
	cache Cache
}

type Cache interface {
	GetOrder(ctx context.Context, orderID string) (*model.Order, error)
	SetOrder(ctx context.Context, order *model.Order) error
	DeleteOrder(ctx context.Context, orderID string, userID string) error
}

func NewOrderRepository(db *gorm.DB, cache Cache) *OrderRepository {
	return &OrderRepository{
		db:    db,
		cache: cache,
	}
}

// ListByUserID 获取用户的订单列表
func (r *OrderRepository) ListByUserID(ctx context.Context, userID string) ([]model.Order, error) {
	var orders []model.Order
	if err := r.db.Where("user_id = ?", userID).Find(&orders).Error; err != nil {
		return nil, errors.Wrap(err, "failed to list orders")
	}
	return orders, nil
}

// Create 创建订单
func (r *OrderRepository) Create(ctx context.Context, order *model.Order) error {
	order.ID = uuid.New().String()
	
	if err := r.db.Create(order).Error; err != nil {
		return errors.Wrap(err, "failed to create order")
	}

	return r.cache.SetOrder(ctx, order)
}

// GetByID 根据ID获取订单
func (r *OrderRepository) GetByID(ctx context.Context, orderID string) (*model.Order, error) {
	// 先尝试从缓存获取
	order, err := r.cache.GetOrder(ctx, orderID)
	if err == nil {
		return order, nil
	}

	// 缓存未命中，从数据库获取
	var dbOrder model.Order
	if err := r.db.Preload("Items").First(&dbOrder, "id = ?", orderID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.ErrOrderNotFound
		}
		return nil, errors.Wrap(err, "failed to get order")
	}

	// 写入缓存
	if err := r.cache.SetOrder(ctx, &dbOrder); err != nil {
		return nil, errors.Wrap(err, "failed to cache order")
	}

	return &dbOrder, nil
}

// Update 更新订单
func (r *OrderRepository) Update(ctx context.Context, order *model.Order) error {
	if err := r.db.Save(order).Error; err != nil {
		return errors.Wrap(err, "failed to update order")
	}

	return r.cache.SetOrder(ctx, order)
}

// Delete 删除订单
func (r *OrderRepository) Delete(ctx context.Context, orderID string) error {
	var order model.Order
	if err := r.db.First(&order, "id = ?", orderID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return errors.ErrOrderNotFound
		}
		return errors.Wrap(err, "failed to find order")
	}

	if err := r.db.Delete(&order).Error; err != nil {
		return errors.Wrap(err, "failed to delete order")
	}

	return r.cache.DeleteOrder(ctx, orderID, order.UserID)
}