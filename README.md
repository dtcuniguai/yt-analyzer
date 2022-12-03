# youtube-command-line-tool

簡單的 youtube data api串接

[Youtube Data Api 官方文件](https://developers.google.com/youtube/v3/docs)


### How to use

- [申請GOOGLE API KEY](https://console.cloud.google.com/apis/api/youtube.googleapis.com/overview)

![Imgur](https://i.imgur.com/gd8N9JZ.png)

- 複製.env
```shell
cp .env.example  .env   
```

- API KEY 放入.env
```shell
GOOGLE_API_KEY={YOUR_API_KEY}
```

- run command (取得 youtube channel id [ep:](https://www.youtube.com/channel/UCoSrY_IQQVpmIRZ9Xf-y93g))

```shell
go run main.go -c UCoSrY_IQQVpmIRZ9Xf-y93g
```


- 接著可以看到頻道相關的資料在json檔案裡面

![Imgur](https://i.imgur.com/mehk3ub.png)

### command flag

- `-o dir, --output dir` : json 資料產生的位置
- `-c channels, --channels channels` : youtube頻道ID 多個頻道以","做分隔
- `-s {file location}, --source {file location}` : youtube頻道ID檔案，以換行符號分割


### 尚未完成功能

- youtube api error handle(優先：部分channel權限未開放會拿到403無法取得資料)
- server hosting
- 部分功能參數化
    - channels read file
    - more output file type support


### youtube response

#### 頻道資料
- [資料來源](https://developers.google.com/youtube/v3/docs/channels#resource)
    |          |  id   | title  | description | customUrl | publishedAt | thumbnails | localized | country | viewCount | subscriberCount | hiddenSubscriberCount | videoCount |
    |  ----    | ----  | ----  | ----  | ----  | ----  | ----  | ----  | ----  |  ----  |  ----  |  ----  |  ----  |
    | 資料名稱  | id | 頻道名稱 | 頻道說明 | 自定義連結？ | 加入時間 | 頻道頭像url | ? | 地區位置 | 總觀看次數 | 訂閱人數 | 是否顯示訂閱人數 | 總觀看次數
    | 儲存類型  | info | info | info | info | info | info | info | info |
    | 是否記錄  | ✅ | ✅ | 

#### 影片資料
- [資料來源](https://developers.google.com/youtube/v3/docs/videos)

    |          |  id   | title  | description | customUrl | publishedAt | thumbnails | localized | country |
    |  ----    | ----  | ----  | ----  | ----  | ----  | ----  | ----  | ----  |
    | 資料名稱  | id | 頻道名稱 | 頻道說明 | 自定義連結？ | 加入時間 | 頻道頭像url | ? | 地區位置
    | 儲存類型  | info | info | info | info | info | info | info | info |
    | 是否記錄  | ✅ | ✅ | 


## DataBase

### Sqlite 
- youtuber

|           | channel_id | title  | description | customUrl | publishedAt | thumbnails | localized | country | token | refresh_token |
|  ----    | ----  | ----  | ----  | ----  | ----  | ----  | ----  | ----  | ---- | ---- | ---- |
| 資料名稱  | youtube 頻道id | 頻道名稱 | 頻道說明 | 自定義連結？ | 加入時間 | 頻道頭像url | ? | 地區位置 | oauth授權token | oauth授權刷新用token
| 儲存類型  |  | TEXT | TEXT | TEXT | INTEGER | TEXT | info | TEXT | TEXT | TEXT |

- video

|          | video_id | channel_id | title  | description | default_thumb | medium_thumb | high_thumb | standard_thumb | maxres_thumb | tags | language | duration | dimension | definition | caption | publishedAt | view_count | like | dislike | comment |
|  ----    | ----  | ----  | ----  | ----  | ----  | ----  | ----  | ----  | ---- | ---- | ---- |  ---- |  ---- |  ---- | ---- | ---- | ---- | ---- | ---- | ---- |
| 資料名稱 | youtube 影片ID  | youtube 頻道id | 影片名稱 | 影片說明 | 縮圖 | 縮圖 | 縮圖 | 縮圖 | 縮圖 | 標籤 | 語系 | 影片長度 | 影片是否是2d or 3d | 影片解析度 | 是否支援字幕 | 發佈時間 | 觀看數 |  讚數 |  倒讚數 |  評論數 |
| 儲存類型  | TEXT | TEXT | TEXT | TEXT | TEXT | TEXT | info | TEXT | TEXT | TEXT | TEXT | TEXT | TEXT | TEXT | INTEGER | INTEGER | INTEGER | INTEGER | INTEGER | INTEGER |









### InfluxDB
- channel

|      | 頻道ID | 總觀看數 | 訂閱者數量 | 總影片數 |
| ---  | ---  | --- | --- | --- |
| key |  channel_id | view_count | subscriber | video_count | 
| type |  tag | field | field | field | 

- video

|      | 頻道ID | 影片ID | 觀看數 | 讚數 | 倒讚數 | 評論數量 |
| ---  | ---  | --- | --- | --- | --- | --- | --- |
| key |  channel_id | video_id | view | like | dislike | comment |
| type |  tag | tag | field | field | field | field |

