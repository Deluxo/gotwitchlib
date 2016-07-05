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

type onlineStreams struct {
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

func GetOnlineStreams(oauthToken string) onlineStreams {
	var output onlineStreams
	u := url + "/streams/followed?oauth_token={oauthToken}&stream_type=live"
	u = strings.Replace(u, "{oauthToken}", oauthToken, 1)
	json.Unmarshal(query(u), &output)
	return output
}

func GetTopGames(limit, offset *int) TopGames {
	var output TopGames
	u := url + "/games/top"
	if *limit == 0 {
		u += "?limit=10"
	} else {
		if *limit > 100 {
			*limit = 100
		}
		u += "?limit=" + strconv.Itoa(*limit)
	}
	if *offset != 0 {
		u += "&offset=" + strconv.Itoa(*offset)
	}
	json.Unmarshal(query(u), &output)
	return output
}
