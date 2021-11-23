package riot

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"path"

	"github.com/tnek/multilimiter"
	"golang.org/x/time/rate"
)

const (
	summonerAPIPrefix = "/tft/summoner/v1/summoners/"
	matchAPIPrefix    = "/tft/match/v1/matches"
	rankedAPIPrefix   = "/tft/league/v1/"
)

var (
	// DevLimiter is a rate limiter that follows the rate limits imposed on RAPI dev Keys.
	DevLimiter = multilimiter.New([]*rate.Limiter{
		rate.NewLimiter(1, 20),
		rate.NewLimiter(120, 100),
	})
)

type Client interface {
	// SummonerByName retrieves a Summoner object by username.
	SummonerByName(ctx context.Context, platform string, name string) (*Summoner, error)

	// League retrieves personal Ranked information about a given summoner
	League(ctx context.Context, s *Summoner) (*LeagueEntryDTO, error)

	// Matches returns the match IDs of the last matches of a given Summoner
	Matches(ctx context.Context, s *Summoner, count int) ([]string, error)

	// Match fetches data about a specific match given a Match ID.
	Match(ctx context.Context, region string, matchID string) (*Match, error)
}

// baseClient is the base wrapper for querying the Riot API.
// It does not perform any rate limiting or request caching.
type baseClient struct {
	Client
	// Key is your RAPI key.
	Key string

	limiter *multilimiter.Limiter
}

func New(key string, limiter *multilimiter.Limiter) *baseClient {
	c := &baseClient{Key: key}
	c.limiter = limiter
	return c
}

// get is a wrapper for directly fetching from RAPI.
func (c *baseClient) get(ctx context.Context, routing string, endpoint string, obj interface{}) error {
	ep := "https://" + path.Join(fmt.Sprintf("%s.api.riotgames.com", routing), endpoint)

	if err := c.limiter.Wait(ctx); err != nil {
		return fmt.Errorf("get(%q): failed on waiting for rate limit: %v", ep, err)
	}

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

func (c *baseClient) SummonerByName(ctx context.Context, platform string, name string) (*Summoner, error) {
	ep := path.Join(summonerAPIPrefix, "by-name", name)
	s := &Summoner{}

	if err := c.get(ctx, platform, ep, s); err != nil {
		return nil, err
	}
	s.Platform = platform
	s.Region = PlatformToRegion[platform]
	return s, nil
}

func (c *baseClient) League(ctx context.Context, s *Summoner) (*LeagueEntryDTO, error) {
	ep := path.Join(rankedAPIPrefix, "entries/by-summoner", s.ID)

	// The TFT League API returns a list of LeagueEntryDTOs despite the list always
	// only ever containing one entry, since it's copy-pasted from the regular
	// League Ranked API.
	var leagues []LeagueEntryDTO
	if err := c.get(ctx, s.Platform, ep, &leagues); err != nil {
		return nil, err
	}

	return &leagues[0], nil
}

func (c *baseClient) Matches(ctx context.Context, s *Summoner, count int) ([]string, error) {
	ep := path.Join(matchAPIPrefix, fmt.Sprintf("by-puuid/%s/ids?count=%d", s.PUUID, count))

	var ids []string
	if err := c.get(ctx, s.Region, ep, &ids); err != nil {
		return nil, err
	}

	return ids, nil
}

func (c *baseClient) Match(ctx context.Context, region string, matchID string) (*Match, error) {
	ep := path.Join(matchAPIPrefix, matchID)
	m := &Match{}
	if err := c.get(ctx, region, ep, &m); err != nil {
		return nil, err
	}

	return m, nil
}
