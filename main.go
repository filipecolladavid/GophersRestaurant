/*
Program to simulate the process of a (concurrent) working restaurant.

Usage:

	go build
	go run main.go
*/
package main

import (
	"pclt/restaurant"
	"time"
)

func main() {
	openTimeout := time.Duration(20 * time.Second)
	restaurant.OpenRestaurant(openTimeout)
}
