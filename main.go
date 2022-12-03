package main

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
	"os/signal"
	"strings"
	"time"

	"ytanalyzer/app/cron"
	"ytanalyzer/lib/youtube"
	"ytanalyzer/route"

	"github.com/go-co-op/gocron"
	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
	"github.com/urfave/cli"
)

type (
	// Config information.
	Config struct {
		channels string
		output   string
		source   string
	}
)

var config Config

func main() {

	godotenv.Load()
	s := gocron.NewScheduler(time.UTC)
	s.Every(5).Minutes().Do(cron.TaskSyncYT)
	s.StartAsync()

	app := fiber.New()
	//shutdown gracefully
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		_ = <-c
		fmt.Println("app shut down ing...")
		s.Remove(cron.TaskSyncYT)
		_ = app.Shutdown()
	}()

	route.RegisterRoute(app)
	app.Listen(fmt.Sprintf(":%v", os.Getenv("APP_PORT")))
}

func cmd() {

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
		cli.StringFlag{
			Name:  "source,s",
			Usage: "channels file read",
		},
	}

	app.Run(os.Args)
}

func run(c *cli.Context) error {
	config = Config{
		channels: c.String("channels"),
		output:   c.String("output"),
		source:   c.String("source"),
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
	channels, err := readChannels(config)
	if err != nil {
		log.Fatal(err.Error())
	}

	infos := make(map[string]youtube.ChannelInfo)

	for _, channel := range *channels {

		var info youtube.ChannelInfo

		// 訂閱資料
		subDetail, subErr := youtube.GetSubscriptionsDetail(ytApiKey, channel)
		if subErr != nil {
			fmt.Println(subErr.Error())
		} else {
			info.ID = subDetail.Items[0].SubscriberSnippet.ChannelID
			info.Title = subDetail.Items[0].SubscriberSnippet.Title
			info.Description = subDetail.Items[0].SubscriberSnippet.Description
			info.Thumbnails = map[string]string{
				"default": subDetail.Items[0].SubscriberSnippet.Thumbnails.Default.URL,
				"medium":  subDetail.Items[0].SubscriberSnippet.Thumbnails.Medium.URL,
				"high":    subDetail.Items[0].SubscriberSnippet.Thumbnails.High.URL,
			}
		}

		//頻道內容
		chDetail, chErr := youtube.GetChannelDetail(ytApiKey, channel)
		if chErr != nil {
			fmt.Println(chErr.Error())
		} else {
			info.SubscriberCount = chDetail.Items[0].Statistics.SubscriberCount
			info.VideoCount = chDetail.Items[0].Statistics.VideoCount
			info.ViewCount = chDetail.Items[0].Statistics.ViewCount
		}

		//channel all info skip
		if chErr != nil && subErr != nil {
			fmt.Println("channel : " + channel + " skip")
			continue
		}

		if subDetail != nil {
			infos[subDetail.Items[0].SubscriberSnippet.Title] = info
		} else {
			infos[channel] = info
		}
	}
	outputData(infos, config.output)

	return nil
}

func readChannels(config Config) (*[]string, error) {

	if len(config.channels) != 0 {
		channels := strings.Split(config.channels, ",")
		return &channels, nil
	} else if len(config.source) != 0 {
		file, err := os.Open(config.source)
		if err != nil {
			return nil, err
		}
		defer file.Close()

		var channels []string
		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			channels = append(channels, scanner.Text())
		}

		return &channels, scanner.Err()
	} else {
		return nil, errors.New("missing channel's given")
	}
}

//
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
