package handler

import (
	"order_api/errors"
	"order_api/model"
	"order_api/service"

	"github.com/gin-gonic/gin"
)

type OrderHandler struct {
	orderService *service.OrderService
}

func NewOrderHandler(orderService *service.OrderService) *OrderHandler {
	return &OrderHandler{
		orderService: orderService,
	}
}

// ListOrders 获取订单列表
func (h *OrderHandler) ListOrders(c *gin.Context) {
	userID := c.GetString("user_id")
	orders, err := h.orderService.ListOrders(c.Request.Context(), userID)
	if err != nil {
		ServerError(c, err)
		return
	}
	Success(c, orders)
}

// CreateOrder 创建订单
func (h *OrderHandler) CreateOrder(c *gin.Context) {
	var order model.Order
	if err := c.ShouldBindJSON(&order); err != nil {
		ValidationError(c, []string{"请求数据格式错误"})
		return
	}

	userID := c.GetString("user_id")
	if userID == "" {
		Unauthorized(c)
		return
	}
	order.UserID = userID

	if err := h.orderService.CreateOrder(c.Request.Context(), &order); err != nil {
		ServerError(c, err)
		return
	}

	Created(c, order)
}

// GetOrder 获取订单详情
func (h *OrderHandler) GetOrder(c *gin.Context) {
	orderID := c.Param("id")
	userID := c.GetString("user_id")

	order, err := h.orderService.GetOrder(c.Request.Context(), orderID, userID)
	if err != nil {
		switch {
		case errors.Is(err, errors.ErrOrderNotFound):
			NotFound(c, "订单不存在")
		case errors.Is(err, errors.ErrUnauthorized):
			Forbidden(c)
		default:
			ServerError(c, err)
		}
		return
	}

	Success(c, order)
}

// UpdateOrder 更新订单状态
func (h *OrderHandler) UpdateOrder(c *gin.Context) {
	orderID := c.Param("id")
	userID := c.GetString("user_id")

	var updateReq struct {
		Status string `json:"status" binding:"required,oneof=pending paid shipped delivered cancelled"`
	}

	if err := c.ShouldBindJSON(&updateReq); err != nil {
		ValidationError(c, []string{"无效的订单状态"})
		return
	}

	if err := h.orderService.UpdateOrderStatus(c.Request.Context(), orderID, userID, updateReq.Status); err != nil {
		switch {
		case errors.Is(err, errors.ErrOrderNotFound):
			NotFound(c, "订单不存在")
		case errors.Is(err, errors.ErrUnauthorized):
			Forbidden(c)
		case errors.Is(err, errors.ErrInvalidOrderStatus):
			ValidationError(c, []string{"订单状态变更无效"})
		default:
			ServerError(c, err)
		}
		return
	}

	Success(c, gin.H{"message": "订单状态更新成功"})
}

// DeleteOrder 删除订单
func (h *OrderHandler) DeleteOrder(c *gin.Context) {
	orderID := c.Param("id")
	userID := c.GetString("user_id")

	if err := h.orderService.DeleteOrder(c.Request.Context(), orderID, userID); err != nil {
		switch {
		case errors.Is(err, errors.ErrOrderNotFound):
			NotFound(c, "订单不存在")
		case errors.Is(err, errors.ErrUnauthorized):
			Forbidden(c)
		case errors.Is(err, errors.ErrInvalidOrderStatus):
			ValidationError(c, []string{"当前订单状态不允许删除"})
		default:
			ServerError(c, err)
		}
		return
	}

	Success(c, gin.H{"message": "订单删除成功"})
}
