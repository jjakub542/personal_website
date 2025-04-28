package main

import (
	"fmt"
	"fullsteak/internal/app"
	"fullsteak/internal/user"
	"log"

	_ "github.com/joho/godotenv/autoload"
)

func main() {
	db := app.PostgresClient()
	var err error
	newUser := user.User{IsSuperuser: true}

	fmt.Println("Email: ")
	fmt.Scanln(&newUser.Email)
	fmt.Println("Password: ")
	fmt.Scanln(&newUser.Password)

	err = newUser.Validate()
	if err != nil {
		log.Fatal(err)
	}

	newUser.CreatePasswordHash()

	repo := user.NewRepository(db)

	err = repo.CreateOne(&newUser)

	if err != nil {
		log.Fatal(err)
	}
}
