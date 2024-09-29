package App

import (
	"context"
	"fmt"
	"lilSpotify/backend/spotApi"
	"log"
	"math"
	"os/exec"
	"strconv"
	"strings"

	"github.com/godbus/dbus"
)

type App struct {
	ctx context.Context
}

func NewApp() *App {
	return &App{}
}

func (a *App) Startup(ctx context.Context) {
	a.ctx = ctx
}

func (a *App) Greet(name string) string {
	return fmt.Sprintf("Hello %s, It's show time!", name)
}

var Method string

var api spotApi.SpotApi

func (a *App) Init(method string) bool {

	if method == "api" {
		api.Config = spotApi.Connect(a.ctx)
	}
	fmt.Printf("lmaoConfig %+v", api)
	api.CurrentSong(api.AccessToken)
	return true
}

func (a *App) CurrentSong() interface{} {

	if Method == "api" {
		xdd := api.CurrentSong(api.AccessToken)
		return xdd
	} else {

	}
	conn, err := dbus.SessionBus()
	if err != nil {
		log.Fatalf("Failed to connect to session bus: %v", err)
	}

	spotifyService := "org.mpris.MediaPlayer2.spotify"
	spotifyPath := dbus.ObjectPath("/org/mpris/MediaPlayer2")

	obj := conn.Object(spotifyService, spotifyPath)

	variant, err := obj.GetProperty("org.mpris.MediaPlayer2.Player.Metadata")

	xdd := map[string]string{
		"status": "closed",
	}

	if err != nil {
		return xdd
	}

	metadata := variant.Value().(map[string]dbus.Variant)

	pos, err := obj.GetProperty("org.mpris.MediaPlayer2.Player.Position")

	statusV, _ := obj.GetProperty("org.mpris.MediaPlayer2.Player.PlaybackStatus")

	if err != nil {
		log.Fatalf("Failed to get Position: %v", err)
	}

	albumCover := strings.Trim(metadata["mpris:artUrl"].String(), `"'`)
	artist := metadata["xesam:artist"].Value().([]string)
	trackId := metadata["mpris:trackid"].String()
	status := strings.Trim(statusV.String(), `"`)
	title := strings.Trim(metadata["xesam:title"].String(), `"`)
	album := strings.Trim(metadata["xesam:album"].String(), `"`)
	length := stringToFloat(metadata["mpris:length"].String(), "@t ")
	lengthF := timeParse(metadata["mpris:length"].String(), "@t ")
	lengthR := strings.Trim(metadata["mpris:length"].String(), "@t ")
	position := stringToFloat(pos.String(), "@x ")
	positionF := timeParse(pos.String(), "@x ")

	xdd = map[string]string{
		"appStatus":  "opened",
		"status":     status,
		"artist":     artist[0],
		"albumCover": albumCover,
		"title":      title,
		"album":      album,
		"length":     strconv.FormatFloat(length, 'f', 2, 32),
		"lengthF":    lengthF,
		"lengthR":    lengthR,
		"position":   strconv.FormatFloat(position, 'f', 2, 32),
		"positionF":  positionF,
		"trackId":    trackId,
	}

	return xdd
}

func (a *App) OpenApp() bool {
	cmd := exec.Command("spotify-launcher")

	cmd.Start()

	return true
}

func (a *App) ChangeState() {
	conn, err := dbus.SessionBus()
	if err != nil {
		log.Fatalf("Failed to connect to session bus: %v", err)
	}

	spotifyService := "org.mpris.MediaPlayer2.spotify"
	spotifyPath := dbus.ObjectPath("/org/mpris/MediaPlayer2")

	obj := conn.Object(spotifyService, spotifyPath)

	obj.Call("org.mpris.MediaPlayer2.Player.PlayPause", 0)
}

func (a *App) PrevSong() {
	conn, err := dbus.SessionBus()
	if err != nil {
		log.Fatalf("Failed to connect to session bus: %v", err)
	}

	spotifyService := "org.mpris.MediaPlayer2.spotify"
	spotifyPath := dbus.ObjectPath("/org/mpris/MediaPlayer2")

	obj := conn.Object(spotifyService, spotifyPath)

	obj.Call("org.mpris.MediaPlayer2.Player.Previous", 0)
}

func (a *App) NextSong() {
	conn, err := dbus.SessionBus()
	if err != nil {
		log.Fatalf("Failed to connect to session bus: %v", err)
	}

	spotifyService := "org.mpris.MediaPlayer2.spotify"
	spotifyPath := dbus.ObjectPath("/org/mpris/MediaPlayer2")

	obj := conn.Object(spotifyService, spotifyPath)

	obj.Call("org.mpris.MediaPlayer2.Player.Next", 0)
}

func (a *App) Seek(seekPosition int, trackId string) {
	conn, err := dbus.SessionBus()
	if err != nil {
		log.Fatalf("Failed to connect to session bus: %v", err)
	}

	// Define the Spotify service and object path
	spotifyService := "org.mpris.MediaPlayer2.spotify"
	playerPath := dbus.ObjectPath("/org/mpris/MediaPlayer2")

	// Get the D-Bus object for the Spotify player
	obj := conn.Object(spotifyService, playerPath)

	// Call the SetPosition method
	trackId = strings.Trim(trackId, "\"")
	err = obj.Call("org.mpris.MediaPlayer2.Player.SetPosition", 0, dbus.ObjectPath(trackId), int64(seekPosition)).Err
	if err != nil {
		log.Printf("Failed to call SetPosition: %v", err)
	} else {
		log.Println("Seek position set successfully")
	}
}

func stringToFloat(float string, cut string) float64 {
	lengthTemp, _ := strings.CutPrefix(strings.Trim(float, `"`), cut)
	lengthFloat, _ := strconv.ParseFloat(lengthTemp, 64)
	lengthFloat /= 1000000 * 60
	return lengthFloat
}

func timeParse(float string, cut string) string {
	lengthFloat := stringToFloat(float, cut)
	intP, fracP := math.Modf(lengthFloat)

	aloo := strconv.FormatFloat(intP, 'f', 0, 32)

	fracPS := strconv.FormatInt(int64(fracP*60), 10)

	if len(fracPS) == 1 {
		fracPS = "0" + fracPS
	}

	aloo += ":" + fracPS

	return aloo
}
