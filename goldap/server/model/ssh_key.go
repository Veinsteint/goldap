package model

import "gorm.io/gorm"

// SSHKey represents a user's SSH public key
type SSHKey struct {
	gorm.Model
	UserID  uint   `gorm:"type:int;not null;index;comment:'User ID'" json:"userId"`
	Title   string `gorm:"type:varchar(100);not null;comment:'Key title'" json:"title"`
	Key     string `gorm:"type:text;not null;comment:'SSH public key'" json:"key"`
	KeyType string `gorm:"type:varchar(20);comment:'Key type: rsa, ed25519'" json:"keyType"`
	User    User   `gorm:"foreignKey:UserID" json:"user,omitempty"`
}
