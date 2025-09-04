package models

import "time"

// модели данных для json и валидации данных
type Order struct {
	OrderUID          string    `json:"order_uid" validate:"required,min=10,max=50"`
	TrackNumber       string    `json:"track_number" validate:"required,min=1,max=50"`
	Entry             string    `json:"entry" validate:"required,min=1,max=10"`
	Delivery          Delivery  `json:"delivery" validate:"required"`
	Payment           Payment   `json:"payment" validate:"required"`
	Items             []Item    `json:"items" validate:"required,min=1,dive"`
	Locale            string    `json:"locale" validate:"required,len=2"`
	InternalSignature string    `json:"internal_signature"`
	CustomerID        string    `json:"customer_id" validate:"required"`
	DeliveryService   string    `json:"delivery_service" validate:"required,min=1,max=50"`
	Shardkey          string    `json:"shardkey" validate:"required,min=1,max=10"`
	SmID              int       `json:"sm_id" validate:"required,gt=0"`
	DateCreated       time.Time `json:"date_created" validate:"required"`
	OofShard          string    `json:"oof_shard" validate:"required,min=1,max=10"`
}

type Delivery struct {
	OrderUID string `json:"order_uid"`
	Name     string `json:"name" validate:"required,min=2,max=100"`
	Phone    string `json:"phone" validate:"required,e164"`
	Zip      string `json:"zip" validate:"required,min=5,max=10"`
	City     string `json:"city" validate:"required,min=2,max=50"`
	Address  string `json:"address" validate:"required,min=5,max=200"`
	Region   string `json:"region" validate:"required,min=2,max=50"`
	Email    string `json:"email" validate:"required,email"`
}

type Payment struct {
	OrderUID     string `json:"order_uid"`
	Transaction  string `json:"transaction" validate:"required"`
	RequestID    string `json:"request_id"`
	Currency     string `json:"currency" validate:"required,len=3,uppercase"`
	Provider     string `json:"provider" validate:"required,min=1,max=50"`
	Amount       int    `json:"amount" validate:"required,gt=0"`
	PaymentDt    int    `json:"payment_dt" validate:"required,gt=0"`
	Bank         string `json:"bank" validate:"required,min=2,max=50"`
	DeliveryCost int    `json:"delivery_cost" validate:"gte=0"`
	GoodsTotal   int    `json:"goods_total" validate:"required,gt=0"`
	CustomFee    int    `json:"custom_fee" validate:"gte=0"`
}

type Item struct {
	OrderUID    string `json:"order_uid"`
	ChrtID      int    `json:"chrt_id" validate:"required,gt=0"`
	TrackNumber string `json:"track_number" validate:"required,min=1,max=50"`
	Price       int    `json:"price" validate:"required,gt=0"`
	Rid         string `json:"rid" validate:"required"`
	Name        string `json:"name" validate:"required,min=1,max=200"`
	Sale        int    `json:"sale" validate:"gte=0,lte=100"`
	Size        string `json:"size" validate:"required,min=1,max=10"`
	TotalPrice  int    `json:"total_price" validate:"required,gt=0"`
	NmID        int    `json:"nm_id" validate:"required,gt=0"`
	Brand       string `json:"brand" validate:"required,min=1,max=100"`
	Status      int    `json:"status" validate:"required,gte=0"`
}
