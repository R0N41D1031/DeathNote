package models

import "gorm.io/gorm"

type Victimas struct {
	gorm.Model

	Nombre string `gorm:"not null"`
	Foto   string
}
