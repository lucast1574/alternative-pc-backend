package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type SellRequestStatus string

const (
	SellPending  SellRequestStatus = "pending"
	SellReviewed SellRequestStatus = "reviewed"
	SellAccepted SellRequestStatus = "accepted"
	SellRejected SellRequestStatus = "rejected"
)

type PartType string

const (
	// Equipos completos
	PartPCDesktop   PartType = "pc_desktop"
	PartLaptop      PartType = "laptop"
	PartAllInOne    PartType = "all_in_one"
	// Componentes
	PartCPU         PartType = "cpu"
	PartGPU         PartType = "gpu"
	PartRAM         PartType = "ram"
	PartMotherboard PartType = "motherboard"
	PartStorage     PartType = "storage"
	PartPSU         PartType = "psu"
	PartCase        PartType = "case"
	PartCooler      PartType = "cooler"
	// Periféricos
	PartMonitor     PartType = "monitor"
	PartKeyboard    PartType = "keyboard"
	PartMouse       PartType = "mouse"
	PartOther       PartType = "otro"
)

// Estimados base por tipo (en soles)
var PriceEstimates = map[PartType][2]float64{
	PartPCDesktop:   {100, 600},
	PartLaptop:      {80, 500},
	PartAllInOne:    {80, 450},
	PartCPU:         {30, 150},
	PartGPU:         {50, 300},
	PartRAM:         {15, 60},
	PartMotherboard: {25, 120},
	PartStorage:     {15, 80},
	PartPSU:         {15, 60},
	PartCase:        {10, 40},
	PartCooler:      {10, 30},
	PartMonitor:     {30, 150},
	PartKeyboard:    {5, 25},
	PartMouse:       {5, 15},
	PartOther:       {5, 100},
}

type SellRequest struct {
	ID            primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	PartType      PartType           `bson:"partType" json:"partType"`
	Brand         string             `bson:"brand" json:"brand"`
	Model         string             `bson:"model" json:"model"`
	Condition     string             `bson:"condition" json:"condition"`
	Description   string             `bson:"description" json:"description"`
	SellerName    string             `bson:"sellerName" json:"sellerName"`
	SellerEmail   string             `bson:"sellerEmail" json:"sellerEmail"`
	SellerPhone   string             `bson:"sellerPhone" json:"sellerPhone"`
	EstimateMin   float64            `bson:"estimateMin" json:"estimateMin"`
	EstimateMax   float64            `bson:"estimateMax" json:"estimateMax"`
	OfferedAmount float64            `bson:"offeredAmount,omitempty" json:"offeredAmount,omitempty"`
	Status        SellRequestStatus  `bson:"status" json:"status"`
	AdminNotes    string             `bson:"adminNotes,omitempty" json:"adminNotes,omitempty"`
	UserID        primitive.ObjectID `bson:"userId,omitempty" json:"userId,omitempty"`
	CreatedAt     time.Time          `bson:"createdAt" json:"createdAt"`
	UpdatedAt     time.Time          `bson:"updatedAt" json:"updatedAt"`
}

type CreateSellRequestInput struct {
	PartType    PartType `json:"partType" binding:"required"`
	Brand       string   `json:"brand" binding:"required"`
	Model       string   `json:"model" binding:"required"`
	Condition   string   `json:"condition" binding:"required"`
	Description string   `json:"description"`
}
