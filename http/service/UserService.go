package service

import (
	"fmt"

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