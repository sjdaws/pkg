package modelmock

import (
	"gorm.io/gorm"
)

// ModelMock models.Model compliant model for testing.
type ModelMock struct {
	ID        int
	DeletedAt *gorm.DeletedAt `sql:"index"`
	Test      bool
}

// TableName return the database table for this model.
func (m ModelMock) TableName() string {
	return "model_mocks"
}
