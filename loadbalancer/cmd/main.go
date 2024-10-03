package main

import (
	"gopkg.in/yaml.v2"
	configuration2 "load-balancer/configuration"
	"log/slog"
	"net"
	"os"
	"sync"
)

var configuration configuration2.Configuration
var wg = sync.WaitGroup{}

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

// handleConnection - go routine to handle a connection
func handleConnection(connection net.Conn) {
	defer connection.Close()

	slog.Info("Handling connection.", "address", connection.RemoteAddr().String())

	// Read data from the connection
	//
	buffer := make([]byte, 1024)
	_, err := connection.Read(buffer)
	if err != nil {
		slog.Error("Failed to read data from connection.", "address", connection.RemoteAddr().String(), "error", err.Error())
		return
	}

	// Respond to the connection with http 1.1 200 OK (Hello world)
	//
	_, err = connection.Write([]byte("HTTP/1.1 200 OK\r\nContent-Length: 11\r\n\r\nHello world"))
	if err != nil {
		slog.Error("Failed to write data to connection.", "address", connection.RemoteAddr().String(), "error", err.Error())
		return
	}
}

// listen - go routine to listen on a listener specified in the configuration
func listen(listener configuration2.Listener) {
	slog.Info("Starting listener.", "name", listener.Name, "address", listener.Socket.Address, "port", listener.Socket.Port, "protocol", listener.Socket.Protocol)

	socket, err := net.Listen(listener.Socket.Protocol, listener.Socket.Address+":"+listener.Socket.Port)
	if err != nil {
		slog.Error("Failed to start listener.", "name", listener.Name, "error", err.Error())
		return
	}

	slog.Info("Listener started.", "name", listener.Name, "address", listener.Socket.Address, "port", listener.Socket.Port, "protocol", listener.Socket.Protocol)

	// Wait for incoming connections
	//
	for {
		connection, err := socket.Accept()
		if err != nil {
			slog.Error("Failed to accept connection.", "name", listener.Name, "error", err.Error())
			continue
		}

		wg.Add(1)
		go handleConnection(connection)
	}

	wg.Done()
}

func main() {

	// Start listening on the defined listeners
	//
	for _, listener := range configuration.Listeners {
		wg.Add(1)
		go listen(listener)
	}

	// Wait for all listeners to finish
	//
	wg.Wait()
}
