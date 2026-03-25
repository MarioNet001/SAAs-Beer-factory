package inventory

import "time"

type ProductInventory struct {
	ProductID        string    `json:"product_id"`
	StockLevel       int       `json:"stock_level"`
	MinimumThreshold int       `json:"minimum_threshold"`
	UpdatedAt        time.Time `json:"updated_at"`
}
