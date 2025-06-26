package main

import (
	"fmt"
	"log"
	"math"
	"net/http"
	"os"
	"strconv"
)

// Mass represents an object that holds the density of a material
type Mass struct {
	Density float64
}

type MassVolume interface {
	density() float64
	volume(dimension float64) float64
}

type Sphere struct {
	Mass
}

func (s Sphere) density() float64 {
	return s.Density
}

func (s Sphere) volume(d float64) float64 {
	return (4.0 / 3.0) * math.Pi * math.Pow(d/2.0, 3)
}

type Cube struct {
	Mass
}

func (c Cube) density() float64 {
	return c.Density
}

func (c Cube) volume(d float64) float64 {
	return math.Pow(d, 3)
}

// Handler handles volume calculation requests with logging
func Handler(massVolume MassVolume) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Printf("Received request: %s %s", r.Method, r.URL.Path)

		dimensionStr := r.URL.Query().Get("dimension")
		if dimensionStr == "" {
			log.Printf("Missing 'dimension' parameter")
			http.Error(w, "'dimension' query param is required", http.StatusBadRequest)
			return
		}

		dimension, err := strconv.ParseFloat(dimensionStr, 64)
		if err != nil {
			log.Printf("Invalid 'dimension' parameter: %v", err)
			http.Error(w, "Invalid 'dimension' parameter", http.StatusBadRequest)
			return
		}

		weight := massVolume.density() * massVolume.volume(dimension)
		response := fmt.Sprintf("%.2f", math.Round(weight*100)/100)
		w.Write([]byte(response))
		log.Printf("Response: %s", response)
	}
}

// Health endpoints
func livenessHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("ok"))
}

func readinessHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("ready"))
}

func main() {
	port, err := strconv.Atoi(os.Args[1])
	if err != nil {
		log.Fatalf("Invalid port argument: %v", err)
	}

	aluminiumSphere := Sphere{Mass{Density: 2.710}} // g/cm³
	ironCube := Cube{Mass{Density: 7.874}}          // g/cm³

	// Business logic endpoints
	http.HandleFunc("/aluminium/sphere", Handler(aluminiumSphere))
	http.HandleFunc("/iron/cube", Handler(ironCube))

	// Health probe endpoints
	http.HandleFunc("/healthz", livenessHandler)
	http.HandleFunc("/readyz", readinessHandler)

	log.Printf("Starting server on port %d...", port)
	if err := http.ListenAndServe(fmt.Sprintf(":%d", port), nil); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}
