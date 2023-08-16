package model

import "time"

type Post struct {
	ID        uint `json:"id" gorm:"primaryKey"`
	CreatedAt time.Time
	Title     string `json:"title"`
	Body      string `json:"body"`
	UserId    uint   `json:"userId"`
	User      User   `gorm:"foreignKey:UserId" json:"-"`
}
