package models

import (
	"gorm.io/gorm"
)

type Company struct {
	gorm.Model
	CompanyName    string `json:"company_name"`
	CompanyDetails string `json:"company_details"`
}
