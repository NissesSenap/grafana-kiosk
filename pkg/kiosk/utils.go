package kiosk

import (
	"crypto/tls"
	"log"
	"net/http"
	"net/url"
	"path"

	grapi "github.com/grafana/grafana-api-golang-client"
)

// GenerateURL constructs URL with appropriate parameters for kiosk mode
func GenerateURL(anURL string, kioskMode string, autoFit bool, isPlayList bool) string {
	u, _ := url.ParseRequestURI(anURL)
	q, _ := url.ParseQuery(u.RawQuery)

	switch kioskMode {
	case "tv": // TV
		q.Set("kiosk", "tv") // no sidebar, topnav without buttons
		log.Printf("KioskMode: TV")
	case "full": // FULLSCREEN
		q.Set("kiosk", "1") // sidebar and topnav always shown
		log.Printf("KioskMode: Fullscreen")
	case "disabled": // FULLSCREEN
		log.Printf("KioskMode: Disabled")
	default: // disabled
		q.Set("kiosk", "1") // sidebar and topnav always shown
		log.Printf("KioskMode: Fullscreen")
	}
	// a playlist should also go inactive immediately
	if isPlayList {
		q.Set("inactive", "1")
	}
	u.RawQuery = q.Encode()
	if autoFit {
		u.RawQuery += "&autofitpanels"
	}
	return u.String()
}

func NewGrafanaClient(anURL, username, password string, ignoreCertErrors bool) (*grapi.Client, error) {
	userinfo := url.UserPassword(username, password)
	clientConfig := grapi.Config{
		APIKey:    "",
		BasicAuth: userinfo,
		Client: &http.Client{
			Transport: &http.Transport{
				TLSClientConfig: &tls.Config{
					InsecureSkipVerify: ignoreCertErrors,
				},
			},
		},
		OrgID:      0,
		NumRetries: 0,
	}

	u, err := url.Parse(anURL)
	if err != nil {
		return nil, err
	}
	u.Path = ""
	u.RawQuery = ""
	u.Fragment = ""

	grafanaClient, err := grapi.New(u.String(), clientConfig)
	if err != nil {
		return nil, err
	}

	return grafanaClient, nil
}

// getPlayListUID, get the UID of a playlist from an id
func GetPlayListUID(anURL string, client *grapi.Client) (string, error) {
	id, err := url.Parse(anURL)
	if err != nil {
		return "", err
	}

	platList, err := client.Playlist(path.Base(id.Path))
	if err != nil {
		return "", err
	}
	return platList.UID, nil
}

/*
func ChangeIDtoUID(anURL, uid string) {
	urlA, _ := url.Parse(anURL)
	urlPATH := urlA.Path
	// TODO grab the output and change the last value.
	Some split
}
*/
