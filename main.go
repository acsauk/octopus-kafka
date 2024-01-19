package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/acsauk/octopus-kafka/internal/octopus"
	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatalf("Error loading .env file")
	}

	octopusAPIKey := os.Getenv("OCTOPUS_API_KEY")
	electricityMPAN := os.Getenv("ELECTRICITY_MPAN")

	httpClient := http.Client{}

	client := octopus.New(octopusAPIKey, "https://api.octopus.energy/v1", &httpClient)

	mp, err := client.ElectricityMeterPoints(electricityMPAN)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("GSP is %v\n", mp.GSP)
	fmt.Printf("MPAN is %v\n", mp.MPAN)
	fmt.Printf("ProfileClass is %v\n", mp.ProfileClass)

	account, err := client.Account("A-93DD6C62")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("%v\n", account)
}
