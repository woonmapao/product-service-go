package initializer

import (
	"fmt"
	"log"
	"os"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func DBInitializer() {

	dsn := os.Getenv("DB_URL")
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("failed to connect database: %v", err)
	}

	fmt.Println(`
	 ***********************
	*  Services Started !!  *
	 ***********************


	 ██████╗ ██████╗  ██████╗ ██████╗ ██╗   ██╗ ██████╗████████╗
	 ██╔══██╗██╔══██╗██╔═══██╗██╔══██╗██║   ██║██╔════╝╚══██╔══╝
	 ██████╔╝██████╔╝██║   ██║██║  ██║██║   ██║██║        ██║   
	 ██╔═══╝ ██╔══██╗██║   ██║██║  ██║██║   ██║██║        ██║   
	 ██║     ██║  ██║╚██████╔╝██████╔╝╚██████╔╝╚██████╗   ██║   
	 ╚═╝     ╚═╝  ╚═╝ ╚═════╝ ╚═════╝  ╚═════╝  ╚═════╝   ╚═╝   
	 
	 
	 
	`)

	var now time.Time
	db.Raw("SELECT NOW()").Scan(&now)

	log.Printf("Current database time: %v", now)

	DB = db
}
