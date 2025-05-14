package models

import "gorm.io/gorm"

type Muerte struct {
	gorm.Model

	Causa       string `json:"causa"`
	Descripcion string `json:"descripcion"`
	VictimaID   uint   `json:"victima_id"`
}
