package lib

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
)

type steamRespGetPublishedFileDetails struct {
	Response struct {
		PublishedFileDetails []struct {
			Title      string `json:"title"`
			PreviewURL string `json:"preview_url"`
		} `json:"publishedfiledetails"`
	} `json:"response"`
}

func FetchSteamTitle(id string) (string, error) {
	body := bytes.NewBufferString(
		fmt.Sprintf("itemcount=1&publishedfileids[0]=%s", id),
	)

	req, err := http.NewRequest("POST",
		"https://api.steampowered.com/ISteamRemoteStorage/GetPublishedFileDetails/v1/",
		body)
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	var sr steamRespGetPublishedFileDetails
	if err := json.NewDecoder(resp.Body).Decode(&sr); err != nil {
		return "", err
	}

	if len(sr.Response.PublishedFileDetails) == 0 {
		return "", fmt.Errorf("not found")
	}

	return sr.Response.PublishedFileDetails[0].Title, nil
}

func FetchSteamCoverURL(id string) (string, error) {
	body := bytes.NewBufferString(
		fmt.Sprintf("itemcount=1&publishedfileids[0]=%s", id),
	)

	req, _ := http.NewRequest(
		"POST",
		"https://api.steampowered.com/ISteamRemoteStorage/GetPublishedFileDetails/v1/",
		body,
	)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	var sr steamRespGetPublishedFileDetails
	if err := json.NewDecoder(resp.Body).Decode(&sr); err != nil {
		return "", err
	}

	if len(sr.Response.PublishedFileDetails) == 0 {
		return "", fmt.Errorf("not found")
	}

	url := sr.Response.PublishedFileDetails[0].PreviewURL
	if url == "" {
		return "", fmt.Errorf("no preview")
	}

	return url, nil
}

func ModCoverCachePath(appID, folder string) (string, error) {
	base, err := os.UserCacheDir()
	if err != nil {
		return "", err
	}

	dir := filepath.Join(
		base,
		appID,
		folder,
	)

	if err := os.MkdirAll(dir, 0755); err != nil {
		return "", err
	}

	return filepath.Join(dir, "cover"), nil
}

func GetOrDownloadCover(appID string, mod ModEntry) (string, error) {
	cachePath, err := ModCoverCachePath(appID, mod.Folder)
	if err != nil {
		return "", err
	}

	if _, err := os.Stat(cachePath); err == nil {
		return cachePath, nil
	}

	url, err := FetchSteamCoverURL(mod.Folder)
	if err != nil {
		return "", err
	}

	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return "", fmt.Errorf("bad status")
	}

	out, err := os.Create(cachePath)
	if err != nil {
		return "", err
	}
	defer out.Close()

	_, err = io.Copy(out, resp.Body)
	return cachePath, err
}
