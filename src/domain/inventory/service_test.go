package inventory

import (
	"errors"
	"testing"
)

// MockInventoryRepository is a mock of the InventoryRepository interface
type MockInventoryRepository struct {
	GetFunc    func(productID string) (*ProductInventory, error)
	UpdateFunc func(inventory *ProductInventory) error
}

func (m *MockInventoryRepository) Get(productID string) (*ProductInventory, error) {
	return m.GetFunc(productID)
}

func (m *MockInventoryRepository) Update(inventory *ProductInventory) error {
	return m.UpdateFunc(inventory)
}

func TestInventoryService_AdjustStock(t *testing.T) {
	tests := []struct {
		name       string
		productID  string
		amount     int
		mockGet    func(productID string) (*ProductInventory, error)
		mockUpdate func(inventory *ProductInventory) error
		wantErr    bool
	}{
		{
			name:      "success - add stock",
			productID: "p1",
			amount:    10,
			mockGet: func(productID string) (*ProductInventory, error) {
				return &ProductInventory{ProductID: "p1", StockLevel: 20, MinimumThreshold: 5}, nil
			},
			mockUpdate: func(inventory *ProductInventory) error {
				if inventory.StockLevel != 30 {
					t.Errorf("expected 30, got %d", inventory.StockLevel)
				}
				return nil
			},
			wantErr: false,
		},
		{
			name:      "success - consume stock",
			productID: "p1",
			amount:    -5,
			mockGet: func(productID string) (*ProductInventory, error) {
				return &ProductInventory{ProductID: "p1", StockLevel: 20, MinimumThreshold: 5}, nil
			},
			mockUpdate: func(inventory *ProductInventory) error {
				if inventory.StockLevel != 15 {
					t.Errorf("expected 15, got %d", inventory.StockLevel)
				}
				return nil
			},
			wantErr: false,
		},
		{
			name:      "error - insufficient stock",
			productID: "p1",
			amount:    -25,
			mockGet: func(productID string) (*ProductInventory, error) {
				return &ProductInventory{ProductID: "p1", StockLevel: 20, MinimumThreshold: 5}, nil
			},
			mockUpdate: func(inventory *ProductInventory) error {
				return nil
			},
			wantErr: true,
		},
		{
			name:      "error - repo get fails",
			productID: "p1",
			amount:    10,
			mockGet: func(productID string) (*ProductInventory, error) {
				return nil, errors.New("db error")
			},
			mockUpdate: func(inventory *ProductInventory) error {
				return nil
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := &MockInventoryRepository{
				GetFunc:    tt.mockGet,
				UpdateFunc: tt.mockUpdate,
			}
			service := NewInventoryService(repo)
			err := service.AdjustStock(tt.productID, tt.amount)

			if (err != nil) != tt.wantErr {
				t.Errorf("AdjustStock() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
