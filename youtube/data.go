package youtube

type ChannelInfo struct {
	ID              string
	Title           string
	Description     string
	SubscriberCount string
	VideoCount      string
	ViewCount       string
	Thumbnails      map[string]string
}

//youtube subscriptions api response
type SubRsp struct {
	Kind          string `json:"kind"`
	Etag          string `json:"etag"`
	NextPageToken string `json:"nextPageToken"`
	PageInfo      struct {
		TotalResults   int `json:"totalResults"`
		ResultsPerPage int `json:"resultsPerPage"`
	} `json:"pageInfo"`
	Items []struct {
		Kind              string `json:"kind"`
		Etag              string `json:"etag"`
		ID                string `json:"id"`
		SubscriberSnippet struct {
			Title       string `json:"title"`
			Description string `json:"description"`
			ChannelID   string `json:"channelId"`
			Thumbnails  struct {
				Default struct {
					URL string `json:"url"`
				} `json:"default"`
				Medium struct {
					URL string `json:"url"`
				} `json:"medium"`
				High struct {
					URL string `json:"url"`
				} `json:"high"`
			} `json:"thumbnails"`
		} `json:"subscriberSnippet"`
	} `json:"items"`
}

type ChannelRsp struct {
	Kind     string `json:"kind"`
	Etag     string `json:"etag"`
	PageInfo struct {
		TotalResults   int `json:"totalResults"`
		ResultsPerPage int `json:"resultsPerPage"`
	} `json:"pageInfo"`
	Items []struct {
		Kind       string `json:"kind"`
		Etag       string `json:"etag"`
		ID         string `json:"id"`
		Statistics struct {
			ViewCount             string `json:"viewCount"`
			SubscriberCount       string `json:"subscriberCount"`
			HiddenSubscriberCount bool   `json:"hiddenSubscriberCount"`
			VideoCount            string `json:"videoCount"`
		} `json:"statistics"`
	} `json:"items"`
}
