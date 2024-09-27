package spotApi

import (
	"context"
	"encoding/base64"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"strings"

	"github.com/wailsapp/wails/v2/pkg/runtime"
)

type PostBody struct {
	grant_type   string
	code         interface{}
	redirect_uri string
}

var clientId string = "0a23cf5a3c2546a58c2a19ba680ffef7"
var clientSec string = "4875510c44b94a11937c656c2d4e06e2"
var redirectUri string = "http://localhost:8080/auth"
var scopes string = "user-read-private"

var urla string = fmt.Sprintf(
	"https://accounts.spotify.com/authorize?client_id=%s&response_type=code&redirect_uri=%s&scope=%s",
	clientId,
	redirectUri,
	scopes,
)

func startLocalServer(ctx context.Context) {
	http.HandleFunc("/auth", func(w http.ResponseWriter, r *http.Request) {
		code := r.URL.Query().Get("code")

		runtime.EventsEmit(ctx, "auth", code)

		fmt.Fprintf(w, "Authentication successful. You can close this window.")
	})

	go func() {
		fmt.Println("xdd")
		log.Fatal(http.ListenAndServe(":8080", nil))
	}()
}

func Connect(ctx context.Context) string {

	runtime.BrowserOpenURL(ctx, urla)
	startLocalServer(ctx)

	runtime.EventsOn(ctx, "auth", func(optionalData ...interface{}) {
		form := url.Values{}
		form.Add("grant_type", "authorization_code")
		form.Add("code", optionalData[0].(string)) // Replace with the actual code received
		form.Add("redirect_uri", redirectUri)      // Replace with your redirect URI
		form.Add("client_id", clientId)            // Replace with your client ID
		form.Add("client_secret", clientSec)       // Replace with your client secret

		req, err := http.NewRequest("POST", "https://accounts.spotify.com/api/token", strings.NewReader(form.Encode()))
		if err != nil {
			fmt.Println("Error creating request:", err)
			return
		}

		basicAuth := base64.StdEncoding.EncodeToString([]byte(clientId + ":" + clientSec))

		req.Header.Set("Authorization", "Basic "+basicAuth)
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			fmt.Println(err)
			return
		}

		bodx, _ := io.ReadAll(resp.Body)
		runtime.LogPrintf(ctx, "xdd %s", string(bodx))

	})

	// cmd := exec.Command("xdg-open", url)
	//
	// err := cmd.Run()
	//
	// if err != nil {
	// 	fmt.Println("Failed Running Command")
	// 	return
	// }

}
