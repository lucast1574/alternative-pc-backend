package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type MaintenanceStatus string

const (
	MaintPending   MaintenanceStatus = "pending"
	MaintContacted MaintenanceStatus = "contacted"
	MaintScheduled MaintenanceStatus = "scheduled"
	MaintCompleted MaintenanceStatus = "completed"
	MaintCancelled MaintenanceStatus = "cancelled"
)

type ServiceType string

const (
	ServiceClean      ServiceType = "limpieza"
	ServiceRepair     ServiceType = "reparacion"
	ServiceUpgrade    ServiceType = "upgrade"
	ServiceDiagnostic ServiceType = "diagnostico"
	ServiceFormat     ServiceType = "formateo"
	ServiceOther      ServiceType = "otro"
)

type MaintenanceRequest struct {
	ID           primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	ServiceType  ServiceType        `bson:"serviceType" json:"serviceType"`
	Description  string             `bson:"description" json:"description"`
	DeviceType   string             `bson:"deviceType" json:"deviceType"`
	DeviceBrand  string             `bson:"deviceBrand" json:"deviceBrand"`
	ClientName   string             `bson:"clientName" json:"clientName"`
	ClientEmail  string             `bson:"clientEmail" json:"clientEmail"`
	ClientPhone  string             `bson:"clientPhone" json:"clientPhone"`
	Status       MaintenanceStatus  `bson:"status" json:"status"`
	AdminNotes   string             `bson:"adminNotes,omitempty" json:"adminNotes,omitempty"`
	UserID       primitive.ObjectID `bson:"userId,omitempty" json:"userId,omitempty"`
	CreatedAt    time.Time          `bson:"createdAt" json:"createdAt"`
	UpdatedAt    time.Time          `bson:"updatedAt" json:"updatedAt"`
}

type CreateMaintenanceInput struct {
	ServiceType ServiceType `json:"serviceType" binding:"required"`
	Description string      `json:"description" binding:"required"`
	DeviceType  string      `json:"deviceType" binding:"required"`
	DeviceBrand string      `json:"deviceBrand"`
	ClientName  string      `json:"clientName" binding:"required"`
	ClientEmail string      `json:"clientEmail" binding:"required,email"`
	ClientPhone string      `json:"clientPhone" binding:"required"`
}
