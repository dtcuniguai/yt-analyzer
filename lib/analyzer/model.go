package analyzer

type Youtuber struct {
	ID           string `db:"id"`
	Title        string `db:"title"`
	Description  string `db:"description"`
	CustomURL    string `db:"custom_url"`
	DThumb       string `db:"default_thumb"`
	MThumb       string `db:"medium_thumb"`
	HThumb       string `db:"high_thumb"`
	Country      string `db:"country"`
	Token        string `db:"token"`
	RefreshToken string `db:"refresh_token"`
	PublishedAt  int64  `db:"publish_at"`
	CreateAt     int64  `db:"create_at"`
	UpdateAt     int64  `db:"update_at"`
}

type Video struct {
	ID           string `db:"id"`
	ChannelID    string `db:"channel_id"`
	Title        string `db:"title"`
	Description  string `db:"description"`
	DefaultThumb string `db:"default_thumb"`
	MThumb       string `db:"medium_thumb"`
	HThumb       string `db:"high_thumb"`
	SThumb       string `db:"standard_thumb"`
	MaxThumb     string `db:"maxres_thumb"`
	Tags         string `db:"tags"`
	Language     string `db:"language"`
	Duration     string `db:"duration"`
	Dimension    string `db:"dimension"`
	Definition   string `db:"definition"`
	Caption      bool   `db:"caption"`
	View         int    `db:"view"`
	Like         int    `db:"like"`
	Dislike      int    `db:"dislike"`
	Comment      int    `db:"comment"`
	PublishAt    int64  `db:"publish_at"`
	CreateAt     int64  `db:"create_at"`
	UpdateAt     int64  `db:"update_at"`
}
