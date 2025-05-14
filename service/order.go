package service

import (
	"context"
	repo "unit-of-work/repository"
)

type OrderService struct {
	orderRepo repo.OrderRepository
	itemRepo  repo.ItemRepository
}

func NewOrderService(orderRepo repo.OrderRepository, itemRepo repo.ItemRepository) *OrderService {
	return &OrderService{
		orderRepo: orderRepo,
		itemRepo:  itemRepo,
	}
}

func (s *OrderService) DeleteOrder(ctx context.Context, orderID int) error {
	// First delete the items related to the order
	if err := s.itemRepo.DeleteItemsByOrderID(ctx, orderID); err != nil {
		return err
	}

	// Then delete the order itself
	if err := s.orderRepo.DeleteOrder(ctx, orderID); err != nil {
		return err
	}

	return nil
}
