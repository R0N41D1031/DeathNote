package models

import "gorm.io/gorm"

type Muertes struct {
	gorm.Model

	Causa       string
	Descripcion string
}
