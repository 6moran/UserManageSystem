package model

import "time"

type User struct {
	ID         int
	Username   string
	Email      string
	Password   string
	Status     int
	Role       string
	Avatar     string
	CreateTime *time.Time
	LastTime   *time.Time
}
