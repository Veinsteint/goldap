package model

import (
	"gorm.io/gorm"
)

// OperationLog records API operation history
type OperationLog struct {
	gorm.Model
	Username   string `gorm:"type:varchar(20);comment:'Username'" json:"username"`
	Ip         string `gorm:"type:varchar(20);comment:'IP address'" json:"ip"`
	IpLocation string `gorm:"type:varchar(20);comment:'IP location'" json:"ipLocation"`
	Method     string `gorm:"type:varchar(20);comment:'HTTP method'" json:"method"`
	Path       string `gorm:"type:varchar(100);comment:'API path'" json:"path"`
	Remark     string `gorm:"type:varchar(100);comment:'Description'" json:"remark"`
	Status     int    `gorm:"type:int(4);comment:'Response status'" json:"status"`
	StartTime  string `gorm:"type:varchar(2048);comment:'Start time'" json:"startTime"`
	TimeCost   int64  `gorm:"type:int(6);comment:'Duration (ms)'" json:"timeCost"`
	UserAgent  string `gorm:"type:varchar(2048);comment:'User agent'" json:"userAgent"`
}
