package main

import (
	"gorm.io/gorm"
	"gorm.io/driver/sqlite"
	"time"
)

type Users struct {
	ID           uint      `gorm:"primaryKey"`
	Username     string    `json:"username" gorm:"type:varchar(40)"`
	PasswordHash string    `json:"password" gorm:"type:varchar(255)"`
	Role         string    `json:"role" gorm:"type:varchar(20)"`
	CreatedAt    time.Time
}


type Task struct {
    ID			uint   `gorm:"primaryKey"`
	Title		string `json:"title" gorm:"unique;not null"`
	Description	string `json:"description"`
    Completed	bool   `json:"completed"`
}

func setupDatabase() *gorm.DB {
	db, err := gorm.Open(sqlite.Open("task.db"), &gorm.Config{})
	if err != nil {
		panic("Failed to connect to the database")
	}
	db.AutoMigrate(&Task{}, &Users{})
	return db
}