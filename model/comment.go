package model

import "time"

type Comment struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	CreatedAt time.Time `json:"-"`
	Title     string    `json:"title"`
	Body      string    `json:"body"`
	UserId    uint      `json:"-"`
	User      User      `gorm:"foreignKey:UserId" json:"-"`
	PostId    uint      `json:"postId"`
	Post      Post      `gorm:"foreignKey:PostId" json:"-"`
}
