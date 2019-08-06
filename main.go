package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"strings"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"gopkg.in/alecthomas/kingpin.v2"
)

var (
	brokerURL = kingpin.Flag("broker", "MQTT Server to connect to").Default("tcp://localhost:1883").String()
	clientID  = kingpin.Flag("client_id", "MQTT Client ID").Default("serial-mqtt-client").String()
	port      = kingpin.Flag("port", "Serial Port").Default("/dev/ttyUSB0").String()
	baudrate  = kingpin.Flag("baudrate", "Baudrate").Default("115200").Int()

	client     mqtt.Client
	serialPort io.ReadWriteCloser
)

func mainWithErrors() error {
	kingpin.Parse()

	opt := mqtt.NewClientOptions()
	opt = opt.AddBroker(*brokerURL)
	opt = opt.SetClientID(*clientID)

	client = mqtt.NewClient(opt)

	token := client.Connect()
	token.Wait()
	if err := token.Error(); err != nil {
		return fmt.Errorf("Unable to connect to broker: %v", err)
	}
	log.Printf("Connected to %v", *brokerURL)
	defer client.Disconnect(250)

	var err error
	serialPort, err = initSerial(*port, *baudrate)
	if err != nil {
		return fmt.Errorf("Unable to connect to serial connection: %v", err)
	}
	log.Printf("Connected to %v", *brokerURL)

	scanner := bufio.NewScanner(serialPort)
	for scanner.Scan() {
		lineSplit := strings.SplitN(scanner.Text(), " ", 2)
		if len(lineSplit) != 3 {
			log.Println("Nonesense line...continue")
			continue
		}
		switch command := lineSplit[0]; command {
		case "publish":
			topic := lineSplit[1]
			content := lineSplit[2]
			token := client.Publish(topic, 0, false, content)
			token.Wait()
			if err := token.Error(); err != nil {
				log.Printf("Unable to publish to %v: %v\n", topic, err)
			}
			break
		default:
			log.Println("Unknown Command")
		}
	}

	return nil
}

func main() {
	err := mainWithErrors()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
