package handler

import (
	"fmt"
	"ytanalyzer/lib/analyzer"
	"ytanalyzer/lib/oauth"
	"ytanalyzer/lib/youtube"

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

	//save user to db
	chInfo, err := youtube.FetchSelfChannelDetail(account.AccessToken)
	if err != nil {
		return c.SendString(err.Error())
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
		Token:        account.AccessToken,
		RefreshToken: account.RefreshToken,
	}

	err = analyzer.RegisterYoutuber(ytinfo)
	if err != nil {
		return c.SendString(err.Error())
	}

	fmt.Printf("頻道[%v] 新增成功ID: %v\n", ytinfo.Title, ytinfo.ID)
	msg := fmt.Sprintf("授權成功，您好 %v", ytinfo.Title)
	return c.SendString(msg)
}
