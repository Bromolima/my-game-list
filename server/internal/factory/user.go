package factory

import (
	"github.com/Bromolima/my-game-list/internal/entities"
	"github.com/Bromolima/my-game-list/internal/http/dto"
	"github.com/google/uuid"
)

func NewUser(email, hashedPassword, username, avatarURL string, roleID uint) *entities.User {
	return &entities.User{
		ID:        uuid.New(),
		Email:     email,
		Password:  hashedPassword,
		Username:  username,
		AvatarURL: avatarURL,
		RoleID:    roleID,
	}
}

func NewUserUpdate(id uuid.UUID, email, hashedPassword, username, avatarURL string) *entities.User {
	return &entities.User{
		ID:        id,
		Email:     email,
		Password:  hashedPassword,
		Username:  username,
		AvatarURL: avatarURL,
	}
}

func NewUserClaims(id uuid.UUID, roleID uint) *entities.User {
	return &entities.User{
		ID:     id,
		RoleID: roleID,
	}
}

func NewResponseFromUser(user *entities.User) *dto.UserResponse {
	return &dto.UserResponse{
		Username: user.Username,
	}
}
