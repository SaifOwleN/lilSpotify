package spotApi

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type SpotApi struct {
	Config
}

type AlbumCover struct {
	Url []string `json:"url"`
}

type Artist struct {
	Name string `json:"name"`
}

type Album struct {
	Name  string       `json:"name"`
	Cover []AlbumCover `json:"images"`
}

type Track struct {
	Title    string   `json:"name"`
	Duration int      `json:"duration_ms"`
	Album    Album    `json:"album"`
	Artist   []Artist `json:"artists"`
}

type PlaybackState struct {
	IsPlaying bool  `json:"is_playing"`
	Progress  int   `json:"progress_ms"`
	Song      Track `json:"Item"`
}

func (api SpotApi) NewReqAuth(method, url, body string) *http.Response {
	req, _ := http.NewRequest("GET", "https://api.spotify.com/v1/me/player/currently-playing", nil)

	req.Header.Set("Authorization", "Bearer "+api.Config.AccessToken)

	client := &http.Client{}

	res, _ := client.Do(req)
	return res
}

func (Spot SpotApi) CurrentSong(Token string) PlaybackState {

	req, _ := http.NewRequest("GET", "https://api.spotify.com/v1/me/player/currently-playing", nil)

	req.Header.Set("Authorization", "Bearer "+Token)

	client := &http.Client{}

	res, err := client.Do(req)

	if err != nil {
		fmt.Println(err)
	}

	body, _ := io.ReadAll(res.Body)

	var state PlaybackState

	json.Unmarshal(body, &state)

	fmt.Printf("Playback State %+v", state)

	return state
}

func (api SpotApi) ChangeState() {

	api.NewReqAuth("GET", "", "")

}
