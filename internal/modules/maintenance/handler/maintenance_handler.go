package handler

import (
	"context"
	"net/http"
	"time"

	"github.com/alternative/backend/internal/modules/maintenance/model"
	"github.com/alternative/backend/internal/platform/database"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func CreateMaintenanceRequest(c *gin.Context) {
	var input model.CreateMaintenanceInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	req := model.MaintenanceRequest{
		ServiceType: input.ServiceType,
		Description: input.Description,
		DeviceType:  input.DeviceType,
		DeviceBrand: input.DeviceBrand,
		ClientName:  input.ClientName,
		ClientEmail: input.ClientEmail,
		ClientPhone: input.ClientPhone,
		Status:      model.MaintPending,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	if userID, exists := c.Get("userId"); exists {
		objID, err := primitive.ObjectIDFromHex(userID.(string))
		if err == nil {
			req.UserID = objID
		}
	}

	coll := database.DB.Collection("maintenance_requests")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	result, err := coll.InsertOne(ctx, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "error al crear solicitud"})
		return
	}

	req.ID = result.InsertedID.(primitive.ObjectID)

	c.JSON(http.StatusCreated, gin.H{
		"request": req,
		"message": "Solicitud de mantenimiento recibida. Nuestro experto te contactará pronto.",
	})
}

func GetMaintenanceRequests(c *gin.Context) {
	coll := database.DB.Collection("maintenance_requests")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	opts := options.Find().SetSort(bson.D{{Key: "createdAt", Value: -1}})
	cursor, err := coll.Find(ctx, bson.M{}, opts)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "error al obtener solicitudes"})
		return
	}
	defer cursor.Close(ctx)

	var requests []model.MaintenanceRequest
	if err := cursor.All(ctx, &requests); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "error al decodificar"})
		return
	}

	if requests == nil {
		requests = []model.MaintenanceRequest{}
	}

	c.JSON(http.StatusOK, requests)
}

func UpdateMaintenanceRequest(c *gin.Context) {
	id := c.Param("id")
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "id inválido"})
		return
	}

	var input struct {
		Status     model.MaintenanceStatus `json:"status"`
		AdminNotes string                  `json:"adminNotes"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	update := bson.M{"updatedAt": time.Now()}
	if input.Status != "" {
		update["status"] = input.Status
	}
	if input.AdminNotes != "" {
		update["adminNotes"] = input.AdminNotes
	}

	coll := database.DB.Collection("maintenance_requests")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	_, err = coll.UpdateOne(ctx, bson.M{"_id": objID}, bson.M{"$set": update})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "error al actualizar"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "solicitud actualizada"})
}
