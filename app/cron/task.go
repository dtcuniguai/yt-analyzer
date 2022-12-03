package cron

import (
	"fmt"
	"time"
	"ytanalyzer/lib/analyzer"
	"ytanalyzer/lib/oauth"
	"ytanalyzer/lib/youtube"
)

/* youtube排程紀錄資料
 *
 */
func TaskSyncYT() {
	err := SyncYTbers()
	if err != nil {
		fmt.Println(err)
	}
}

/* youtube 頻道主資料更新
 *
 */
func SyncYTbers() error {

	page := 1
	for {
		offset := (page - 1) * analyzer.QUERY_PERPAGE
		ytbs, err := analyzer.GetAllYtbers(offset, analyzer.QUERY_PERPAGE)
		if err != nil {
			return err
		}

		//youtuber 資料同步 & 紀錄log
		for _, ytb := range ytbs {

			//refresh token
			rt, err := oauth.RefreshToken(ytb.RefreshToken)
			if err != nil {
				return err
			}
			ytb.Token = rt.AccessToken

			err = LogYTberChannel(ytb)
			if err != nil {
				return err
			}

			err = LogYTberVideo(ytb)
			if err != nil {
				return err
			}
		}

		//fetch完畢
		if len(ytbs) < analyzer.QUERY_PERPAGE {
			fmt.Println("同步完畢")
			break
		}
	}

	return nil
}

/* youtuber 頻道相關log數據紀錄
 *
 */
func LogYTberChannel(ytb *analyzer.Youtuber) error {

	chInfo, err := youtube.FetchSelfChannelDetail(ytb.Token)
	if err != nil {
		return err
	}

	measurement := "channel"
	fieldData := map[string]interface{}{
		"view_count":  chInfo.Items[0].Statistics.ViewCount,
		"subscriber":  chInfo.Items[0].Statistics.SubscriberCount,
		"video_count": chInfo.Items[0].Statistics.VideoCount,
	}

	tagData := map[string]string{
		"channel_id": chInfo.Items[0].ID,
	}

	fmt.Printf("create channel log: %v\n", ytb.ID)
	err = analyzer.CreateLog(measurement, tagData, fieldData)
	if err != nil {
		return err
	}
	return nil
}

/* youtuber 影片相關log數據紀錄
 *
 */
func LogYTberVideo(ytb *analyzer.Youtuber) error {

	var protect uint
	var pageToken string
	for {

		query := map[string]string{
			"access_token": ytb.Token,
			"channelId":    ytb.ID,
			"part":         "snippet",
			"formine":      "true",
			"type":         "video",
			"maxResults":   "50",
		}
		if len(pageToken) != 0 {
			query["pageToken"] = pageToken
		}

		//query video
		vList, err := youtube.SearchVideoList(query)
		if err != nil {
			return err
		}

		//generate video ids for query video detail
		var vIDs string
		for _, item := range vList.Items {
			if len(vIDs) != 0 {
				vIDs = fmt.Sprintf("%v,%v", vIDs, item.ID.VideoID)
			} else {
				vIDs = item.ID.VideoID
			}
		}

		vRsp, err := youtube.FetchVideoDetail(ytb.Token, vIDs)
		if err != nil {
			return err
		}
		for _, video := range vRsp.Items {
			measurement := "video"
			fieldData := map[string]interface{}{
				"view":    video.Statistics.ViewCount,
				"like":    video.Statistics.LikeCount,
				"dislike": video.Statistics.DislikeCount,
				"comment": video.Statistics.CommentCount,
			}

			tagData := map[string]string{
				"id":         video.ID,
				"channel_id": ytb.ID,
			}

			fmt.Printf("create video log: %v\n", video.ID)
			err = analyzer.CreateLog(measurement, tagData, fieldData)
			if err != nil {
				return err
			}
		}

		if len(vList.NextPageToken) != 0 {
			pageToken = vList.NextPageToken
		} else {
			fmt.Printf("頻道主%v 影片資料同步完成\n", ytb.Title)
			break
		}

		//避免err無窮迴圈保護機制
		if protect > 1000000 {
			break
		}
		time.Sleep(3 * time.Second)

	}

	return nil
}
