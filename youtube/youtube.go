package youtube

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
)

/*
youtube api interface
*/

//https://www.googleapis.com/youtube/v3/channels?part=statistics&id={{id}}&key={{API_KEY}}
func GetChannelDetail(apiKey string, channelId string) (*ChannelRsp, error) {

	client := &http.Client{}
	url := "https://www.googleapis.com/youtube/v3/channels"

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Print(err)
		return nil, err
	}

	q := req.URL.Query()
	q.Add("part", "statistics")
	q.Add("id", channelId)
	q.Add("key", apiKey)
	req.URL.RawQuery = q.Encode()

	resp, err := client.Do(req)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	defer resp.Body.Close()
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	var ch ChannelRsp
	err = json.Unmarshal([]byte(data), &ch)
	if err != nil {
		return nil, err
	}

	return &ch, nil
}

//https://www.googleapis.com/youtube/v3/subscriptions?part=subscriberSnippet&channelId={{CHANNEL_ID}}&key={{API_KEY}}&maxResults=1
func GetSubscriptionsDetail(apiKey string, channelId string) (*SubRsp, error) {

	client := &http.Client{}
	url := "https://www.googleapis.com/youtube/v3/subscriptions"

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Print(err)
		return nil, err
	}

	q := req.URL.Query()
	q.Add("part", "subscriberSnippet")
	q.Add("channelId", channelId)
	q.Add("key", apiKey)
	q.Add("maxResults", "1")
	req.URL.RawQuery = q.Encode()

	resp, err := client.Do(req)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	defer resp.Body.Close()
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	var sub SubRsp
	err = json.Unmarshal([]byte(data), &sub)
	if err != nil {
		return nil, err
	}

	return &sub, nil
}
