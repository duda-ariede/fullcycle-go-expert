package usecase

import "CleanArch/internal/entity"

type ListOrdersOutputDTO struct {
	Orders []entity.Order
	Total  int32
}

type ListOrdersUseCase struct {
	OrderRepository entity.OrderRepositoryInterface
}

func NewListOrdersUseCase(OrderRepository entity.OrderRepositoryInterface) *ListOrdersUseCase {
	return &ListOrdersUseCase{OrderRepository: OrderRepository}
}

func (uc *ListOrdersUseCase) ListOrders() (ListOrdersOutputDTO, error) {
	orders, err := uc.OrderRepository.List()
	if err != nil {
		return ListOrdersOutputDTO{}, err
	}

	total, err := uc.OrderRepository.GetTotal()
	if err != nil {
		return ListOrdersOutputDTO{}, err
	}

	return ListOrdersOutputDTO{
		Orders: orders,
		Total:  int32(total),
	}, nil
}