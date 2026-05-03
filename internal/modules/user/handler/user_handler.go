package handler

import (
	"context"
	"net/http"
	"time"

	"github.com/alternative/backend/internal/modules/user/model"
	"github.com/alternative/backend/internal/platform/auth"
	"github.com/alternative/backend/internal/platform/database"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
)

func Register(c *gin.Context) {
	var input model.RegisterInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	coll := database.DB.Collection("users")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Check if email exists
	var existing model.User
	err := coll.FindOne(ctx, bson.M{"email": input.Email}).Decode(&existing)
	if err == nil {
		c.JSON(http.StatusConflict, gin.H{"error": "el email ya está registrado"})
		return
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "error al procesar contraseña"})
		return
	}

	user := model.User{
		Name:      input.Name,
		Email:     input.Email,
		Password:  string(hashedPassword),
		Role:      model.RoleUser,
		Phone:     input.Phone,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	result, err := coll.InsertOne(ctx, user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "error al crear usuario"})
		return
	}

	user.ID = result.InsertedID.(primitive.ObjectID)

	token, err := auth.GenerateToken(user.ID.Hex(), user.Email, string(user.Role))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "error al generar token"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"token": token,
		"user":  user.ToResponse(),
	})
}

func Login(c *gin.Context) {
	var input model.LoginInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	coll := database.DB.Collection("users")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var user model.User
	err := coll.FindOne(ctx, bson.M{"email": input.Email}).Decode(&user)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "credenciales inválidas"})
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.Password))
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "credenciales inválidas"})
		return
	}

	token, err := auth.GenerateToken(user.ID.Hex(), user.Email, string(user.Role))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "error al generar token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"token": token,
		"user":  user.ToResponse(),
	})
}

func GetMe(c *gin.Context) {
	userID, _ := c.Get("userId")

	coll := database.DB.Collection("users")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	objID, err := primitive.ObjectIDFromHex(userID.(string))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "id inválido"})
		return
	}

	var user model.User
	err = coll.FindOne(ctx, bson.M{"_id": objID}).Decode(&user)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "usuario no encontrado"})
		return
	}

	c.JSON(http.StatusOK, user.ToResponse())
}

func UpdateRole(c *gin.Context) {
	targetID := c.Param("id")

	var input struct {
		Role model.Role `json:"role" binding:"required"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if input.Role != model.RoleUser && input.Role != model.RoleSeller && input.Role != model.RoleAdmin {
		c.JSON(http.StatusBadRequest, gin.H{"error": "rol inválido"})
		return
	}

	coll := database.DB.Collection("users")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	objID, err := primitive.ObjectIDFromHex(targetID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "id inválido"})
		return
	}

	_, err = coll.UpdateOne(ctx, bson.M{"_id": objID}, bson.M{
		"$set": bson.M{"role": input.Role, "updatedAt": time.Now()},
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "error al actualizar rol"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "rol actualizado"})
}

func ListUsers(c *gin.Context) {
	coll := database.DB.Collection("users")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	cursor, err := coll.Find(ctx, bson.M{})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "error al listar usuarios"})
		return
	}
	defer cursor.Close(ctx)

	var users []model.UserResponse
	for cursor.Next(ctx) {
		var user model.User
		if err := cursor.Decode(&user); err != nil {
			continue
		}
		users = append(users, user.ToResponse())
	}

	if users == nil {
		users = []model.UserResponse{}
	}

	c.JSON(http.StatusOK, users)
}
