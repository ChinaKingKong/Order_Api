// Package model 定义了订单系统的数据模型
package model

import (
	"errors"
	"time"

	"gorm.io/gorm"
)

// 订单状态常量定义
const (
	StatusPending   = "pending"   // 待支付
	StatusPaid      = "paid"      // 已支付
	StatusShipped   = "shipped"   // 已发货
	StatusDelivered = "delivered" // 已送达
	StatusCancelled = "cancelled" // 已取消
)

// Order 订单模型
type Order struct {
	ID        string         `json:"id" gorm:"primaryKey;type:varchar(36)" label:"订单ID"`
	UserID    string         `json:"user_id" gorm:"type:varchar(36);index;not null" validate:"required" label:"用户ID"`
	Status    string         `json:"status" gorm:"type:varchar(20);default:pending" validate:"required,order_status" label:"订单状态"`
	Amount    float64        `json:"amount" gorm:"type:decimal(10,2)" validate:"gte=0" label:"订单金额"`
	Items     []OrderItem    `json:"items" gorm:"foreignKey:OrderID" validate:"required,dive" label:"订单项"`
	CreatedAt time.Time      `json:"created_at" label:"创建时间"`
	UpdatedAt time.Time      `json:"updated_at" label:"更新时间"`
	DeletedAt gorm.DeletedAt `json:"deleted_at,omitempty" gorm:"index" label:"删除时间"`
}

// OrderItem 订单项模型
type OrderItem struct {
	ID        string    `json:"id" gorm:"primaryKey;type:varchar(36)" label:"订单项ID"`
	OrderID   string    `json:"order_id" gorm:"type:varchar(36);index;not null" label:"订单ID"`
	ProductID string    `json:"product_id" gorm:"type:varchar(36);not null" validate:"required" label:"商品ID"`
	Quantity  int       `json:"quantity" gorm:"not null" validate:"required,gt=0" label:"商品数量"`
	Price     float64   `json:"price" gorm:"type:decimal(10,2);not null" validate:"required,gt=0" label:"商品价格"`
	CreatedAt time.Time `json:"created_at" label:"创建时间"`
	UpdatedAt time.Time `json:"updated_at" label:"更新时间"`
}

// CalculateAmount 计算订单总金额
func (o *Order) CalculateAmount() {
	var total float64
	for _, item := range o.Items {
		total += item.Price * float64(item.Quantity)
	}
	o.Amount = total
}

// Validate 验证订单数据
func (o *Order) Validate() error {
	if o.UserID == "" {
		return errors.New("用户ID不能为空")
	}

	if len(o.Items) == 0 {
		return errors.New("订单项不能为空")
	}

	validStatus := map[string]bool{
		StatusPending:   true,
		StatusPaid:      true,
		StatusShipped:   true,
		StatusDelivered: true,
		StatusCancelled: true,
	}

	if !validStatus[o.Status] {
		return errors.New("无效的订单状态")
	}

	for _, item := range o.Items {
		if item.ProductID == "" {
			return errors.New("商品ID不能为空")
		}
		if item.Quantity <= 0 {
			return errors.New("商品数量必须大于0")
		}
		if item.Price <= 0 {
			return errors.New("商品价格必须大于0")
		}
	}

	return nil
}

// IsValidStatusTransition 检查订单状态转换是否有效
func IsValidStatusTransition(currentStatus, newStatus string) bool {
	// 定义有效的状态转换
	validTransitions := map[string]map[string]bool{
		StatusPending: {
			StatusPaid:      true,
			StatusCancelled: true,
		},
		StatusPaid: {
			StatusShipped:   true,
			StatusCancelled: true,
		},
		StatusShipped: {
			StatusDelivered: true,
		},
		StatusDelivered: {},      // 终态，不能转换
		StatusCancelled: {},      // 终态，不能转换
	}

	if transitions, exists := validTransitions[currentStatus]; exists {
		return transitions[newStatus]
	}
	return false
}

// CanCancel 检查订单是否可以取消
func CanCancel(status string) bool {
	return status == StatusPending || status == StatusPaid
}

// CanDelete 检查订单是否可以删除
func CanDelete(status string) bool {
	return status == StatusCancelled || status == StatusDelivered
}

// BeforeCreate GORM 钩子，在创建前执行
func (o *Order) BeforeCreate(tx *gorm.DB) error {
	if o.Status == "" {
		o.Status = StatusPending
	}
	return nil
}

// BeforeUpdate GORM 钩子，在更新前执行
func (o *Order) BeforeUpdate(tx *gorm.DB) error {
	return o.Validate()
}

// String 实现 Stringer 接口，用于日志输出
func (o *Order) String() string {
	return "Order{ID: " + o.ID + ", UserID: " + o.UserID + ", Status: " + o.Status + "}"
}