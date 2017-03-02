package twitch

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strconv"
)

// Full path to twitch.tv service
const (
	protocol  = "https://"
	host      = "api.twitch.tv"
	dir       = "/kraken"
	url       = protocol + host + dir
	TwitchURL = protocol + "twitch.tv/"
)

// Streams https://github.com/justintv/Twitch-API/blob/master/v3_resources/streams.md
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
// Game https://github.com/justintv/Twitch-API/blob/master/v3_resources/games.md
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
// TopGames same as games. Almost...
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
// Channel https://github.com/justintv/Twitch-API/blob/master/v3_resources/channels.md
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

type OnlineSubs struct {
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

//FollowingChannels the list of follows
type FollowingChannels struct {
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

// FollowChannel is the output by the methods
// that control the following of a streamer.
type FollowChannel struct {
	notifications bool
	Channel       Channel
}

func query(url, method string) []byte {
	client := &http.Client{}
	if method == "" {
		method = "GET"
	}
	request, _ := http.NewRequest(method, url, nil)
	res, _ := client.Do(request)
	ret, _ := ioutil.ReadAll(res.Body)
	defer res.Body.Close()
	return ret
}

// GetFollowedChannels simple getter
func GetFollowedChannels(oauthToken, username string) FollowingChannels {
	var output FollowingChannels
	u := url + "/users/" + username + "/follows/channels?direction=DESC&sortby=created_at"
	json.Unmarshal(query(u, ""), &output)
	return output
}

// GetLiveSubs gets the list of the streamers that are currently live
func GetLiveSubs(oauthToken string) OnlineSubs {
	var output OnlineSubs
	u := url + "/streams/followed?oauth_token=" + oauthToken + "&stream_type=live"
	json.Unmarshal(query(u, ""), &output)
	return output
}

// GetStreams simply retrieves live steams that are currently trending
func GetStreams(oauthToken, game, streamType string, limit, offset int) Streams {
	u := url + "/streams?oauth_token=" + oauthToken
	if game != "" {
		u += "&game=" + game
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
	json.Unmarshal(query(u, ""), &output)
	return output
}

// GetTopGames gets the list of games played the most rigth now
func GetTopGames(oauthToken string, limit, offset *int) TopGames {
	var output TopGames
	u := url + "/games/top?oauth_token=" + oauthToken
	if *limit == 0 {
		u += "&limit=10"
	} else {
		if *limit > 100 {
			*limit = 100
		}
		u += "&limit=" + strconv.Itoa(*limit)
	}
	if *offset != 0 {
		u += "&offset=" + strconv.Itoa(*offset)
	}
	json.Unmarshal(query(u, ""), &output)
	return output
}

func follow(oauthToken, username, channel, method string, notification bool) FollowChannel {
	var output FollowChannel
	u := url + "/users/" + username + "/follows/channels/" + channel + "?oauth_token=" + oauthToken
	if notification == true {
		u += "&notification"
	}
	json.Unmarshal(query(u, "PUT"), &output)
	return output
}

// Follow makes you to follow the streamer
func Follow(oauthToken, username, channel string, notification bool) FollowChannel {
	return follow(oauthToken, username, channel, "PUT", notification)
}

// Unfollow makes you to unfollow the streamer
func Unfollow(oauthToken, username, channel string) FollowChannel {
	return follow(oauthToken, username, channel, "DELETE", false)
}
