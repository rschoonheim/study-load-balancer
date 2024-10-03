package main

import (
	"gopkg.in/yaml.v2"
	configuration2 "load-balancer/configuration"
	"load-balancer/linux"
	"log/slog"
	"net"
	"os"
	"sync"
)

var configuration configuration2.Configuration
var waitGroup sync.WaitGroup = sync.WaitGroup{}

func init() {

	// Check if `MAIN_CONFIGURATION` is set
	//
	mainConfigurationPath := os.Getenv("MAIN_CONFIGURATION")
	if mainConfigurationPath == "" {
		println("MAIN_CONFIGURATION environment variable is not set")
		os.Exit(1)
	}

	// Check if the main configuration file exists
	//
	_, err := os.Stat(mainConfigurationPath)
	if os.IsNotExist(err) {
		println("Main configuration file does not exist")
		os.Exit(1)
	}

	// Read contents of the main configuration file
	//
	mainConfigurationFile, err := os.Open(mainConfigurationPath)
	if err != nil {
		println(err.Error())
		os.Exit(1)
	}

	// Parse the main configuration file
	//
	if err := yaml.NewDecoder(mainConfigurationFile).Decode(&configuration); err != nil {
		println(err.Error())
		os.Exit(1)
	}
}

func socketListenerStart(listener configuration2.Listener) error {

	// Find address of the interface
	//
	var address string
	for _, networkInterface := range configuration.Interfaces {
		if networkInterface.Name == listener.Interface {
			address = networkInterface.Ip
			break
		}
	}

	if address == "" {
		slog.Error("Interface not found.", "name", listener.Interface)
		return nil
	}

	slog.Info("Starting socket listener.", "name", listener.Name, "address", address, "port", listener.Socket.Port, "protocol", listener.Socket.Protocol)

	// Start the socket listener
	//
	listner, err := net.Listen(listener.Socket.Protocol, address+":"+listener.Socket.Port)
	if err != nil {
		return err
	}

	slog.Info("Socket listener started.", "name", listener.Name, "address", address, "port", listener.Socket.Port, "protocol", listener.Socket.Protocol)

	// Accept incoming connections
	//
	for {
		conn, err := listner.Accept()
		if err != nil {
			return err
		}

		// Handle the connection
		//
		go func(conn net.Conn) {
			defer conn.Close()

			// Read the request
			//
			buffer := make([]byte, 1024)
			_, err := conn.Read(buffer)
			if err != nil {
				return
			}

			// Send the response
			//
			_, err = conn.Write([]byte("Hello, World!"))
			if err != nil {
				return
			}
		}(conn)
	}
}

func main() {

	// Create network interfaces for the load balancer
	//
	for _, networkInterface := range configuration.Interfaces {
		err := linux.NetworkInterfaceCreate(networkInterface.Name, networkInterface.Ip)
		if err != nil {
			println(err.Error())
			os.Exit(1)
		}

		slog.Info("Network interface created.", "name", networkInterface.Name, "ip", networkInterface.Ip)
	}

	// Start socket listeners for the load balancer
	//
	for _, listener := range configuration.Listeners {
		waitGroup.Add(1)

		go func(listener configuration2.Listener) {
			defer waitGroup.Done()

			err := socketListenerStart(listener)
			if err != nil {
				println(err.Error())
				os.Exit(1)
			}
		}(listener)
	}

	// Wait for the socket listeners to finish
	//
	waitGroup.Wait()
}
