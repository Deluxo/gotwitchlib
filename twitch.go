package twitch

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
)

const (
	protocol  = "https://"
	host      = "api.twitch.tv"
	dir       = "/kraken"
	url       = protocol + host + dir
	TwitchUrl = protocol + "twitch.tv/"
)

type Streams struct {
	Links struct {
		Featured string
		Followed string
		Next     string
		Self     string
		Summary  string
	}
	Total   float64
	Streams []struct {
		ID    float64
		Links struct {
			Self string
		}
		AverageFps float64
		Channel    Channel
		CreatedAt  string
		Delay      float64
		Game       string
		IsPlaylist bool
		Preview    struct {
			Large    string
			Medium   string
			Small    string
			Template string
		}
		VideoHeight float64
		Viewers     float64
	}
}

type Game struct {
	ID    float64
	Links struct{}
	Box   struct {
		Large    string
		Medium   string
		Small    string
		Template string
	}
	GiantbombID float64
	Logo        struct {
		Large    string
		Medium   string
		Small    string
		Template string
	}
	Name string
}

type TopGames struct {
	Links struct {
		Next string
		Self string
	}
	Total float64
	Top   []struct {
		Channels float64
		Game     Game
		Viewers  float64
	}
}

type Channel struct {
	_id    int
	_links struct {
		Chat          string
		Commercial    string
		Editors       string
		Features      string
		Follows       string
		Self          string
		StreamKey     string
		Subscriptions string
		Teams         string
		Videos        string
	}
	Background                   interface{}
	Banner                       interface{}
	BroadcasterLanguage          string
	CreatedAt                    string
	Delay                        interface{}
	DisplayName                  string
	Followers                    int
	Game                         string
	Language                     string
	Logo                         string
	Mature                       bool
	Name                         string
	Partner                      bool
	ProfileBanner                string
	ProfileBannerBackgroundColor string
	Status                       string
	UpdatedAt                    string
	URL                          string
	VideoBanner                  string
	Views                        int
}

type onlineSubs struct {
	_links struct {
		Next string
		Self string
	}
	_total  int
	Streams []struct {
		_id    int
		_links struct {
			Self string
		}
		Channel    Channel
		AverageFps float64
		CreatedAt  string
		Delay      int
		Game       string
		IsPlaylist bool
		Preview    struct {
			Large    string
			Medium   string
			Small    string
			Template string
		}
		VideoHeight int
		Viewers     int
	}
}

type followingChannels struct {
	_links struct {
		Next string
		Self string
	}
	_total  int
	Follows []struct {
		_links struct {
			Self string
		}
		Channel       Channel
		CreatedAt     string
		Notifications bool
	}
}

func query(url string) []byte {
	res, _ := http.Get(url)
	ret, _ := ioutil.ReadAll(res.Body)
	res.Body.Close()
	return ret
}

func GetFollowedChannels(username string) followingChannels {
	var output followingChannels
	u := url + "/users/{username}/follows/channels?direction=DESC&sortby=created_at"
	u = strings.Replace(u, "{username}", username, 1)
	json.Unmarshal(query(u), &output)
	return output
}

func GetLiveSubs(oauthToken string) onlineSubs {
	var output onlineSubs
	u := url + "/streams/followed?oauth_token=" + oauthToken + "&stream_type=live"
	json.Unmarshal(query(u), &output)
	return output
}

func GetStreams(game, streamType string, limit, offset int) Streams {
	u := url + "/streams?"
	if game != "" {
		u += "game=" + game
	}
	if streamType != "" {
		u += "&stream_type=" + streamType
	}
	if limit != 0 {
		u += "&limit=" + strconv.Itoa(limit)
	} else {
		u += "&limit=10"
	}
	if offset != 0 {
		u += "&offset=" + strconv.Itoa(offset)
	}
	var output Streams
	json.Unmarshal(query(u), &output)
	return output
}

func GetTopGames(limit, offset *int) TopGames {
	var output TopGames
	u := url + "/games/top?"
	if *limit == 0 {
		u += "limit=10"
	} else {
		if *limit > 100 {
			*limit = 100
		}
		u += "limit=" + strconv.Itoa(*limit)
	}
	if *offset != 0 {
		u += "&offset=" + strconv.Itoa(*offset)
	}
	json.Unmarshal(query(u), &output)
	return output
}
