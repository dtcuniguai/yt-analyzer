package analyzer

import (
	"fmt"
	"strings"
	"time"

	"github.com/jmoiron/sqlx"
)

const QUERY_PERPAGE = 30

/* 新增youtuber資料
 * table: sqlite.youtuber
 * param: Youtuber (obj)
 */
func CreateYoutuber(ytber *Youtuber) error {

	db, err := GetDB()
	if err != nil {
		return err
	}

	n := time.Now().Unix()
	dataMap := map[string]interface{}{
		"id":            ytber.ID,
		"title":         ytber.Title,
		"description":   ytber.Description,
		"custom_url":    ytber.CustomURL,
		"default_thumb": ytber.DThumb,
		"medium_thumb":  ytber.MThumb,
		"high_thumb":    ytber.HThumb,
		"country":       ytber.Country,
		"token":         ytber.Token,
		"refresh_token": ytber.RefreshToken,
		"publish_at":    ytber.PublishedAt,
		"view":          ytber.View,
		"subscriber":    ytber.Subscriber,
		"video_count":   ytber.VideoCount,
		"create_at":     n,
		"update_at":     n,
	}

	query := "INSERT INTO youtuber (`id`,`title`,`description`,`custom_url`,`default_thumb`,`medium_thumb`,`high_thumb`,`country`,`token`, `view`, `subscriber`, `video_count`,`refresh_token`,`publish_at`,`create_at`,`update_at`) VALUES (:id,:title,:description,:custom_url,:default_thumb,:medium_thumb,:high_thumb,:country,:token, :view, :subscriber, :video_count,:refresh_token,:publish_at,:create_at,:update_at)"
	_, err = db.NamedExec(query, dataMap)
	if err != nil {
		return err
	}

	return nil
}

/* 更新youtuber資料
 * table: sqlite.youtuber
 * param: vid、updateField(當前數據更新)
 */
func UpdateYtber(id string, updateField map[string]interface{}) error {

	db, err := GetDB()
	if err != nil {
		return err
	}

	updateField["id"] = id
	updateField["update_at"] = time.Now().Unix()
	query := "UPDATE `youtuber` SET  view=:view, subscriber=:subscriber, video_count=:video_count, update_at=:update_at WHERE id = :id"
	_, err = db.NamedExec(query, updateField)
	if err != nil {
		return err
	}

	return nil
}

/* 建立video資料
 * table: sqlite.video
 * param: Video (obj)
 */
func CreateVideo(v *Video) error {

	db, err := GetDB()
	if err != nil {
		return err
	}

	n := time.Now().Unix()
	dataMap := map[string]interface{}{
		"id":             v.ID,
		"channel_id":     v.ChannelID,
		"title":          v.Title,
		"description":    v.Description,
		"default_thumb":  v.DefaultThumb,
		"medium_thumb":   v.MThumb,
		"high_thumb":     v.HThumb,
		"standard_thumb": v.SThumb,
		"maxres_thumb":   v.MaxThumb,
		"tags":           v.Tags,
		"language":       v.Language,
		"duration":       v.Duration,
		"dimension":      v.Dimension,
		"definition":     v.Definition,
		"caption":        v.Caption,
		"view":           v.View,
		"like":           v.Like,
		"dislike":        v.Dislike,
		"comment":        v.Comment,
		"publish_at":     v.PublishAt,
		"create_at":      n,
		"update_at":      n,
	}

	query := "INSERT INTO video (`id`, `channel_id`, `title`, `description`, `default_thumb`, `medium_thumb`, `high_thumb`, `standard_thumb`, `maxres_thumb`, `tags`, `language`, `duration`, `dimension`, `definition`, `caption`, `view`, `like`, `dislike`, `comment`, `publish_at`, `create_at`, `update_at`) VALUES (:id, :channel_id, :title, :description, :default_thumb, :medium_thumb, :high_thumb, :standard_thumb, :maxres_thumb, :tags, :language, :duration, :dimension, :definition, :caption, :view, :like, :dislike, :comment, :publish_at, :create_at, :update_at)"
	_, err = db.NamedExec(query, dataMap)
	if err != nil {
		fmt.Println("err here")
		return err
	}

	return nil
}

/* 更新video資料
 * table: sqlite.video
 * param: vid、updateField(當前數據更新)
 */
func UpdateVideo(vid string, updateField map[string]interface{}) error {

	db, err := GetDB()
	if err != nil {
		return err
	}

	updateField["vid"] = vid
	updateField["update_at"] = time.Now().Unix()
	query := "UPDATE `video` SET  view=:view, like=:like, dislike=:dislike, comment=:comment, update_at=:update_at WHERE id = :vid"
	_, err = db.NamedExec(query, updateField)
	if err != nil {
		return err
	}

	return nil
}

/* 取得已授權的youtuber列表
 * table: sqlite.youtube
 * param: offset(分頁定位)、limit(分頁數)
 */
func GetAllYtbers(offset int, limit int) ([]*Youtuber, error) {

	db, err := GetDB()
	if err != nil {
		return nil, err
	}

	var ytbers []*Youtuber
	query := "SELECT * FROM `youtuber` LIMIT :limit OFFSET :offset"
	rows, err := db.NamedQuery(query, map[string]interface{}{
		"offset": offset,
		"limit":  limit,
	})
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var ytb Youtuber
		err := rows.StructScan(&ytb)
		if err != nil {
			return nil, err
		}

		ytbers = append(ytbers, &ytb)
	}

	return ytbers, nil
}

/* 影片ID取得影片列表
 * table: sqlite.video
 * param: video id (id字串以,分隔)
 */
func GetVideosByIDs(vids string) ([]*Video, error) {

	db, err := GetDB()
	if err != nil {
		return nil, err
	}

	var videos []*Video
	ids := strings.Split(vids, ",")
	query, args, err := sqlx.In("SELECT * FROM `video` WHERE id in (?);", ids)
	if err != nil {
		return nil, err
	}

	query = db.Rebind(query)
	rows, err := db.Queryx(query, args...)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var v Video
		err := rows.StructScan(&v)
		if err != nil {
			return nil, err
		}

		videos = append(videos, &v)
	}

	return videos, nil
}
