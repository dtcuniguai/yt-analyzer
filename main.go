package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"
	"time"
	"you-cli/youtube"

	"github.com/joho/godotenv"
	"github.com/urfave/cli"
)

type (
	// Config information.
	Config struct {
		channels string
		output   string
	}
)

var config Config

func main() {

	app := cli.NewApp()
	app.Name = "Youtube-Cli"
	app.Usage = "application youtube command line"
	app.Action = run
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "channels,c",
			Usage: "youtube channel ids you want to get info (support multiple channels separate with \",\"  E,g: you-cli -c 'channel1,channel2,channel3'  )",
		},
		cli.StringFlag{
			Name:  "output,o",
			Usage: "location for output data saving type:json",
		},
	}

	app.Run(os.Args)
}

func run(c *cli.Context) error {
	config = Config{
		channels: c.String("channels"),
		output:   c.String("output"),
	}

	return exec()
}

func exec() error {

	//load env file
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	//get yt api key
	ytApiKey := os.Getenv("GOOGLE_API_KEY")
	channels := strings.Split(config.channels, ",")

	infos := make(map[string]youtube.ChannelInfo)

	for _, channel := range channels {
		// get yt data
		subDetail, err := youtube.GetSubscriptionsDetail(ytApiKey, channel)
		if err != nil {
			fmt.Println(err.Error())
		}
		chDetail, err := youtube.GetChannelDetail(ytApiKey, channel)

		if err != nil {
			fmt.Println(err.Error())
		}

		infos[subDetail.Items[0].SubscriberSnippet.Title] = youtube.ChannelInfo{
			ID:              subDetail.Items[0].SubscriberSnippet.ChannelID,
			Title:           subDetail.Items[0].SubscriberSnippet.Title,
			Description:     subDetail.Items[0].SubscriberSnippet.Description,
			SubscriberCount: chDetail.Items[0].Statistics.SubscriberCount,
			VideoCount:      chDetail.Items[0].Statistics.VideoCount,
			ViewCount:       chDetail.Items[0].Statistics.ViewCount,
			Thumbnails: map[string]string{
				"default": subDetail.Items[0].SubscriberSnippet.Thumbnails.Default.URL,
				"medium":  subDetail.Items[0].SubscriberSnippet.Thumbnails.Medium.URL,
				"high":    subDetail.Items[0].SubscriberSnippet.Thumbnails.High.URL,
			},
		}
	}
	outputData(infos, config.output)

	return nil
}

func outputData(data map[string]youtube.ChannelInfo, dir string) {

	//default directory
	if len(dir) == 0 {
		dir = "./"
	}

	fileName := "yt-" + time.Now().Format("20060102150405") + ".json"

	f, err := os.Create(fileName)

	if err != nil {
		log.Fatal(err)
	}

	defer f.Close()

	jsonString, err := json.Marshal(data)
	if err != nil {
		log.Fatal(err)
	}
	_, err2 := f.WriteString(string(jsonString))

	if err2 != nil {
		log.Fatal(err2)
	}

}
