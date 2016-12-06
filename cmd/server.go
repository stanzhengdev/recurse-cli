package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"

	"github.com/davecgh/go-spew/spew"

	"golang.org/x/oauth2"
)

type BearerToken struct {
	AccessToken  string `json:"access_token"`
	TokenType    string `json:"token_type"`
	ExpiresIn    int    `json:"expires_in"`
	RefreshToken string `json:"refresh_token"`
	Scope        string `json:"scope"`
	CreatedAt    int    `json:"created_at"`
}


/// Get an OAuth Token from
func main() {
	oauthConfig := oauth2.Config{
		ClientID:    os.Getenv("ClientID"),
		ClientSecret: os.Getenv("ClientSecret"),
		RedirectURL:  os.Getenv("RedirectURL"),
		Endpoint: oauth2.Endpoint{
			TokenURL: "https://www.recurse.com/oauth/token",
			AuthURL:  "https://www.recurse.com/oauth/authorize",
		},
		Scopes: []string{},
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "<a href=%s>login</a>", buildAuthUrl(oauthConfig))
	})

	http.HandleFunc("/authed", func(w http.ResponseWriter, r *http.Request) {
		authcode := r.URL.Query().Get("code")
		// spew.Dump(r)
		fmt.Println(authorizeToken(oauthConfig, authcode))

		io.WriteString(w, authcode)
		w.Write([]byte(authcode))
	})
	fmt.Println("starting server")
	http.ListenAndServe(":4000", nil)
}


func buildAuthUrl(state oauth2.Config) *url.URL {
	u, err := url.Parse(state.Endpoint.AuthURL)
	if err != nil {
		log.Fatal(err)
	}
	uq := u.Query()
	uq.Add("response_type", "code")
	uq.Add("client_id", state.ClientID)
	uq.Add("redirect_uri", state.RedirectURL)
	uq.Add("access_type", "online")
	u.RawQuery = uq.Encode()
	return u
}

func authorizeToken(c oauth2.Config, code string) *BearerToken {
	// Exchange Code with token endpoint to receive a bearer token
    // Returns a Bearer Token
	u, _ := url.Parse(c.Endpoint.TokenURL)
	uq := u.Query()
	uq.Add("grant_type", "authorization_code")
	uq.Add("code", code)
	uq.Add("redirect_uri", c.RedirectURL)

	u.RawQuery = uq.Encode()

	// Build request
	hc := http.Client{}
	req, err := http.NewRequest("POST", u.String(), nil)
	if err != nil {
		fmt.Println(err)
	}
    // Apply approppriate headers
	req.Header.Add("Content-Type", "application/x-www.form-urlencoded")
	req.SetBasicAuth(c.ClientID, c.ClientSecret)

	resp, err := hc.Do(req)
	if err != nil {
		fmt.Println(err)
	}

	bear := &BearerToken{}
	err = json.NewDecoder(resp.Body).Decode(bear)
)
	if err != nil {
		fmt.Println(err)
	}
	return bear
}
