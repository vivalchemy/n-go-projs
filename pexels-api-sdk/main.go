package main

import (
	"encoding/json"
	"fmt"
	"io"
	"math/rand/v2"
	"net/http"
	"os"
	"strconv"
	"time"

	dotenv "github.com/joho/godotenv"
)

const (
	PhotosApiUrl = "https://api.pexels.com/v1"
	VideosApiUrl = "https://api.pexels.com/videos"
)

type Client struct {
	Token          string
	hc             *http.Client
	RemainingTimes int32
}

func NewClient(token string) *Client {
	return &Client{
		Token: token,
		hc: &http.Client{
			Timeout: time.Second * 10,
		},
	}
}

type SearchResult struct {
	Page         int32   `json:"page"`
	PerPage      int32   `json:"per_page"`
	TotalResults int32   `json:"total_results"`
	NextPage     string  `json:"next_page"`
	Photos       []Photo `json:"photos"`
}

type CuratedResult struct {
	Page     int32   `json:"page"`
	PerPage  int32   `json:"per_page"`
	NextPage string  `json:"next_page"`
	Photos   []Photo `json:"photos"`
}

type Photo struct {
	ID              int32       `json:"id"`
	Width           int32       `json:"width"`
	Height          int32       `json:"height"`
	URL             string      `json:"url"`
	Photographer    string      `json:"photographer"`
	PhotographerURL string      `json:"photographer_url"`
	Src             PhotoSource `json:"src"`
}

type PhotoSource struct {
	Original  string `json:"original"`
	Large     string `json:"large"`
	Large2x   string `json:"large2x"`
	Medium    string `json:"medium"`
	Small     string `json:"small"`
	Portrait  string `json:"portrait"` // Fixed typo (was "Potrait")
	Square    string `json:"square"`
	Landscape string `json:"landscape"`
	Tiny      string `json:"tiny"`
}

type VideoSearchResult struct {
	Page         int      `json:"page"`
	PerPage      int      `json:"per_page"`
	TotalResults int      `json:"total_results"`
	NextPage     string   `json:"next_page"`
	Videos       []Videos `json:"videos"`
}

type PopularVideosResult struct {
	Page         int      `json:"page"`
	PerPage      int      `json:"per_page"`
	TotalResults int      `json:"total_results"`
	Url          string   `json:"url"`
	Videos       []Videos `json:"videos"`
}

type Videos struct {
	ID            int             `json:"id"`
	Width         int             `json:"width"`
	Height        int             `json:"height"`
	URL           string          `json:"url"`
	Image         string          `json:"image"`
	Duration      float64         `json:"duration"`
	FullRes       any             `json:"full_res"`
	VideoFiles    []VideoFiles    `json:"video_files"`
	VideoPictures []VideoPictures `json:"video_pictures"`
}

type VideoFiles struct {
	ID       int    `json:"id"`
	Quality  string `json:"quality"`
	FileType string `json:"file_type"`
	Width    int    `json:"width"`
	Height   int    `json:"height"`
	Link     string `json:"link"`
}

type VideoPictures struct {
	ID      int    `json:"id"`
	Picture string `json:"picture"`
	Nr      int    `json:"nr"`
}

func (c *Client) SearchPhotos(query string, perPage int, page int) (*SearchResult, error) {
	url := fmt.Sprintf("%s/search?query=%s&per_page=%d&page=%d", PhotosApiUrl, query, perPage, page)
	res, err := c.requestDoWithAuth("GET", url)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close() // Ensure res is not nil before using

	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API error: %d", res.StatusCode)
	}

	data, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	var result SearchResult
	err = json.Unmarshal(data, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

func (c *Client) CuratedPhotos(perPage, page int) (*CuratedResult, error) {
	url := fmt.Sprintf("%s/curated?per_page=%d&page=%d", PhotosApiUrl, perPage, page)
	res, err := c.requestDoWithAuth("GET", url)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close() // Ensure res is not nil before using

	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API error: %d", res.StatusCode)
	}

	data, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	var result CuratedResult
	err = json.Unmarshal(data, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil

}

func (c *Client) GetPhotos(id int32) (*Photo, error) {
	url := fmt.Sprintf("%s/photos/%d", PhotosApiUrl, id)
	res, err := c.requestDoWithAuth("GET", url)
	if err != nil {
		return nil, err
	}

	data, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	var result Photo
	err = json.Unmarshal(data, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

func (c *Client) GetRandomPhoto() (*Photo, error) {
	randNum := rand.IntN(1001)
	result, err := c.CuratedPhotos(1, randNum)
	if err == nil && len(result.Photos) > 0 {
		return &result.Photos[0], nil
	}
	return nil, err
}

func (c *Client) SearchVideos(query string, perPage int, page int) (*VideoSearchResult, error) {
	url := fmt.Sprintf("%s/search?query=%s&per_page=%d&page=%d", VideosApiUrl, query, perPage, page)
	res, err := c.requestDoWithAuth("GET", url)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close() // Ensure res is not nil before using

	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API error: %d", res.StatusCode)
	}

	data, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	var result VideoSearchResult
	err = json.Unmarshal(data, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

func (c *Client) PopularVideos(perPage int, page int) (*VideoSearchResult, error) {
	url := fmt.Sprintf("%s/search?per_page=%d&page=%d", VideosApiUrl, perPage, page)
	res, err := c.requestDoWithAuth("GET", url)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close() // Ensure res is not nil before using

	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API error: %d", res.StatusCode)
	}

	data, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	var result VideoSearchResult
	err = json.Unmarshal(data, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

func (c *Client) GetVideo(id int) (*Videos, error) {
	url := fmt.Sprintf("%s/videos/%d", VideosApiUrl, id)
	res, err := c.requestDoWithAuth("GET", url)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close() // Ensure res is not nil before using

	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API error: %d", res.StatusCode)
	}

	data, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	var result Videos
	err = json.Unmarshal(data, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

func (c *Client) PopularVideo(perPage int, page int) (*PopularVideosResult, error) {
	url := fmt.Sprintf("%s/search?per_page=%d&page=%d", VideosApiUrl, perPage, page)
	res, err := c.requestDoWithAuth("GET", url)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close() // Ensure res is not nil before using

	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API error: %d", res.StatusCode)
	}

	data, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	var result PopularVideosResult
	err = json.Unmarshal(data, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

func (c *Client) GetRandomVideo() (*Videos, error) {
	randNum := rand.IntN(1001)
	result, err := c.PopularVideos(1, randNum)
	if err == nil && len(result.Videos) > 0 {
		return &result.Videos[0], nil
	}
	return nil, err
}

func (c *Client) requestDoWithAuth(method string, url string) (*http.Response, error) {
	req, err := http.NewRequest(method, url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", c.Token)

	res, err := c.hc.Do(req)
	if err != nil {
		return nil, err
	}

	// Check for rate limit header
	ratelimitStr := res.Header.Get("X-RateLimit-Remaining")
	if ratelimitStr != "" {
		ratelimit, err := strconv.Atoi(ratelimitStr)
		if err == nil {
			c.RemainingTimes = int32(ratelimit)
		}
	}

	return res, nil
}

func main() {
	err := dotenv.Load()
	if err != nil {
		fmt.Println("Error loading .env file:", err)
		return
	}

	TOKEN := os.Getenv("PEXELS_API_TOKEN")
	if TOKEN == "" {
		fmt.Println("PEXELS_API_TOKEN is missing from the environment variables")
		return
	}

	client := NewClient(TOKEN)

	result, err := client.SearchVideos("anime desktop wallpapers", 15, 1)
	if err != nil {
		fmt.Println("Error searching photos:", err)
		return
	}

	if result.Page == 0 {
		fmt.Println("No photos found")
		return
	}

	fmt.Println("Search Results:", result)

	// jsonResult, err := json.Marshal(result)
	// if err != nil {
	// 	fmt.Println("Error marshalling result:", err)
	// 	return
	// }
	// fmt.Println(string(jsonResult))
}
