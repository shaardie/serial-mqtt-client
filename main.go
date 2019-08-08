package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"gopkg.in/alecthomas/kingpin.v2"

	parser "github.com/shaardie/serial-mqtt-client/parser"
)

var (
	brokerURL = kingpin.Flag("broker", "MQTT Server to connect to").Default("tcp://localhost:1883").String()
	clientID  = kingpin.Flag("client_id", "MQTT Client ID").Default("serial-mqtt-client").String()
	port      = kingpin.Flag("port", "Serial Port").Default("/dev/ttyUSB0").String()
	baudrate  = kingpin.Flag("baudrate", "Baudrate").Default("115200").Int()

	client mqtt.Client
	rwc    io.ReadWriteCloser
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
	rwc, err = NewSerial(*port, *baudrate)
	if err != nil {
		return fmt.Errorf("Unable to connect to connect: %v", err)
	}
	defer rwc.Close()
	log.Printf("Connected to %v", *port)

	scanner := bufio.NewScanner(rwc)
	for scanner.Scan() {
		line := scanner.Text()
		log.Println(line)
		cmd, err := parser.ParseLine(line)
		if err != nil {
			log.Println(err)
			continue
		}
		if cmd == nil {
			continue
		}
		switch cmd.Command {
		case parser.PUBLISH:
			log.Printf("Publish %v to %v", cmd.Value, cmd.Topic)
			token := client.Publish(cmd.Topic, 0, false, cmd.Value)
			token.Wait()
			if err := token.Error(); err != nil {
				log.Printf("Failure during publishing %v", err)
				continue
			}
			break
		case parser.SUBSCRIBE:
			break
		default:
			log.Printf("Unknown Command %v", cmd.Command)
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
