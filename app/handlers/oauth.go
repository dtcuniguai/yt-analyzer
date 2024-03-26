package handlers

import (
	"fmt"
	"strconv"
	"ytanalyzer/internal/analyzer"
	"ytanalyzer/internal/oauth"
	"ytanalyzer/internal/youtube"

	"github.com/gofiber/fiber/v2"
)

type Oauth struct {
	Handler
}

/* google app oauth 授權連結
 *
 */
func (o Oauth) Auth(c *fiber.Ctx) error {
	url := oauth.AuthrozieUrl()
	return c.Redirect(url)
}

/* google app oauth授權完成返回
 * 透過google回傳 code 取得授權用戶token並儲存在資料庫中
 */
func (o Oauth) Redirect(c *fiber.Ctx) error {

	code := c.Query("code")
	account, err := oauth.TokenInfo(code)
	if err != nil {
		return c.SendString(err.Error())
	}

	chInfo, err := youtube.FetchSelfChannelDetail(account.AccessToken)
	if err != nil {
		return c.SendString(err.Error())
	}

	view, err := strconv.Atoi(chInfo.Items[0].Statistics.ViewCount)
	if err != nil {
		return err
	}

	videoCount, err := strconv.Atoi(chInfo.Items[0].Statistics.VideoCount)
	if err != nil {
		return err
	}

	subscriber, err := strconv.Atoi(chInfo.Items[0].Statistics.SubscriberCount)
	if err != nil {
		return err
	}

	//data formatting
	ytinfo := analyzer.Youtuber{
		ID:           chInfo.Items[0].ID,
		Title:        chInfo.Items[0].Snippet.Title,
		Description:  chInfo.Items[0].Snippet.Description,
		CustomURL:    chInfo.Items[0].Snippet.CustomURL,
		PublishedAt:  chInfo.Items[0].Snippet.PublishedAt.Unix(),
		DThumb:       chInfo.Items[0].Snippet.Thumbnails.Default.URL,
		MThumb:       chInfo.Items[0].Snippet.Thumbnails.Medium.URL,
		HThumb:       chInfo.Items[0].Snippet.Thumbnails.High.URL,
		Country:      chInfo.Items[0].Snippet.Country,
		View:         view,
		VideoCount:   videoCount,
		Subscriber:   subscriber,
		Token:        account.AccessToken,
		RefreshToken: account.RefreshToken,
	}

	//save user to db
	err = analyzer.CreateYoutuber(&ytinfo)
	if err != nil {
		return c.SendString(err.Error())
	}
	fmt.Printf("頻道[%v] 新增成功ID: %v\n", ytinfo.Title, ytinfo.ID)

	//create user log
	measurement := "channel"
	fieldData := map[string]interface{}{
		"view_count":  chInfo.Items[0].Statistics.ViewCount,
		"subscriber":  chInfo.Items[0].Statistics.SubscriberCount,
		"video_count": chInfo.Items[0].Statistics.VideoCount,
	}

	tagData := map[string]string{
		"channel_id": chInfo.Items[0].ID,
	}

	fmt.Printf("create channel log: %v\n", ytinfo.ID)
	err = analyzer.CreateLog(measurement, tagData, fieldData)
	if err != nil {
		return err
	}

	msg := fmt.Sprintf("授權成功，您好 %v", ytinfo.Title)
	return c.SendString(msg)
}
