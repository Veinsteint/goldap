package model

import (
	"time"

	"gorm.io/gorm"
)

// PendingUser represents a user registration pending approval
type PendingUser struct {
	gorm.Model
	Username     string     `gorm:"type:varchar(50);not null;unique;comment:'Username'" json:"username"`
	Password     string     `gorm:"type:text;not null;comment:'Password (encrypted)'" json:"password"`
	Nickname     string     `gorm:"type:varchar(50);not null;comment:'Full name'" json:"nickname"`
	Mail         string     `gorm:"type:varchar(100);not null;unique;comment:'Email'" json:"mail"`
	Remark       string     `gorm:"type:varchar(255);comment:'Registration note'" json:"remark"`
	Status       uint       `gorm:"type:tinyint(1);default:0;comment:'Status: 0=Pending, 1=Approved, 2=Rejected'" json:"status"`
	Reviewer     string     `gorm:"type:varchar(20);comment:'Reviewer'" json:"reviewer"`
	ReviewRemark string     `gorm:"type:varchar(255);comment:'Review note'" json:"reviewRemark"`
	ReviewedAt   *time.Time `gorm:"type:datetime(3);comment:'Review time'" json:"reviewedAt"`
}
