package main

import (
	"context"
	"fmt"
	"log"
	"math"
	"os/exec"
	"strconv"
	"strings"

	"github.com/godbus/dbus"
)

// App struct
type App struct {
	ctx context.Context
}

// NewApp creates a new App application struct
func NewApp() *App {
	return &App{}
}

// startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
}

// Greet returns a greeting for the given name
func (a *App) Greet(name string) string {
	return fmt.Sprintf("Hello %s, It's show time!", name)
}

func (a *App) CurrentSong() map[string]string {
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
	status := strings.Trim(statusV.String(), `"`)
	title := strings.Trim(metadata["xesam:title"].String(), `"`)
	album := strings.Trim(metadata["xesam:album"].String(), `"`)
	length := stringToFloat(metadata["mpris:length"].String(), "@t ")
	lengthF := timeParse(metadata["mpris:length"].String(), "@t ")
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
		"position":   strconv.FormatFloat(position, 'f', 2, 32),
		"positionF":  positionF,
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
