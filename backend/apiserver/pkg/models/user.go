package models

import (
	"context"
	"log"
)

type User struct {
	BaseModel
	Name string
}

func (u *User) Save(ctx context.Context) {
	db := GetDB()
	log.Println("name:", u.Name, "db:", db)
	err := db.Create(u).Error
	if err != nil {
		log.Println("create user failed.", "err:", err.Error())
	}
}

func SaveUserInBatch(ctx context.Context, users *[]User) error {
	db := GetDB()
	err := db.Create(users).Error
	if err != nil {
		log.Println("create user failed.", "err:", err.Error())
	}
	return nil
}
