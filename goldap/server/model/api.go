package model

import "gorm.io/gorm"

// Api represents an API endpoint for authorization
type Api struct {
	gorm.Model
	Method   string `gorm:"type:varchar(20);comment:'HTTP method'" json:"method"`
	Path     string `gorm:"type:varchar(100);comment:'API path'" json:"path"`
	Category string `gorm:"type:varchar(50);comment:'Category'" json:"category"`
	Remark   string `gorm:"type:varchar(100);comment:'Description'" json:"remark"`
	Creator  string `gorm:"type:varchar(20);comment:'Creator'" json:"creator"`
}
