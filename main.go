package main

import (
	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
	"github.com/orange-cpp/virtualizersdk-redux/v2"
	"hammer-down-server/db"
	"hammer-down-server/handlers"
	"log"
	"net/http"
)

func init() {
	if err := godotenv.Load(); err != nil {
		panic("No .env file found")
	}
}

var stealth_area = [virtualizersdk.STEALTH_ONE_MB * 2]uint32{
	virtualizersdk.STEALTH_ARRAY_FLAG1,
	virtualizersdk.STEALTH_ARRAY_FLAG2,
	virtualizersdk.STEALTH_ARRAY_FLAG3,
	virtualizersdk.STEALTH_ARRAY_FLAG4,
	virtualizersdk.STEALTH_ARRAY_FLAG5,
	virtualizersdk.STEALTH_ARRAY_FLAG6,
	virtualizersdk.STEALTH_ARRAY_FLAG7,
	virtualizersdk.STEALTH_ARRAY_FLAG8}

func main() {
	virtualizersdk.Macro(virtualizersdk.FALCON_TINY_START)

	database, err := db.NewDB()
	if err != nil {
		log.Fatalf("DB init error: %v", err)
	}

	http.HandleFunc("/detect", func(w http.ResponseWriter, r *http.Request) {
		handlers.DetectCheat(w, r, database)
	})
	log.Println("Starting server on :80")
	log.Fatal(http.ListenAndServe(":80", nil))

	virtualizersdk.Macro(virtualizersdk.FALCON_TINY_END)
}
