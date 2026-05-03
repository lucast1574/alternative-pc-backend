package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/alternative/backend/internal/modules/user/model"
	"github.com/alternative/backend/internal/platform/auth"
	"github.com/alternative/backend/internal/platform/database"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type googleUserInfo struct {
	Sub       string `json:"sub"`
	Email     string `json:"email"`
	Name      string `json:"name"`
	Picture   string `json:"picture"`
	Verified  bool   `json:"email_verified"`
}

func GoogleAuth(c *gin.Context) {
	var input model.GoogleAuthInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Verify token with Google
	resp, err := http.Get(fmt.Sprintf("https://www.googleapis.com/oauth2/v3/userinfo?access_token=%s", input.Token))
	if err != nil || resp.StatusCode != 200 {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "token de Google inválido"})
		return
	}
	defer resp.Body.Close()

	var gUser googleUserInfo
	if err := json.NewDecoder(resp.Body).Decode(&gUser); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "error al verificar Google"})
		return
	}

	coll := database.DB.Collection("users")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Check if user exists by email or googleId
	var user model.User
	err = coll.FindOne(ctx, bson.M{"$or": []bson.M{
		{"email": gUser.Email},
		{"googleId": gUser.Sub},
	}}).Decode(&user)

	if err != nil {
		// New user - create
		user = model.User{
			Name:      gUser.Name,
			Email:     gUser.Email,
			GoogleID:  gUser.Sub,
			AvatarURL: gUser.Picture,
			Role:      model.RoleUser,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}
		result, err := coll.InsertOne(ctx, user)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "error al crear usuario"})
			return
		}
		user.ID = result.InsertedID.(primitive.ObjectID)

		// Google users need to complete profile
		token, _ := auth.GenerateToken(user.ID.Hex(), user.Email, string(user.Role))
		c.JSON(http.StatusCreated, gin.H{
			"token":           token,
			"user":            user.ToResponse(),
			"profileComplete": false,
			"message":         "Completa tu perfil para poder usar todos los servicios",
		})
		return
	}

	// Existing user - update google info if needed
	if user.GoogleID == "" {
		coll.UpdateOne(ctx, bson.M{"_id": user.ID}, bson.M{
			"$set": bson.M{"googleId": gUser.Sub, "avatarUrl": gUser.Picture, "updatedAt": time.Now()},
		})
	}

	profileComplete := user.DNI != "" && user.Phone != "" && user.Address != ""
	token, _ := auth.GenerateToken(user.ID.Hex(), user.Email, string(user.Role))
	c.JSON(http.StatusOK, gin.H{
		"token":           token,
		"user":            user.ToResponse(),
		"profileComplete": profileComplete,
	})
}

func UpdateProfile(c *gin.Context) {
	userID, _ := c.Get("userId")

	var input struct {
		Phone    string `json:"phone"`
		DNI      string `json:"dni"`
		Address  string `json:"address"`
		City     string `json:"city"`
		District string `json:"district"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	objID, _ := primitive.ObjectIDFromHex(userID.(string))
	coll := database.DB.Collection("users")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	update := bson.M{"updatedAt": time.Now()}
	if input.Phone != "" {
		update["phone"] = input.Phone
	}
	if input.DNI != "" {
		update["dni"] = input.DNI
	}
	if input.Address != "" {
		update["address"] = input.Address
	}
	if input.City != "" {
		update["city"] = input.City
	}
	if input.District != "" {
		update["district"] = input.District
	}

	coll.UpdateOne(ctx, bson.M{"_id": objID}, bson.M{"$set": update})

	var user model.User
	coll.FindOne(ctx, bson.M{"_id": objID}).Decode(&user)

	c.JSON(http.StatusOK, user.ToResponse())
}
