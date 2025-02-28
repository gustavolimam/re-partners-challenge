package models

type Order struct {
	Items int `json:"items"`
}

type OrderResponse struct {
	OrderItems int          `json:"order_items"`
	OrderPacks []OrderPacks `json:"order_packs"`
}

type OrderPacks struct {
	Size  int `json:"size"`
	Count int `json:"count"`
}
