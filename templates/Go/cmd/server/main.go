// package main

import (
	"fmt"
	"os"
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	}
}

func run() error {
	service := service.NewService()
	server := api.NewServer(service)
	r := server.Routes(r)

	log.Println("serving application at port :8080")
	err := http.ListenAndServe(":8080", r)
	return err
}
