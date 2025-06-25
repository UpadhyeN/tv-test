package main

import (
	"fmt"
	"math"
	"net/http"
	"os"
	"strconv"
)

// mass is object that holds the density of a material
type Mass struct {
	Density float64
}

// --- ADD YOUR CODE ---

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

// --- BETWEEN THOSE LINES ---

func Handler(massVolume MassVolume) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if dimension, err := strconv.ParseFloat(r.URL.Query().Get("dimension"), 64); err == nil {
			weight := massVolume.density() * massVolume.volume(dimension)
			w.Write([]byte(fmt.Sprintf("%.2f", math.Round(weight*100)/100)))
			return
		}
		w.WriteHeader(http.StatusBadRequest)
	}
}

func main() {
	port, err := strconv.Atoi(os.Args[1])
	if err != nil {
		panic(err)
	}

	aluminiumSphere := Sphere{Mass{Density: 2.710}} // g/cm³
	ironCube := Cube{Mass{Density: 7.874}}          // g/cm³

	http.HandleFunc("/aluminium/sphere", Handler(aluminiumSphere))
	http.HandleFunc("/iron/cube", Handler(ironCube))

	if err := http.ListenAndServe(fmt.Sprintf(":%d", port), nil); err != nil {
		panic(err)
	}
}
