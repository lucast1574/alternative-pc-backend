package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Role string

const (
	RoleUser   Role = "user"
	RoleSeller Role = "seller"
	RoleAdmin  Role = "admin"
)

type User struct {
	ID         primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Name       string             `bson:"name" json:"name"`
	Email      string             `bson:"email" json:"email"`
	Password   string             `bson:"password" json:"-"`
	Role       Role               `bson:"role" json:"role"`
	Phone      string             `bson:"phone,omitempty" json:"phone,omitempty"`
	DNI        string             `bson:"dni,omitempty" json:"dni,omitempty"`
	Address    string             `bson:"address,omitempty" json:"address,omitempty"`
	City       string             `bson:"city,omitempty" json:"city,omitempty"`
	District   string             `bson:"district,omitempty" json:"district,omitempty"`
	GoogleID   string             `bson:"googleId,omitempty" json:"-"`
	AvatarURL  string             `bson:"avatarUrl,omitempty" json:"avatarUrl,omitempty"`
	CreatedAt  time.Time          `bson:"createdAt" json:"createdAt"`
	UpdatedAt  time.Time          `bson:"updatedAt" json:"updatedAt"`
}

type RegisterInput struct {
	Name     string `json:"name" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
	Phone    string `json:"phone" binding:"required"`
	DNI      string `json:"dni" binding:"required"`
	Address  string `json:"address" binding:"required"`
	City     string `json:"city" binding:"required"`
	District string `json:"district"`
}

type LoginInput struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type UserResponse struct {
	ID        primitive.ObjectID `json:"id"`
	Name      string             `json:"name"`
	Email     string             `json:"email"`
	Role      Role               `json:"role"`
	Phone     string             `json:"phone,omitempty"`
	DNI       string             `json:"dni,omitempty"`
	Address   string             `json:"address,omitempty"`
	City      string             `json:"city,omitempty"`
	District  string             `json:"district,omitempty"`
	AvatarURL string             `json:"avatarUrl,omitempty"`
}

func (u *User) ToResponse() UserResponse {
	return UserResponse{
		ID:        u.ID,
		Name:      u.Name,
		Email:     u.Email,
		Role:      u.Role,
		Phone:     u.Phone,
		DNI:       u.DNI,
		Address:   u.Address,
		City:      u.City,
		District:  u.District,
		AvatarURL: u.AvatarURL,
	}
}

type GoogleAuthInput struct {
	Token string `json:"token" binding:"required"`
}
