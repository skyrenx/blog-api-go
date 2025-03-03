package service

import (
	"fmt"

	"github.com/skyrenx/blog-api-go/http/entities"
	"github.com/skyrenx/blog-api-go/http/entities/dto"
	"github.com/skyrenx/blog-api-go/http/repository"
)

func GetUserByUsername(username string) (*dto.UserWithoutPassword, error) {
	r, err := repository.GetUserByUsername(username)
	if err != nil {
		fmt.Printf("Error in GetUserByUsername: %v\n", err.Error())
		return nil, fmt.Errorf("could not get the username of the user: %v", username)
	}
	return r, nil
}

func Register(user entities.User) error {
	err := repository.RegisterUser(user)
	if err != nil {
		fmt.Printf("Error in Register: %v\n", err.Error())
		return fmt.Errorf("could not register the user: %v", user.Username)
	}
	return nil
}

func Login(user entities.User) (*string, error) {
	token, err := repository.Login(user)
	if err != nil {
		fmt.Printf("Error in Login: %v\n", err.Error())
		return nil, fmt.Errorf("could not login the user: %v", user.Username)
	}
	return token, nil
}