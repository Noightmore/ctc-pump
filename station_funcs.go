package main

import (
	"fmt"
	"log"
	"math/rand"
	"sync"
	"time"
)

func spawnCars(stations *map[string]*Station, carCount int, carCfg *Config, wg *sync.WaitGroup) {
	defer wg.Done()

	for i := 0; i < carCount; i++ {
		// Parse car configuration to time.Duration
		arrivalTimeMin, err := time.ParseDuration(carCfg.Cars.ArrivalTimeMin)
		if err != nil {
			log.Fatalf("Error parsing arrival time min for cars: %v", err)
		}

		arrivalTimeMax, err := time.ParseDuration(carCfg.Cars.ArrivalTimeMax)
		if err != nil {
			log.Fatalf("Error parsing arrival time max for cars: %v", err)
		}

		// Create a car instance
		car := InitializeCar(arrivalTimeMin, arrivalTimeMax)

		// Wait a random duration within specified min and max before trying to add a car
		waitDuration := arrivalTimeMin + time.Duration(rand.Int63n(int64(arrivalTimeMax-arrivalTimeMin)))
		time.Sleep(waitDuration)

		// Attempt to add the car to the first available station's queue
		added := false
		for _, station := range *stations {
			if len(station.StationQueue) < cap(station.StationQueue) && station.IsAvailable {
				station.StationQueue <- car
				fmt.Println("Car added to station queue")
				added = true
				break
			}
		}

		if !added {
			fmt.Println("All station queues are full, waiting...")
			time.Sleep(10 * time.Millisecond) // Wait before trying again if no queue was available
		}
	}
}

func processPump(station *Station, registers *[]*Register, wg *sync.WaitGroup, quit chan struct{}) {
	defer wg.Done()
	for {
		select {
		case car, ok := <-station.StationQueue:
			if !ok {
				return // Exit if the StationQueue is closed
			}
			// Process the car here
			placed := false

		processLoop: // Label for the outer loop
			for _, register := range *registers {
				select {
				case register.RegisterQueue <- car:
					// Handle the car once successfully placed
					t := time.Duration(rand.Int63n(int64(station.ServeTimeMax)-int64(station.ServeTimeMin)) + int64(station.ServeTimeMin))
					time.Sleep(t) // Simulate processing time

					// Update station data
					station.TotalCars++
					station.TotalTime += t

					placed = true
					break processLoop // Break out of the for-loop, not just the select
				case <-quit:
					return // Exit immediately if quit signal is received
				default:
					// Print message and continue to the next register
					fmt.Println("Register queue is full. Trying next...")
				}
			}

			if !placed {
				fmt.Println("All register queues are full. Retrying...")
				time.Sleep(1 * time.Second) // Wait and retry or exit based on specific needs
			}
		case <-quit:
			return // Exit immediately if quit signal is received
		}
	}
}

func processRegister(register *Register, wg *sync.WaitGroup, quit chan struct{}) {
	defer wg.Done()
	for {
		select {
		case car, ok := <-register.RegisterQueue:
			if !ok {
				return // Exit if the RegisterQueue is closed
			}
			// Process the car here
			t := time.Duration(rand.Int63n(int64(register.HandleTimeMax)-int64(register.HandleTimeMin)) + int64(register.HandleTimeMin))

			register.Mutex.Lock()
			register.TotalCars++
			register.TotalTime += t
			register.Mutex.Unlock()
			// sleep for t time
			time.Sleep(t)
			// print the car
			fmt.Printf("Car: %v\n", car)

		case <-quit:
			return // Exit immediately if quit signal is received
		}
	}
}

// Monitor all registers and coordinate shutdown if max car count is reached
func monitorRegisters(registers []*Register, maxCarCount int, quit chan struct{}, wg *sync.WaitGroup) {
	defer wg.Done()

	for {
		// sum total cars processed by all registers
		totalCars := 0
		for _, register := range registers {

			register.Mutex.Lock()
			totalCars += register.TotalCars
			register.Mutex.Unlock()

		}

		// check if max car count is reached
		if totalCars >= maxCarCount {
			// close all go routines
			fmt.Println("Max car count reached. Closing all go routines...")
			close(quit)
			break
		}

		// sleep for a while before checking again
		time.Sleep(1 * time.Second)
	}
}
