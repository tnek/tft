package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/tnek/tft/riot"
)

const (
	rapiKeyPath = "./apikey"
)

func readRAPIKey(path string) (string, error) {
	dat, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(dat)), nil
}

func main() {
	rapiKey, err := readRAPIKey(rapiKeyPath)
	if err != nil {
		log.Fatalf("failed to read rapiKey from %v: %v", rapiKeyPath, err)
	}

	rc := riot.New(rapiKey, riot.DevLimiter)
	ctx := context.Background()
	summoner, err := rc.SummonerByName(ctx, "na1", "tnekk")
	if err != nil {
		log.Fatalf("%v", err)
	}

	fmt.Println(summoner)
}
