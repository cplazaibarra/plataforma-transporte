package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"
)

type TelemetryPayload struct {
	VehicleID    string  `json:"vehicle_id"`
	DriverID     string  `json:"driver_id"`
	Latitude     float64 `json:"latitude"`
	Longitude    float64 `json:"longitude"`
	Speed        float64 `json:"speed"`
	Heading      float64 `json:"heading"`
	BatteryLevel int     `json:"battery_level"`
}

func main() {
	vehicleID := "a1000000-0000-0000-0000-000000000001"
	driverID := "d1000000-0000-0000-0000-000000000001"
	apiURL := "http://localhost:8095/api/v1/telemetry"

	// Coordenadas simuladas en ruta desde Santiago Centro hacia Providencia
	waypoints := [][2]float64{
		{-33.4429, -70.6539}, // Santiago Centro
		{-33.4410, -70.6480},
		{-33.4385, -70.6410},
		{-33.4360, -70.6350},
		{-33.4330, -70.6290},
		{-33.4310, -70.6250},
		{-33.4295, -70.6210},
		{-33.4285, -70.6184}, // Destino: Av Providencia 1240
	}

	log.Println("🚗 Iniciando simulador de telemetría GPS para Móvil KJ-DF-44...")

	for {
		for i, pt := range waypoints {
			payload := TelemetryPayload{
				VehicleID:    vehicleID,
				DriverID:     driverID,
				Latitude:     pt[0],
				Longitude:    pt[1],
				Speed:        35.5 + float64(i%5)*2.5,
				Heading:      68.0,
				BatteryLevel: 88 - (i % 5),
			}

			body, _ := json.Marshal(payload)
			resp, err := http.Post(apiURL, "application/json", bytes.NewBuffer(body))
			if err != nil {
				log.Printf("⚠️ Error enviando telemetría: %v. ¿El backend está corriendo?", err)
			} else {
				resp.Body.Close()
				fmt.Printf("📍 Telemetría enviada [%d/%d]: Lat: %.5f, Lng: %.5f, Vel: %.1f km/h\n",
					i+1, len(waypoints), pt[0], pt[1], payload.Speed)
			}

			time.Sleep(3 * time.Second)
		}
		log.Println("🔄 Reiniciando ciclo de ruta simulada...")
		time.Sleep(2 * time.Second)
	}
}
