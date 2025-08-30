package repository

import (
	"ticert/entity"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type OrderRepository interface {
	CreateOrder(order *entity.Order, orderDetails []*entity.OrderDetail, categoryID uuid.UUID, quantity int) error
	GetOrders(page, limit int, userID uuid.UUID) ([]*entity.Order, int64, error)
	GetOrderById(orderID uuid.UUID) (*entity.Order, error)
	GetOrderDetailByTicketCode(ticketCode string) (*entity.OrderDetail, error)
	CancelOrder(orderID uuid.UUID) error
	VerifyOrderStatus(orderID uuid.UUID) error
	VerifyTicket(id uuid.UUID) error
}

type orderRepository struct {
	db *gorm.DB
}

func NewOrderRepository(db *gorm.DB) OrderRepository {
	return &orderRepository{db: db}
}

func (r *orderRepository) CreateOrder(order *entity.Order, orderDetails []*entity.OrderDetail, categoryID uuid.UUID, quantity int) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		var category entity.Category
		if err := tx.Set("gorm:query_option", "FOR UPDATE").
			Where("id = ?", categoryID).
			First(&category).Error; err != nil {
			return err
		}

		if category.Quantity < quantity {
			return gorm.ErrRecordNotFound
		}

		if err := tx.Model(&entity.Category{}).
			Where("id = ?", categoryID).
			Update("quantity", gorm.Expr("quantity - ?", quantity)).Error; err != nil {
			return err
		}

		if category.Quantity-quantity == 0 {
			if err := tx.Model(&entity.Category{}).
				Where("id = ?", categoryID).
				Update("status", "sold").Error; err != nil {
				return err
			}
		}

		if err := tx.Create(order).Error; err != nil {
			return err
		}

		for _, orderDetail := range orderDetails {
			orderDetail.OrderID = order.ID
			if err := tx.Create(orderDetail).Error; err != nil {
				return err
			}
		}

		if err := tx.Preload("OrderDetails").Preload("Category.Event").Preload("User").First(order).Error; err != nil {
			return err
		}

		return nil
	})
}

func (r *orderRepository) GetOrders(page, limit int, userID uuid.UUID) ([]*entity.Order, int64, error) {
	var orders []*entity.Order
	var total int64

	countQuery := r.db.Model(&entity.Order{}).Where("user_id = ?", userID)
	if err := countQuery.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	query := r.db.Preload("OrderDetails").Preload("Category.Event").Preload("User").Where("user_id = ?", userID).Order("created_at DESC")

	if page > 0 && limit > 0 {
		query = query.Offset((page - 1) * limit).Limit(limit)
	}

	if err := query.Find(&orders).Error; err != nil {
		return nil, 0, err
	}

	return orders, total, nil
}

func (r *orderRepository) GetOrderById(orderID uuid.UUID) (*entity.Order, error) {
	var order entity.Order
	if err := r.db.Preload("OrderDetails").Preload("Category.Event").Preload("User").Where("id = ?", orderID).First(&order).Error; err != nil {
		return nil, err
	}
	return &order, nil
}

func (r *orderRepository) GetOrderDetailByTicketCode(ticketCode string) (*entity.OrderDetail, error) {
	var orderDetail entity.OrderDetail
	if err := r.db.Preload("Order").Where("ticket_code = ?", ticketCode).First(&orderDetail).Error; err != nil {
		return nil, err
	}
	return &orderDetail, nil
}

func (r *orderRepository) CancelOrder(orderID uuid.UUID) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		var order entity.Order
		if err := tx.Where("id = ?", orderID).First(&order).Error; err != nil {
			return err
		}

		if order.Quantity > 0 {
			if err := tx.Model(&entity.Category{}).
				Where("id = ?", order.CategoryID).
				Update("quantity", gorm.Expr("quantity + ?", order.Quantity)).Error; err != nil {
				return err
			}
			var category entity.Category
			if err := tx.Where("id = ?", order.CategoryID).First(&category).Error; err == nil {
				if category.Status == "sold" && category.Quantity > 0 {
					if err := tx.Model(&entity.Category{}).
						Where("id = ?", order.CategoryID).
						Update("status", "available").Error; err != nil {
						return err
					}
				}
			}
		}

		if err := tx.Model(&entity.Order{}).
			Where("id = ?", orderID).
			Update("status", "cancelled").Error; err != nil {
			return err
		}

		return nil
	})
}

func (r *orderRepository) VerifyOrderStatus(orderID uuid.UUID) error {
	if err := r.db.Model(&entity.Order{}).Where("id = ?", orderID).Update("status", "paid").Error; err != nil {
		return err
	}
	return nil
}

func (r *orderRepository) VerifyTicket(id uuid.UUID) error {
	if err := r.db.Model(&entity.OrderDetail{}).Where("id = ?", id).Update("redeemed", true).Error; err != nil {
		return err
	}
	return nil
}
