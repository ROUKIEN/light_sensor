package main

import (
	"bufio"
	"github.com/tarm/serial"
	"log"
	"strconv"
)

type Payload struct {
	Light int `json:"light"`
}

func main() {
	c := &serial.Config{Name: "/dev/ttyACM0", Baud: 9600}
	s, err := serial.OpenPort(c)

	if err != nil {
		log.Fatalf("serial.OpenPort: %v", err)
	}

	scanner := bufio.NewScanner(s)
	for scanner.Scan() {
		data := new(Payload)
		i, err := strconv.Atoi(scanner.Text())
		if err != nil {
			log.Printf("%v", err)
			continue
		}
		data.Light = i
		log.Printf("%v\n", data.Light)
	}
}
