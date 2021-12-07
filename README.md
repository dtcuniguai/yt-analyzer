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


### 尚未完成功能

- youtube api error handle(優先：部分channel權限未開放會拿到403無法取得資料)
- server hosting
- 部分功能參數化
    - channels read file
    - more output file type support
