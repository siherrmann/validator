package main

import (
	"log"
	"time"
)

type User struct {
	ID              string        `vld:"min3,rex^[a-zA-Z0-9]$"`
	Email           string        `vld:"con@"`
	DurationMinutes time.Duration `vld:"-"`
	CreatedAt       time.Time     `vld:"-"`
}

func main() {
	err := Validate(User{
		ID:    "id.@#",
		Email: "test@email.de",
	})
	if err != nil {
		log.Println(err)
	}
}
