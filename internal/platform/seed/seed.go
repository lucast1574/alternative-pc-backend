package seed

import (
	"context"
	"log"
	"time"

	"github.com/alternative/backend/internal/modules/listing/model"
	"github.com/alternative/backend/internal/platform/database"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func SeedMockListings() {
	coll := database.DB.Collection("listings")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	count, _ := coll.CountDocuments(ctx, bson.M{})
	if count > 0 {
		return // Already has listings
	}

	adminID := primitive.NewObjectID()

	listings := []interface{}{
		model.Listing{
			Title:       "PC Gamer Ryzen 5 5600X + RTX 3060",
			Description: "PC armada con Ryzen 5 5600X, RTX 3060 12GB, 16GB RAM DDR4, SSD 512GB NVMe, fuente 650W 80+ Bronze. Ideal para gaming en 1080p a máximos.",
			Price:       1850,
			Category:    model.CategoryPC,
			Condition:   "Buen estado",
			Images:      []string{},
			Status:      model.StatusApproved,
			SellerID:    adminID,
			SellerName:  "Alternative PC",
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		},
		model.Listing{
			Title:       "Laptop HP Pavilion 15 - Core i5 11va Gen",
			Description: "Laptop HP Pavilion con procesador Intel Core i5-1135G7, 8GB RAM, SSD 256GB, pantalla 15.6\" Full HD. Batería dura ~5 horas. Incluye cargador original.",
			Price:       1200,
			Category:    model.CategoryPC,
			Condition:   "Como nuevo",
			Images:      []string{},
			Status:      model.StatusApproved,
			SellerID:    adminID,
			SellerName:  "Alternative PC",
			CreatedAt:   time.Now().Add(-1 * time.Hour),
			UpdatedAt:   time.Now(),
		},
		model.Listing{
			Title:       "GTX 1660 Super 6GB GDDR6",
			Description: "Tarjeta de video NVIDIA GTX 1660 Super, 6GB GDDR6. Nunca minada, usada solo para gaming. Funciona perfectamente.",
			Price:       380,
			Category:    model.CategoryComponent,
			Condition:   "Buen estado",
			Images:      []string{},
			Status:      model.StatusApproved,
			SellerID:    adminID,
			SellerName:  "Alternative PC",
			CreatedAt:   time.Now().Add(-2 * time.Hour),
			UpdatedAt:   time.Now(),
		},
		model.Listing{
			Title:       "Monitor Samsung 24\" Full HD 75Hz",
			Description: "Monitor Samsung de 24 pulgadas, resolución 1920x1080, panel IPS, 75Hz, FreeSync. Con cable HDMI incluido.",
			Price:       320,
			Category:    model.CategoryPeripheral,
			Condition:   "Uso normal",
			Images:      []string{},
			Status:      model.StatusApproved,
			SellerID:    adminID,
			SellerName:  "Alternative PC",
			CreatedAt:   time.Now().Add(-3 * time.Hour),
			UpdatedAt:   time.Now(),
		},
		model.Listing{
			Title:       "Kit RAM Corsair Vengeance 16GB (2x8) DDR4 3200MHz",
			Description: "Kit de memoria RAM DDR4 Corsair Vengeance LPX, 2 módulos de 8GB a 3200MHz. Con perfil XMP.",
			Price:       120,
			Category:    model.CategoryComponent,
			Condition:   "Como nuevo",
			Images:      []string{},
			Status:      model.StatusApproved,
			SellerID:    adminID,
			SellerName:  "Alternative PC",
			CreatedAt:   time.Now().Add(-4 * time.Hour),
			UpdatedAt:   time.Now(),
		},
		model.Listing{
			Title:       "PC Oficina Intel i3 10100 + SSD 240GB",
			Description: "PC ideal para oficina y estudio. Intel Core i3-10100, 8GB RAM DDR4, SSD 240GB, fuente 450W. Windows 11 activado.",
			Price:       750,
			Category:    model.CategoryPC,
			Condition:   "Buen estado",
			Images:      []string{},
			Status:      model.StatusApproved,
			SellerID:    adminID,
			SellerName:  "Alternative PC",
			CreatedAt:   time.Now().Add(-5 * time.Hour),
			UpdatedAt:   time.Now(),
		},
		model.Listing{
			Title:       "Teclado mecánico Redragon Kumara K552 RGB",
			Description: "Teclado mecánico TKL con switches Outemu Blue, retroiluminación RGB. Anti-ghosting. Cable USB trenzado.",
			Price:       85,
			Category:    model.CategoryPeripheral,
			Condition:   "Como nuevo",
			Images:      []string{},
			Status:      model.StatusApproved,
			SellerID:    adminID,
			SellerName:  "Alternative PC",
			CreatedAt:   time.Now().Add(-6 * time.Hour),
			UpdatedAt:   time.Now(),
		},
		model.Listing{
			Title:       "SSD Kingston A400 480GB SATA",
			Description: "Disco SSD Kingston A400 de 480GB, interfaz SATA III, velocidad de lectura 500MB/s. Perfecto para upgrade.",
			Price:       75,
			Category:    model.CategoryComponent,
			Condition:   "Nuevo",
			Images:      []string{},
			Status:      model.StatusApproved,
			SellerID:    adminID,
			SellerName:  "Alternative PC",
			CreatedAt:   time.Now().Add(-7 * time.Hour),
			UpdatedAt:   time.Now(),
		},
	}

	result, err := coll.InsertMany(ctx, listings)
	if err != nil {
		log.Printf("Error seeding listings: %v", err)
		return
	}
	log.Printf("Seeded %d mock listings", len(result.InsertedIDs))
}
