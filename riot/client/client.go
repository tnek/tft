package riotclient

import (
	"context"
	"fmt"

	"github.com/tnek/tft/riot"
)

const (
	summonerAPIPrefix = "/tft/summoner/v1/summoners/"
	matchAPIPrefix    = "/tft/match/v1/matches"
	rankedAPIPrefix   = "/tft/league/v1/"
)

type Client interface {
	// SummonerByName retrieves a Summoner object by username.
	SummonerByName(ctx context.Context, platform string, name string) (*riot.Summoner, error)

	// League retrieves personal Ranked information about a given summoner
	League(ctx context.Context, s *riot.Summoner) (*riot.LeagueEntryDTO, error)

	// Matches returns the match IDs of the last matches of a given Summoner
	Matches(ctx context.Context, s *riot.Summoner, count int) ([]string, error)

	// Match fetches data about a specific match given a Match ID.
	Match(ctx context.Context, region string, matchID string) (*riot.Match, error)
}

func MatchesInSet(ctx context.Context, c Client, s *riot.Summoner, set int) ([]*riot.Match, error) {
	league, err := c.League(ctx, s)
	if err != nil {
		return nil, fmt.Errorf("failed to get league info of %v: %v", s.Name, err)
	}
	totalGames := league.Wins + league.Losses
	ids, err := c.Matches(ctx, s, totalGames)
	if err != nil {
		return nil, fmt.Errorf("failed to get all past matches of %v: %v", s.Name, err)
	}
	var matches []*riot.Match
	idToErr := map[string]error{}
	for _, id := range ids {
		match, err := c.Match(ctx, s.Region, id)
		if err != nil {
			idToErr[id] = err
			continue
		}
		if match.Info.Set != set {
			break
		}
		matches = append(matches, match)
	}

	if len(idToErr) != 0 {
		return matches, fmt.Errorf("errors on fetching matches: %v", idToErr)
	}

	return matches, nil
}
