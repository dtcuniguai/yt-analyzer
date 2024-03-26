package cron

import (
	"fmt"
	"strconv"
	"strings"
	"time"
	"ytanalyzer/internal/analyzer"
	"ytanalyzer/internal/oauth"
	"ytanalyzer/internal/youtube"
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
		fmt.Println(1)
		return err
	}

	view, err := strconv.Atoi(chInfo.Items[0].Statistics.ViewCount)
	if err != nil {
		fmt.Println(2)
		return err
	}

	subscriber, err := strconv.Atoi(chInfo.Items[0].Statistics.SubscriberCount)
	if err != nil {
		fmt.Println(3)
		return err
	}

	videoCount, err := strconv.Atoi(chInfo.Items[0].Statistics.VideoCount)
	if err != nil {
		fmt.Println(4)
		return err
	}

	//update youtuber databse
	statistics := map[string]interface{}{
		"view":        view,
		"subscriber":  subscriber,
		"video_count": videoCount,
	}
	err = analyzer.UpdateYtber(chInfo.Items[0].ID, statistics)
	if err != nil {
		return err
	}

	//create log
	start := time.Now().AddDate(0, 0, -1).Format("2006-01-02")
	end := time.Now().Format("2006-01-02")
	mcs := "views,redViews,comments,likes,dislikes,videosAddedToPlaylists,videosRemovedFromPlaylists,shares,estimatedMinutesWatched,estimatedRedMinutesWatched,averageViewDuration,averageViewPercentage,annotationClickThroughRate,annotationCloseRate,annotationImpressions,annotationClickableImpressions,annotationClosableImpressions,annotationClicks,annotationCloses,cardClickRate,cardTeaserClickRate,cardImpressions,cardTeaserImpressions,cardClicks,cardTeaserClicks,subscribersGained,subscribersLost,estimatedRevenue,estimatedAdRevenue,grossRevenue,estimatedRedPartnerRevenue,monetizedPlaybacks,playbackBasedCpm,adImpressions,cpm"
	qMap := map[string]string{
		"ids":       "channel==MINE",
		"metrics":   mcs,
		"startDate": start,
		"endDate":   end,
	}

	metrics, rows, err := youtube.FetchTYAnalytics(ytb.Token, qMap)
	if err != nil {
		return err
	}

	measurement := "channel"
	fieldData := make(map[string]interface{})
	//make metrics
	if len(rows) == 0 {
		//所有統計資料為0
		headers := strings.Split(mcs, ",")
		for _, h := range headers {
			fieldData[h] = 0
		}
	} else {
		for i, header := range metrics {
			fieldData[header.Name] = rows[0][i]
		}
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

		//db取得影片資料
		vData, err := analyzer.GetVideosByIDs(vIDs)
		if err != nil {
			return err
		}

		//youtube api 取得影片詳細資料
		vRsp, err := youtube.FetchVideoDetail(ytb.Token, vIDs)
		if err != nil {
			return err
		}
		for _, video := range vRsp.Items {

			//TODO:for迴圈比對yt api影片是否已在db中、可能有更有效率的解法找出是否已存在db
			var exist bool
			for _, v := range vData {
				if v.ID == video.ID {
					exist = true
					break
				}
			}

			view, err := strconv.Atoi(video.Statistics.ViewCount)
			if err != nil {
				return err
			}

			like, err := strconv.Atoi(video.Statistics.LikeCount)
			if err != nil {
				return err
			}

			dislike, err := strconv.Atoi(video.Statistics.DislikeCount)
			if err != nil {
				return err
			}

			comment, err := strconv.Atoi(video.Statistics.CommentCount)
			if err != nil {
				return err
			}

			//資料庫更新、新增
			if exist {
				//更新db

				statistics := map[string]interface{}{
					"view":    view,
					"like":    like,
					"dislike": dislike,
					"comment": comment,
				}

				err = analyzer.UpdateVideo(video.ID, statistics)
				if err != nil {
					return err
				}
				fmt.Printf("影片[%v]更新完成\n", video.ID)
			} else {

				vTags := strings.Join(video.Snippet.Tags, ",")
				caption, err := strconv.ParseBool(video.ContentDetails.Caption)
				if err != nil {
					return err
				}

				//新增db
				vinfo := analyzer.Video{
					ID:           video.ID,
					ChannelID:    ytb.ID,
					Title:        video.Snippet.Title,
					Description:  video.Snippet.Description,
					DefaultThumb: video.Snippet.Thumbnails.Default.URL,
					MThumb:       video.Snippet.Thumbnails.Medium.URL,
					HThumb:       video.Snippet.Thumbnails.High.URL,
					SThumb:       video.Snippet.Thumbnails.Standard.URL,
					MaxThumb:     video.Snippet.Thumbnails.Maxres.URL,
					Tags:         vTags,
					Language:     video.Snippet.DefaultLanguage,
					Duration:     video.ContentDetails.Duration,
					Dimension:    video.ContentDetails.Dimension,
					Definition:   video.ContentDetails.Definition,
					Caption:      caption,
					View:         view,
					Like:         like,
					Dislike:      dislike,
					Comment:      comment,
					PublishAt:    video.Snippet.PublishedAt.Unix(),
				}

				err = analyzer.CreateVideo(&vinfo)
				if err != nil {
					return err
				}
				fmt.Printf("影片[%v]新增完成\n", video.ID)
			}

			//create log
			filters := fmt.Sprintf("video==%v", video.ID)
			start := time.Now().AddDate(0, 0, -1).Format("2006-01-02")
			end := time.Now().Format("2006-01-02")
			mcs := "views,redViews,comments,likes,dislikes,videosAddedToPlaylists,videosRemovedFromPlaylists,shares,estimatedMinutesWatched,estimatedRedMinutesWatched,averageViewDuration,averageViewPercentage,annotationClickThroughRate,annotationCloseRate,annotationImpressions,annotationClickableImpressions,annotationClosableImpressions,annotationClicks,annotationCloses,cardClickRate,cardTeaserClickRate,cardImpressions,cardTeaserImpressions,cardClicks,cardTeaserClicks,subscribersGained,subscribersLost,estimatedRevenue,estimatedAdRevenue,grossRevenue,estimatedRedPartnerRevenue,monetizedPlaybacks,playbackBasedCpm,adImpressions,cpm"
			qMap := map[string]string{
				"dimensions": "video",
				"ids":        "channel==MINE",
				"metrics":    mcs,
				"filters":    filters,
				"startDate":  start,
				"endDate":    end,
			}

			metrics, rows, err := youtube.FetchTYAnalytics(ytb.Token, qMap)
			if err != nil {
				return err
			}

			measurement := "video"
			tagData := make(map[string]string)
			fieldData := make(map[string]interface{})

			//make metrics
			if len(rows) == 0 {
				//所有統計資料為0
				tagData["video"] = video.ID
				headers := strings.Split(mcs, ",")
				for _, h := range headers {
					fieldData[h] = 0
				}

			} else {
				for i, header := range metrics {
					if header.DataType == "DIMENSION" {
						tag := fmt.Sprintf("%v", rows[0][i])
						tagData[header.Name] = tag
					} else {
						fieldData[header.Name] = rows[0][i]
					}
				}
			}

			//tag channel
			tagData["channel"] = ytb.ID
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
		// time.Sleep(1 * time.Second)
	}

	return nil
}
