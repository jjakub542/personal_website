package tests

import (
	"fullsteak/internal/user"
	"testing"
)

func TestUserRepository(t *testing.T) {
	var err error

	u := &user.User{
		Email:       "jjakub2d33@gmail.com",
		Password:    "123",
		IsSuperuser: false,
	}

	u.CreatePasswordHash()

	repo := user.NewRepository(TestDB)
	err = repo.CreateOne(u)
	if err != nil {
		t.Fatal(err)
	}

	user2, err := repo.GetOneByEmail(u.Email)

	if err != nil {
		t.Fatal(err)
	}

	if user2.PasswordHash != u.PasswordHash {
		t.Fail()
	}
}
