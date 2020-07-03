package tft

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"path"
)

// Client represents a handler for the Riot API
type Client struct {
	// Key is your RAPI key.
	Key string
}

// get is a wrapper for fetching from RAPI.
func (c *Client) get(ctx context.Context, routing string, endpoint string, obj interface{}) error {
	ep := "https://" + path.Join(fmt.Sprintf("%s.api.riotgames.com", routing), endpoint)

	req, err := http.NewRequest("GET", ep, nil)
	if err != nil {
		return fmt.Errorf("client.Get(%q): failed on initializing http.NewRequest %v", ep, err)
	}

	req.Header.Set("X-Riot-Token", c.Key)
	hc := &http.Client{}
	resp, err := hc.Do(req)
	if err != nil {
		return fmt.Errorf("client.Get(%q): failed on request: %v", ep, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf(`client.Get(%q): request failed with status "%d %s"`, ep, resp.StatusCode, resp.Status)
	}

	if err := json.NewDecoder(resp.Body).Decode(obj); err != nil {
		return fmt.Errorf("client.Get(%q): failed on json parse of request body: %v", ep, err)
	}

	return nil
}
