package inventory

type InventoryRepository interface {
	Get(productID string) (*ProductInventory, error)
	Update(inventory *ProductInventory) error
}
