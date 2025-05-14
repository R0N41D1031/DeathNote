package models

import "gorm.io/gorm"

type Victima struct {
	gorm.Model

	Nombre string   `gorm:"not null" json:"nombre"`
	Foto   string   `json:"foto"`
	Muerte []Muerte `gorm:"foreignKey:VictimaID" json:"muertes"`
}
