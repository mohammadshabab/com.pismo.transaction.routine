package models

import (
	"time"
)

type Transaction struct {
	ID              uint      `gorm:"primaryKey;autoIncrement" json:"transaction_id"`
	AccountID       uint      `gorm:"foreignKey:AccountID" json:"account_id"`
	OperationTypeID int       `json:"operation_type_id"`
	Amount          float64   `json:"amount"`
	EventDate       time.Time `json:"event_date"`
}
