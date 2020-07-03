package riot

// Summoner represents a TFT account.
type Summoner struct {
	Name string `json:"name"`

	// RevisionDate is the date summoner was last modified specified as epoch
	// milliseconds. The following events will update this timestamp: summoner
	// name change, summoner level change, or profile icon change.
	RevisionDate int64 `json:"revisionDate"`

	// SummonerLevel is the summoner level associated with the summoner.
	SummonerLevel int64 `json:"summonerLevel"`

	// ProfileIconID is the ID of the summoner icon associated with the summoner.
	ProfileIconID int `json:"profileIconId"`

	// Riot API's 30 billion ways to refer to an account.
	// AccountID is the encrypted account ID. Max length 56 characters
	AccountID string `json:"accountId"`
	// Id is the encrypted summoner ID. Max length 63 characters.
	ID string `json:"id"`
	// PUUID is the encrypted PUUID. Max length of 78 characters.
	PUUID string `json:"puuid"`

	// RAPI endpoint types to use (https://developer.riotgames.com/docs/tft#_routing-values).
	Region   string
	Platform string
}

// LeagueEntryDTO represents a summoner's League ranking
type LeagueEntryDTO struct {
	LeagueID     string `json:"leagueId"`
	SummonerID   string `json:"summonerId"`
	SummonerName string `json:"summonerName"`

	QueueType    string `json:"queueType"`
	Tier         string `json:"tier"`
	Rank         string `json:"rank"`
	LeaguePoints int    `json:"leaguePoints"`

	// First placement.
	Wins int `json:"wins"`
	// Second through eigth placement
	Losses int `json:"losses"`

	// Interesting internal user account flags
	HotStreak  bool `json:"hotStreak"`
	Veteran    bool `json:"veteran"`
	FreshBlood bool `json:"freshBlood"`
	Inactive   bool `json:"inactive"`

	// miniSeries of type MiniSeriesDTO is another field mentioned in the API Docs,
	// but seems to be unused and a leftover from when the TFT match API was just
	// part of the main	League API.
}
