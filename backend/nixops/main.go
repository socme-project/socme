package nixops

import (
	"fmt"
	"log"

	"github.com/joho/godotenv"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func main() {
	err := godotenv.Load("../../.env")
	if err != nil {
		fmt.Println("Error loading .env file in the root folder: ", err)
		return
	}
	db, err := gorm.Open(sqlite.Open("../backend.db"), &gorm.Config{})
	if err != nil {
		log.Fatal("Error connecting to database:", err)
		return
	}

	clients, err := GetAllClients(db)
	if err != nil {
		log.Fatal("Error getting clients:", err)
	}

	results := SendAll(clients, "uptime")
	for name, output := range results {
		fmt.Printf("[%s] %s\n", name, output)
	}
}
