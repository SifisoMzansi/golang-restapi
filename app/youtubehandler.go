package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"google.golang.org/api/option"
	"google.golang.org/api/youtube/v3"
)

type YoutubeStats struct {
	Subscribers    int    `json:"subscribers"`
	Views          int    `json:"views"`
	MinutesWatched int    `json:"minutes_watched"`
	ChannelName    string `json:"channel_name"`
}

//var k = "AIzaSyBujADWwDBFO2hAd16A6GKwbGFh_WmiDQ0"
// channel id  UCAsarZPd1ULXKqOHnGLdmXw

func getChannelStats(apiKey string, channelId string) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {

		ctx := context.Background()
		yts, err := youtube.NewService(ctx, option.WithAPIKey(apiKey))

		if err != nil {
			fmt.Println("Error creating new Youtube service: ", err)
			w.WriteHeader(http.StatusBadRequest)
			return

		}

		call := yts.Channels.List([]string{"snippet, contentDetails, statistics"})

		response, err := call.Id(channelId).Do()

		if err != nil {
			fmt.Println("Error getting channel stats:", err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		fmt.Println("Title ", response.Items[0].Snippet.Title)
		fmt.Println("Description ", response.Items[0].Snippet.Description)
		fmt.Println("PublishedAt ", response.Items[0].Snippet.PublishedAt)
		yt := YoutubeStats{}

		if len(response.Items) > 0 {
			yt = YoutubeStats{
				Subscribers:    int(response.Items[0].Statistics.SubscriberCount),
				Views:          int(response.Items[0].Statistics.ViewCount),
				MinutesWatched: int(response.Items[0].Statistics.VideoCount),
				ChannelName:    response.Items[0].Snippet.Title,
			}
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			if err := json.NewEncoder(w).Encode(yt); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				w.WriteHeader(http.StatusBadRequest)
				return
			}
		} else {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		if err := json.NewEncoder(w).Encode(yt); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
	}
}

//go run ./app/./...
