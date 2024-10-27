package service

import (
	"context"
	"order_api/errors"
	"order_api/model"
	"order_api/repository"
)

type OrderService struct {
	repo *repository.OrderRepository
}

func NewOrderService(repo *repository.OrderRepository) *OrderService {
	return &OrderService{repo: repo}
}

// ListOrders 获取用户的订单列表
func (s *OrderService) ListOrders(ctx context.Context, userID string) ([]model.Order, error) {
	return s.repo.ListByUserID(ctx, userID)
}

// CreateOrder 创建订单
func (s *OrderService) CreateOrder(ctx context.Context, order *model.Order) error {
	if err := order.Validate(); err != nil {
		return errors.Wrap(err, "order validation failed")
	}

	order.CalculateAmount()
	order.Status = model.StatusPending

	if err := s.repo.Create(ctx, order); err != nil {
		return errors.Wrap(err, "failed to create order")
	}

	return nil
}

// GetOrder 获取订单详情
func (s *OrderService) GetOrder(ctx context.Context, orderID, userID string) (*model.Order, error) {
	order, err := s.repo.GetByID(ctx, orderID)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get order")
	}

	// 验证订单所属权
	if order.UserID != userID {
		return nil, errors.ErrUnauthorized
	}

	return order, nil
}

// UpdateOrderStatus 更新订单状态
func (s *OrderService) UpdateOrderStatus(ctx context.Context, orderID, userID, newStatus string) error {
	order, err := s.GetOrder(ctx, orderID, userID)
	if err != nil {
		return err
	}

	if !model.IsValidStatusTransition(order.Status, newStatus) {
		return errors.Wrap(errors.ErrInvalidOrderStatus, "invalid status transition")
	}

	order.Status = newStatus
	if err := s.repo.Update(ctx, order); err != nil {
		return errors.Wrap(err, "failed to update order")
	}

	return nil
}

// DeleteOrder 删除订单
func (s *OrderService) DeleteOrder(ctx context.Context, orderID, userID string) error {
	order, err := s.GetOrder(ctx, orderID, userID)
	if err != nil {
		return err
	}

	if !model.CanDelete(order.Status) {
		return errors.Wrap(errors.ErrInvalidOrderStatus, "order cannot be deleted")
	}

	if err := s.repo.Delete(ctx, orderID); err != nil {
		return errors.Wrap(err, "failed to delete order")
	}

	return nil
}