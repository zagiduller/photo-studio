package components

import "gorm.io/gorm"

// @project photo-studio
// @created 06.11.2022

type Model struct {
	gorm.Model
	db *gorm.DB
}

func (m *Model) SetDB(db *gorm.DB) {
	m.db = db
}

func (m *Model) GetDB() *gorm.DB {
	if m.db == nil {
		return GetDB()
	}
	return m.db
}
