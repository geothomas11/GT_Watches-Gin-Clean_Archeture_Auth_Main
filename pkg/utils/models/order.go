package models

type CheckoutDetails struct {
	AddressInfoResponse []AddressInfoResponse
	Payment_Method      []PaymentDetails
	Cart                []Cart
	Total_Price         float64
}

type OrderFromCart struct {
	PaymentID uint `json:"payment_id" binding:"required"`
	AddressID uint `json:"address_id" binding:"required"`
}
type OrderSuccessResponse struct {
	OrderID        uint   `json:"order_id"`
	ShipmentStatus string `json:"shipment_status"`
}

type OrderDetails struct {
	OrderId        int
	FinalPrice     float64
	ShipmentStatus string
	PaymentStatus  string
}

type OrderIncoming struct {
	UserID    int `json:"user_id"`
	PaymentID int `json:"payment_id"`
	AddressID int `json:"address_id"`
}

type OrderProductDetails struct {
	ProductID   uint    `json:"product_id"`
	ProductName string  `json:"product_name"`
	Quantity    int     `json:"quantity"`
	TOtoalPrice float64 `json:"total_price"`
}

type FullOrderDetails struct {
	OrderDetails        OrderDetails
	OrderProductDetails []OrderProductDetails
}
