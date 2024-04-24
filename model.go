package main

import (
	"math/rand"
	"sync"
	"time"
)

type Car struct {
	ArrivalTime              time.Time
	PumpQueueEntryTime       time.Time
	PumpServiceEntryTime     time.Time
	RegisterQueueEntryTime   time.Time
	RegisterServiceEntryTime time.Time
}

func InitializeCar(ArrivalTimeMin time.Duration, ArrivalTimeMax time.Duration) *Car {

	ArrivalTime :=
		time.Now().Add(time.Duration(rand.Intn(int(ArrivalTimeMax)-int(ArrivalTimeMin)+1)) + ArrivalTimeMin)

	car := Car{}
	car.ArrivalTime = ArrivalTime
	car.PumpQueueEntryTime = time.Time{}
	car.RegisterQueueEntryTime = time.Time{}
	car.PumpServiceEntryTime = time.Time{}
	car.RegisterServiceEntryTime = time.Time{}

	return &car
}

type Register struct {
	HandleTimeMin time.Duration
	HandleTimeMax time.Duration
	TotalCars     int
	TotalTime     time.Duration
	MaxQueueTime  time.Duration
	RegisterQueue chan *Car
	Mutex         *sync.Mutex
}

func InitializeRegister(HandleTimeMin time.Duration, HandleTimeMax time.Duration) *Register {

	register := Register{}

	register.HandleTimeMin = HandleTimeMin
	register.HandleTimeMax = HandleTimeMax
	register.TotalCars = 0
	register.TotalTime = 0
	register.MaxQueueTime = 0
	register.RegisterQueue = make(chan *Car, 3)
	register.Mutex = &sync.Mutex{}

	return &register

}

type Station struct {
	StationType    string
	ServeTimeMin   time.Duration
	ServeTimeMax   time.Duration
	IsAvailable    bool
	TotalCars      int
	TotalTime      time.Duration
	TotalQueueTime time.Duration
	MaxQueueTime   time.Duration
	StationQueue   chan *Car
}

func InitializeStation(StationType string, ServeTimeMin time.Duration, ServeTimeMax time.Duration) *Station {
	station := Station{}

	station.StationType = StationType
	station.ServeTimeMin = ServeTimeMin
	station.ServeTimeMax = ServeTimeMax
	station.IsAvailable = true
	station.TotalCars = 0
	station.TotalTime = 0
	station.TotalQueueTime = 0
	station.MaxQueueTime = 0
	station.StationQueue = make(chan *Car, 10)

	return &station

}
