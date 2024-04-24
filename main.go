package main

import (
	"fmt"
	"log"
	"sync"
)

// main function

func main() {

	myConfig := loadYamlData()

	// print the config data

	// print the stations data

	//for stationType, cfg := range myConfig.Stations {
	//
	//	fmt.Printf("Station type: %s\n", stationType)
	//	fmt.Printf("Count: %d\n", cfg.Count)
	//	fmt.Printf("Serve time min: %s\n", cfg.ServeTimeMin)
	//	fmt.Printf("Serve time max: %s\n", cfg.ServeTimeMax)
	//}
	//
	//// print the registers data
	//
	//fmt.Printf("Registers count: %d\n", myConfig.Registers.Count)
	//fmt.Printf("Handle time min: %s\n", myConfig.Registers.HandleTimeMin)
	//fmt.Printf("Handle time max: %s\n", myConfig.Registers.HandleTimeMax)
	//
	//// print the car data
	//
	//fmt.Printf("Car count: %d\n", myConfig.Cars.Count)
	//fmt.Printf("Arrival time min: %s\n", myConfig.Cars.ArrivalTimeMin)
	//fmt.Printf("Arrival time max: %s\n", myConfig.Cars.ArrivalTimeMax)
	//

	// INITIALIZATION PART
	//queueLimitsEntry := 10
	///queueLimitsExit := 10

	// initialize cars and registers

	mi, err := ParseDuration(myConfig.Registers.HandleTimeMin)
	if err != nil {
		log.Fatalf("Error parsing handle time min for registers: %v", err)
	}

	ma, err := ParseDuration(myConfig.Registers.HandleTimeMax)
	if err != nil {
		log.Fatalf("Error parsing handle time max for registers: %v", err)
	}

	registers := make([]*Register, 0, myConfig.Registers.Count)
	for i := 0; i < myConfig.Registers.Count; i++ {
		register := InitializeRegister(mi, ma)
		registers = append(registers, register)
	}
	// initialize stations

	stations := make(map[string]*Station)

	// initialize pumps
	for stationType, cfg := range myConfig.Stations {

		serveTimeMin, err := ParseDuration(cfg.ServeTimeMin)
		if err != nil {
			log.Fatalf("Error parsing serve time min for %s: %v", stationType, err)
		}

		serveTimeMax, err := ParseDuration(cfg.ServeTimeMax)
		if err != nil {
			log.Fatalf("Error parsing serve time max for %s: %v", stationType, err)
		}

		stations[stationType] = InitializeStation(stationType, serveTimeMin, serveTimeMax)
	}

	// print the stations data

	for stationType, station := range stations {
		fmt.Printf("Station type: %s\n", stationType)
		fmt.Printf("Serve time min: %s\n", station.ServeTimeMin)
		fmt.Printf("Serve time max: %s\n", station.ServeTimeMax)
	}

	var wg sync.WaitGroup
	quit := make(chan struct{})

	wg.Add(1)
	go spawnCars(&stations, int(myConfig.Cars.Count), myConfig, &wg)

	// process the cars at each pump
	for _, station := range stations {
		wg.Add(1)
		go processPump(station, &registers, &wg, quit)
	}

	// process the cars at each pump
	for _, reg := range registers {
		wg.Add(1)
		go processRegister(reg, &wg, quit)
	}

	wg.Add(1)

	go monitorRegisters(registers, int(myConfig.Cars.Count), quit, &wg)

	wg.Wait()

	// print the statistics
	fmt.Printf("Station Statistics\n")
	for stationType, station := range stations {
		fmt.Printf("Station type: %s\n", stationType)
		fmt.Printf("Total cars: %d\n", station.TotalCars)
		fmt.Printf("Total time: %s\n", station.TotalTime)
		fmt.Printf("Max queue time: %s\n", station.MaxQueueTime)
		fmt.Printf("-----------------------\n")
	}
}
