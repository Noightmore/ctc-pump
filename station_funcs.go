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
				//fmt.Println("Car added to station queue")
				added = true
				break
			}
		}

		if !added {
			//fmt.Println("All station queues are full, waiting...")
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

			// Initialize variables for processing
			placed := false
			t := time.Duration(rand.Int63n(int64(station.ServeTimeMax)-int64(station.ServeTimeMin)) + int64(station.ServeTimeMin))

			// Attempt to place car in an available register
		processLoop: // Label for the outer loop
			for _, register := range *registers {
				select {
				case register.RegisterQueue <- car:
					// Successful placement
					placed = true
					//fmt.Println("Car processed at pump and placed in register queue in", t)
					break processLoop
				case <-quit:
					return // Exit if quit signal is received
				default:
					// This register is full, try the next one
				}
			}

			if !placed {
				//fmt.Println("All register queues are full, waiting for an available register...")
				// Wait until we can place the car or a quit signal is received
				placed = waitForAvailableRegister(registers, car, quit)
				if !placed {
					//fmt.Println("Failed to place car, quitting...")
					return // Optionally handle as needed if the car could not be placed after waiting
				}
			}

			// Update station data only after successful placement
			station.TotalCars++
			station.TotalTime += t
			if station.MaxQueueTime < t {
				station.MaxQueueTime = t
			}

			time.Sleep(t) // Simulate the time taken to process the car at the pump

		case <-quit:
			return // Exit immediately if quit signal is received
		}
	}
}

// Helper function to block until a car can be placed in a register or quit is received
func waitForAvailableRegister(registers *[]*Register, car *Car, quit chan struct{}) bool {
	for {
		for _, register := range *registers {
			select {
			case register.RegisterQueue <- car:
				//fmt.Println("Car placed in register queue after waiting.")
				return true
			case <-quit:
				return false
			default:
				// Continue to next register

			}
		}
		// Optional: small sleep to prevent a tight loop hammering the CPU
		time.Sleep(100 * time.Millisecond)
	}
}

func processRegister(register *Register, wg *sync.WaitGroup, quit chan struct{}) {
	defer wg.Done()
	for {
		select {
		case _, ok := <-register.RegisterQueue:
			if !ok {
				return // Exit if the RegisterQueue is closed
			}
			// Process the car here
			t := time.Duration(rand.Int63n(int64(register.HandleTimeMax)-int64(register.HandleTimeMin)) + int64(register.HandleTimeMin))

			register.Mutex.Lock()
			register.TotalCars++
			register.TotalTime += t

			if register.MaxQueueTime < t {
				register.MaxQueueTime = t
			}
			register.Mutex.Unlock()
			// sleep for t time
			time.Sleep(t)
			// print the car

			//fmt.Printf("Car: %v\n", car)
			//fmt.Printf("Car processed at register in %v\n", t)

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
