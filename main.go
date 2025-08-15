package main

import (
	"context"
	"fmt"
	"os"
)

func main() {
	fmt.Println("Hello, Leaning Server DB!")
	conn, err = pgx.Connect(context.Background(), os.Getenv("DATABASE_URL"))
	if err != nil {
		fmt.Println()
		os.Exit(1)
	}
	defer conn.Close(context.Background())
	fmt.Println("Connected to the database successfully!")

}
