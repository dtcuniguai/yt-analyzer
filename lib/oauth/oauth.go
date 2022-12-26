package oauth

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strings"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

var config *oauth2.Config

func Config() *oauth2.Config {

	if config == nil {
		redirect := fmt.Sprintf("%v/oauth/redirect", os.Getenv("APP_URL"))

		c := oauth2.Config{
			ClientID:     os.Getenv("GOOGLE_OAUTH_CLIENTID"),
			ClientSecret: os.Getenv("GOOGLE_OAUTH_SECRET"),
			RedirectURL:  redirect,
			Scopes: []string{
				"https://www.googleapis.com/auth/userinfo.email",
				"https://www.googleapis.com/auth/userinfo.profile",
				"https://www.googleapis.com/auth/youtube.readonly",
				"https://www.googleapis.com/auth/yt-analytics.readonly",
				"https://www.googleapis.com/auth/yt-analytics-monetary.readonly",
			},
			Endpoint: google.Endpoint,
		}

		config = &c
	}

	return config
}

/* google oauth 授權url產生
 *
 */
func AuthrozieUrl() string {
	c := Config()
	return c.AuthCodeURL("", oauth2.AccessTypeOffline, oauth2.ApprovalForce)
}

/* google oauth 使用者token過期刷新重新取得新token
 *
 */
func RefreshToken(refreshToken string) (*OToken, error) {

	client := &http.Client{}
	u := "https://www.googleapis.com/oauth2/v4/token"
	cfg := Config()

	//set post form
	data := url.Values{}
	data.Set("client_id", cfg.ClientID)
	data.Set("client_secret", cfg.ClientSecret)
	data.Set("refresh_token", refreshToken)
	data.Set("grant_type", "refresh_token")

	req, err := http.NewRequest(http.MethodPost, u, strings.NewReader(data.Encode()))
	if err != nil {
		return nil, err
	}

	//set header if needed
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	// req.Header.Set("Cookie", "")

	//send request
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var user OToken
	err = json.Unmarshal([]byte(b), &user)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

/* google oauth 透過認證回傳的code取得授權者的token及refresh token
 *
 */
func TokenInfo(code string) (*OToken, error) {

	client := &http.Client{}
	u := "https://www.googleapis.com/oauth2/v4/token"
	cfg := Config()

	//set post form
	data := url.Values{}
	data.Set("client_id", cfg.ClientID)
	data.Set("client_secret", cfg.ClientSecret)
	data.Set("redirect_uri", cfg.RedirectURL)
	data.Set("code", code)
	data.Set("grant_type", "authorization_code")

	req, err := http.NewRequest(http.MethodPost, u, strings.NewReader(data.Encode()))
	if err != nil {
		return nil, err
	}

	//set header if needed
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	// req.Header.Set("Cookie", "")

	//send request
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var user OToken
	err = json.Unmarshal([]byte(b), &user)
	if err != nil {
		return nil, err
	}

	return &user, nil
}
