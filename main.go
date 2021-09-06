package main

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"
)

const PUBLIC = "public"
const PUBLIC_PATH = "/" + PUBLIC + "/"

var CHANNELS = []*channel{}

type channel struct {
	ID   string
	Name string
	Live string
}

func main() {
	go func() {
		for {
			loadChannels()
			time.Sleep(time.Hour)
		}
	}()

	http.Handle(PUBLIC_PATH, http.StripPrefix(PUBLIC_PATH, http.FileServer(http.Dir("./"+PUBLIC))))
	http.HandleFunc("/", indexHandler)
	http.ListenAndServe(":3000", nil)
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "text/html; charset=UTF-8")
	views.ExecuteTemplate(w, "index", CHANNELS)
}

func loadChannels() {
	chans := []*channel{}

	for _, id := range channelsIDs {
		name, err := channelName(id)
		if err != nil {
			log.Print(err)
			continue
		}

		vs, err := liveVideos(id)
		if err != nil {
			log.Print(err)
			continue
		}

		for _, v := range vs {
			chans = append(chans, &channel{
				ID:   id,
				Name: name,
				Live: v,
			})
		}
	}

	CHANNELS = chans
}

func liveVideos(id string) ([]string, error) {
	url := fmt.Sprintf(
		"https://youtube.googleapis.com/youtube/v3/search?channelId=%s&eventType=live&type=video&key=%s",
		id,
		os.Getenv("YOUTUBE_API_KEY"),
	)

	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("GET %s failed: %w", id, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("GET %s response not HTTP OK: %s", id, resp.Status)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("Reading body for %s failed: %w", id, err)
	}

	feed := struct {
		Items []struct {
			ID struct {
				VideoID string `json:"videoId"`
			} `json:"id"`
		} `json:"items"`
	}{}

	if err = json.Unmarshal(body, &feed); err != nil {
		return nil, fmt.Errorf("Unmarshal %s failed: %w", id, err)
	}

	ids := []string{}
	for _, item := range feed.Items {
		ids = append(ids, item.ID.VideoID)
	}

	return ids, nil
}

func channelName(id string) (string, error) {
	url := fmt.Sprintf(
		"https://www.youtube.com/feeds/videos.xml?channel_id=%s",
		id,
	)

	resp, err := http.Get(url)
	if err != nil {
		return "", fmt.Errorf("GET %s failed: %w", id, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("GET %s response not HTTP OK: %s", id, resp.Status)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("Reading body for %s failed: %w", id, err)
	}

	feed := struct {
		Title string `xml:"title"`
	}{}

	if err = xml.Unmarshal(body, &feed); err != nil {
		return "", fmt.Errorf("Unmarshal %s failed: %w", id, err)
	}

	return feed.Title, nil
}
