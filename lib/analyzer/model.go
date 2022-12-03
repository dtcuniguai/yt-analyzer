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
