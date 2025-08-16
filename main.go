package main

import (
	"context"
	"fmt"
	"learningServerDB/internal/cfg"
	"learningServerDB/internal/server"
	"os"
	"os/signal"
)

func main() {
	// Load configuration
	cfg := cfg.LoadConfig()
	ctx, cancel := context.WithCancel(context.Background())

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	serverInstans := server.NewServer(ctx, &cfg)
	serverInstans.Serv()

	go func() {
		oscall := <-c
		fmt.Println("Received signal:", oscall)
		cancel()
		serverInstans.Shutdown()
		fmt.Println("Server shutdown gracefully")
	}()
	//fmt.Println("Hello, Leaning Server DB!")
	//conn, err := pgx.Connect(context.Background(), cfg.GetDBString())
	//if err != nil {
	//	fmt.Println("Failed to connect to the database: " + err.Error())
	//	os.Exit(1)
	//}
	//defer conn.Close(context.Background())
	//fmt.Println("Connected to the database successfully!")

}
