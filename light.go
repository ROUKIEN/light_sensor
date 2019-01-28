package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"github.com/eclipse/paho.mqtt.golang"
	"github.com/tarm/serial"
	"log"
	"os"
	"strconv"
	"time"
)

type Payload struct {
	Light int `json:"light"`
}

var f mqtt.MessageHandler = func(client mqtt.Client, msg mqtt.Message) {
	fmt.Printf("TOPIC: %s\n", msg.Topic())
	fmt.Printf("MSG: %s\n", msg.Payload())
}

func main() {
	mqtt.DEBUG = log.New(os.Stdout, "", 0)
	mqtt.ERROR = log.New(os.Stderr, "", 0)
	opts := mqtt.NewClientOptions().AddBroker("localhost:1883").SetClientID("gotrivial")
	opts.SetKeepAlive(2 * time.Second)
	opts.SetDefaultPublishHandler(f)
	opts.SetPingTimeout(1 * time.Second)

	client := mqtt.NewClient(opts)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}

	c := &serial.Config{Name: "/dev/ttyACM0", Baud: 9600}
	s, err := serial.OpenPort(c)

	if err != nil {
		log.Fatalf("serial.OpenPort: %v", err)
	}

	scanner := bufio.NewScanner(s)
	for scanner.Scan() {
		data := new(Payload)
		light, err := strconv.Atoi(scanner.Text())
		if err != nil {
			log.Printf("Error: %v", err)
			continue
		}

		data.Light = light
		jsonStr, _ := json.Marshal(data)
		token := client.Publish("test_channel", 0, false, string(jsonStr))
		token.Wait()
	}
}
