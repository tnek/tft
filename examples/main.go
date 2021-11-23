package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"

	riotclient "github.com/tnek/tft/riot/client"
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

	rc := riotclient.New(rapiKey, riotclient.DevLimiter)
	ctx := context.Background()
	s, err := rc.SummonerByName(ctx, "na1", "tnekk")
	if err != nil {
		log.Fatalf("%v", err)
	}

	fmt.Println(s)
	id, err := rc.Matches(ctx, s, 1)
	if err != nil {
		log.Fatalf("%v", err)
	}
	match, err := rc.Match(ctx, s.Region, id[0])
	if err != nil {
		log.Fatalf("%v", err)
	}
	fmt.Println(match)
}
