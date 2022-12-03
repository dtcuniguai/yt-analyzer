package analyzer

import "time"

/* 新增youtuber基本資料到db裡面去
 *
 */
func RegisterYoutuber(ytber Youtuber) error {

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
		"create_at":     n,
		"update_at":     n,
	}

	query := "INSERT INTO youtuber (`id`,`title`,`description`,`custom_url`,`default_thumb`,`medium_thumb`,`high_thumb`,`country`,`token`,`refresh_token`,`publish_at`,`create_at`,`update_at`) VALUES (:id,:title,:description,:custom_url,:default_thumb,:medium_thumb,:high_thumb,:country,:token,:refresh_token,:publish_at,:create_at,:update_at)"
	_, err = db.NamedExec(query, dataMap)
	if err != nil {
		return err
	}

	return nil
}
