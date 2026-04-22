package model

import (
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

// FieldRelation stores field mapping relationships
type FieldRelation struct {
	gorm.Model
	Flag       string
	Attributes datatypes.JSON
}
