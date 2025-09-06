package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"golang.org/x/net/html"
)

type AkkMensaApi interface {
	GetAvailableDates() ([]time.Time, error)
	GetMenuForDate(date time.Time) (map[string]interface{}, error)
}

type AkkMensaApiImpl struct {
	baseUrl string
}

func NewAkkMensaApi(baseUrl string) AkkMensaApi {
	return &AkkMensaApiImpl{
		baseUrl: baseUrl,
	}
}

func (api *AkkMensaApiImpl) GetAvailableDates() ([]time.Time, error) {
	body, err := makeHttpGetRequest(api.baseUrl + "/")
	if err != nil {
		return []time.Time{}, err
	}

	doc, err := html.Parse(strings.NewReader(string(body)))
	if err != nil {
		return []time.Time{}, fmt.Errorf("error parsing HTML: %v", err)
	}

	files := extractFiles(doc)

	var dates []time.Time
	for _, file := range files {
		split := strings.Split(file, ".")
		date, err := time.Parse("2006-01-02", split[0])
		if err != nil {
			return []time.Time{}, fmt.Errorf("error parsing date: %v", err)
		}
		dates = append(dates, date)
	}

	return dates, nil
}

func extractFiles(n *html.Node) []string {
	var files []string

	if n.Type == html.ElementNode && n.Data == "a" {
		for _, attr := range n.Attr {
			if attr.Key == "href" && !strings.HasPrefix(attr.Val, "?") && !strings.HasSuffix(attr.Val, "/") {
				files = append(files, attr.Val)
			}
		}
	}
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		files = append(files, extractFiles(c)...)
	}

	return files
}

func (api *AkkMensaApiImpl) GetMenuForDate(date time.Time) (map[string]interface{}, error) {
	body, err := makeHttpGetRequest(fmt.Sprintf("%s/%s.json", api.baseUrl, date.Format("2006-01-02")))
	if err != nil {
		return nil, err
	}

	var result map[string]interface{}
	err = json.Unmarshal(body, &result)
	if err != nil {
		return nil, fmt.Errorf("error parsing JSON: %v", err)
	}

	return result, nil
}

func makeHttpGetRequest(url string) ([]byte, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("error fetching URL: %v", err)
	}
	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(resp.Body)

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("error while getting json file: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading response body: %v", err)
	}

	return body, nil
}
