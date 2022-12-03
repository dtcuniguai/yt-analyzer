package youtube

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"net/http"
)

/*
 * youtube api interface
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
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusForbidden {
		return nil, errors.New("channel : " + channelId + " doesn't have permission to get infos")
	}

	if resp.StatusCode == http.StatusNotFound {
		return nil, errors.New("channel : " + channelId + " resource not found")
	}

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
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

	if resp.StatusCode == http.StatusForbidden {
		return nil, errors.New("channel : " + channelId + " doesn't have permission to get infos")
	}

	if resp.StatusCode == http.StatusNotFound {
		return nil, errors.New("channel : " + channelId + " resource not found")
	}

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

/* 取得自身頻道詳細資料
 *
 */
func FetchSelfChannelDetail(token string) (*ChannelRsp, error) {
	client := &http.Client{}
	url := "https://www.googleapis.com/youtube/v3/channels"

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Print(err)
		return nil, err
	}

	q := req.URL.Query()
	q.Add("part", "snippet,contentDetails,statistics,brandingSettings")
	q.Add("access_token", token)
	q.Add("mine", "true")
	req.URL.RawQuery = q.Encode()

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var ch ChannelRsp
	err = json.Unmarshal([]byte(data), &ch)
	if err != nil {
		return nil, err
	}

	return &ch, nil
}

/* 影片條件搜尋
 *
 */
func SearchVideoList(qMap map[string]string) (*SearchRsp, error) {
	client := &http.Client{}
	url := "https://www.googleapis.com/youtube/v3/search"

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Print(err)
		return nil, err
	}

	q := req.URL.Query()
	for k, v := range qMap {
		q.Add(k, v)
	}

	// q.Add("access_token", token)
	// q.Add("part", "snippet")
	// q.Add("formine", "true")
	// q.Add("type", "video")
	// q.Add("maxResults", "5")
	// q.Add("channelId", "")

	req.URL.RawQuery = q.Encode()
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var list SearchRsp
	err = json.Unmarshal([]byte(data), &list)
	if err != nil {
		return nil, err
	}

	return &list, nil
}

/* 影片詳細資料
 *
 */
func FetchVideoDetail(token string, id string) (*VideoRsp, error) {
	client := &http.Client{}
	url := "https://www.googleapis.com/youtube/v3/videos"

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Print(err)
		return nil, err
	}

	q := req.URL.Query()
	q.Add("access_token", token)
	q.Add("part", "snippet,contentDetails,statistics")
	q.Add("id", id)

	req.URL.RawQuery = q.Encode()

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var list VideoRsp
	err = json.Unmarshal([]byte(data), &list)
	if err != nil {
		return nil, err
	}

	return &list, nil
}
