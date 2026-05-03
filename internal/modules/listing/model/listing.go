package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ListingStatus string

const (
	StatusPending  ListingStatus = "pending"
	StatusApproved ListingStatus = "approved"
	StatusRejected ListingStatus = "rejected"
)

type ListingCategory string

const (
	CategoryPC          ListingCategory = "pc_armada"
	CategoryComponent   ListingCategory = "componente"
	CategoryPeripheral  ListingCategory = "periferico"
	CategoryAccessory   ListingCategory = "accesorio"
	CategoryOther       ListingCategory = "otro"
)

type Listing struct {
	ID          primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Title       string             `bson:"title" json:"title"`
	Description string             `bson:"description" json:"description"`
	Price       float64            `bson:"price" json:"price"`
	Category    ListingCategory    `bson:"category" json:"category"`
	Condition   string             `bson:"condition" json:"condition"`
	Images      []string           `bson:"images,omitempty" json:"images"`
	Status      ListingStatus      `bson:"status" json:"status"`
	SellerID    primitive.ObjectID `bson:"sellerId" json:"sellerId"`
	SellerName  string             `bson:"sellerName" json:"sellerName"`
	CreatedAt   time.Time          `bson:"createdAt" json:"createdAt"`
	UpdatedAt   time.Time          `bson:"updatedAt" json:"updatedAt"`
}

type CreateListingInput struct {
	Title       string          `json:"title" binding:"required"`
	Description string          `json:"description" binding:"required"`
	Price       float64         `json:"price" binding:"required,gt=0"`
	Category    ListingCategory `json:"category" binding:"required"`
	Condition   string          `json:"condition" binding:"required"`
	Images      []string        `json:"images"`
}
