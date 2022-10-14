package main

import (
	"context"
	"log"
	"os"
	"path/filepath"

	"github.com/galihsatriawan/eod/internal"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	path, err := os.Getwd()
	if err != nil {
		log.Fatalln(err)
	}
	pathData := filepath.Join(path, "data")
	repo := internal.NewRepository(pathData)
	svc := internal.NewService(repo)
	log.Println("Program processing....")
	err = svc.EndOfDay(ctx)
	if err != nil {
		log.Println("Processing failed")
		log.Fatalln(err)
	}
	log.Println("Processing success")
}
