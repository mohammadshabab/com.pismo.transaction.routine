package models

type Account struct {
	AccountID      int    `json:"account_id" gorm:"primary_key;auto_increment"`
	DocumentNumber string `json:"document_number" gorm:"unique;not null"`
}
