package handler

import (
	"context"
	"net/http"
	"time"

	"github.com/alternative/backend/internal/modules/listing/model"
	usermodel "github.com/alternative/backend/internal/modules/user/model"
	"github.com/alternative/backend/internal/platform/database"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func CreateListing(c *gin.Context) {
	var input model.CreateListingInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID, _ := c.Get("userId")
	role, _ := c.Get("role")

	objID, _ := primitive.ObjectIDFromHex(userID.(string))

	// Get seller name
	var user usermodel.User
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	database.DB.Collection("users").FindOne(ctx, bson.M{"_id": objID}).Decode(&user)

	status := model.StatusPending
	if role == "admin" {
		status = model.StatusApproved
	}

	listing := model.Listing{
		Title:       input.Title,
		Description: input.Description,
		Price:       input.Price,
		Category:    input.Category,
		Condition:   input.Condition,
		Images:      input.Images,
		Status:      status,
		SellerID:    objID,
		SellerName:  user.Name,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	if listing.Images == nil {
		listing.Images = []string{}
	}

	coll := database.DB.Collection("listings")
	result, err := coll.InsertOne(ctx, listing)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "error al crear publicación"})
		return
	}

	listing.ID = result.InsertedID.(primitive.ObjectID)
	c.JSON(http.StatusCreated, listing)
}

func GetApprovedListings(c *gin.Context) {
	coll := database.DB.Collection("listings")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	opts := options.Find().SetSort(bson.D{{Key: "createdAt", Value: -1}})
	cursor, err := coll.Find(ctx, bson.M{"status": model.StatusApproved}, opts)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "error al obtener publicaciones"})
		return
	}
	defer cursor.Close(ctx)

	var listings []model.Listing
	if err := cursor.All(ctx, &listings); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "error al decodificar"})
		return
	}

	if listings == nil {
		listings = []model.Listing{}
	}

	c.JSON(http.StatusOK, listings)
}

func GetPendingListings(c *gin.Context) {
	coll := database.DB.Collection("listings")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	opts := options.Find().SetSort(bson.D{{Key: "createdAt", Value: -1}})
	cursor, err := coll.Find(ctx, bson.M{"status": model.StatusPending}, opts)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "error al obtener pendientes"})
		return
	}
	defer cursor.Close(ctx)

	var listings []model.Listing
	if err := cursor.All(ctx, &listings); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "error al decodificar"})
		return
	}

	if listings == nil {
		listings = []model.Listing{}
	}

	c.JSON(http.StatusOK, listings)
}

func GetMyListings(c *gin.Context) {
	userID, _ := c.Get("userId")
	objID, _ := primitive.ObjectIDFromHex(userID.(string))

	coll := database.DB.Collection("listings")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	opts := options.Find().SetSort(bson.D{{Key: "createdAt", Value: -1}})
	cursor, err := coll.Find(ctx, bson.M{"sellerId": objID}, opts)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "error al obtener publicaciones"})
		return
	}
	defer cursor.Close(ctx)

	var listings []model.Listing
	if err := cursor.All(ctx, &listings); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "error al decodificar"})
		return
	}

	if listings == nil {
		listings = []model.Listing{}
	}

	c.JSON(http.StatusOK, listings)
}

func ApproveListing(c *gin.Context) {
	id := c.Param("id")
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "id inválido"})
		return
	}

	coll := database.DB.Collection("listings")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	_, err = coll.UpdateOne(ctx, bson.M{"_id": objID}, bson.M{
		"$set": bson.M{"status": model.StatusApproved, "updatedAt": time.Now()},
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "error al aprobar"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "publicación aprobada"})
}

func RejectListing(c *gin.Context) {
	id := c.Param("id")
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "id inválido"})
		return
	}

	coll := database.DB.Collection("listings")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	_, err = coll.UpdateOne(ctx, bson.M{"_id": objID}, bson.M{
		"$set": bson.M{"status": model.StatusRejected, "updatedAt": time.Now()},
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "error al rechazar"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "publicación rechazada"})
}

func GetListing(c *gin.Context) {
	id := c.Param("id")
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "id inválido"})
		return
	}

	coll := database.DB.Collection("listings")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var listing model.Listing
	err = coll.FindOne(ctx, bson.M{"_id": objID}).Decode(&listing)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "publicación no encontrada"})
		return
	}

	c.JSON(http.StatusOK, listing)
}
