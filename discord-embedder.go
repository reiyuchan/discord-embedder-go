package discordvideoembedder

import (
	"bytes"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/url"
	"os"
	"regexp"
	"time"
)

const (
	baseURL   = "https://discord.nfp.is/"
	catboxURL = "https://catbox.moe/user/api.php"
)

type DiscordEmbedder struct {
	client *http.Client
}

func New(client *http.Client) *DiscordEmbedder {
	if client == nil {
		return &DiscordEmbedder{client: &http.Client{Timeout: time.Second * 30}}
	}
	return &DiscordEmbedder{client: client}
}

func (de *DiscordEmbedder) UploadToCatBox(path string) (interface{}, error) {
	_, err := os.Stat(path)
	if err != nil {
		return nil, err
	}
	var buf bytes.Buffer
	writer := multipart.NewWriter(&buf)
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	_ = writer.WriteField("reqtype", "fileupload")
	wfile, err := writer.CreateFormFile("fileToUpload", path)
	if err != nil {
		return nil, err
	}
	io.Copy(wfile, file)
	writer.Close()
	req, err := http.NewRequest(http.MethodPost, catboxURL, &buf)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())
	res, err := de.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	data, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	return string(data), err
}

func (de *DiscordEmbedder) GetURL(videoURL string) (interface{}, error) {
	paURL, err := url.ParseRequestURI(videoURL)
	if err != nil {
		return nil, err
	}
	var buf bytes.Buffer
	writer := multipart.NewWriter(&buf)
	_ = writer.WriteField("video", paURL.String())
	writer.Close()
	req, err := http.NewRequest(http.MethodPost, baseURL, &buf)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())
	res, err := de.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	data, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	regex := regexp.MustCompile("<pre>(.*)</pre>")
	match := regex.FindStringSubmatch(string(data))
	if len(match) < 1 {
		return nil, fmt.Errorf("no match found")
	}
	return match[1], nil
}
