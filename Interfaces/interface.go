package main

import "fmt"

type IVehicle interface {
	PrintMake()
	PrintSpeed()
}

type Vehicle struct {
	make string
	speed int
}

func (v Vehicle) PrintMake() {
	fmt.Println(v.make)
}

func (v Vehicle) PrintSpeed() {
	fmt.Println(v.speed)	
}

func main() {
	var v IVehicle
	v = Vehicle{"BMW", 155}
	v.PrintMake()
	v.PrintSpeed()
}

