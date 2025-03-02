package discordvideoembedder

import (
	"bytes"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"regexp"
	"strings"
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

// UploadToCatBox returns URL of uploaded file
func (de *DiscordEmbedder) UploadToCatBox(path string) (string, error) {
	_, err := os.Stat(path)
	if err != nil {
		return "", err
	}
	var buf bytes.Buffer
	writer := multipart.NewWriter(&buf)
	file, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer file.Close()
	_ = writer.WriteField("reqtype", "fileupload")
	wfile, err := writer.CreateFormFile("fileToUpload", path)
	if err != nil {
		return "", err
	}
	io.Copy(wfile, file)
	writer.Close()
	req, err := http.NewRequest(http.MethodPost, catboxURL, &buf)
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())
	res, err := de.client.Do(req)
	if err != nil {
		return "", err
	}
	defer res.Body.Close()
	data, err := io.ReadAll(res.Body)
	if err != nil {
		return "", err
	}
	return string(data), err
}

// GetURL returns the generated embed URL
func (de *DiscordEmbedder) GetURL(videoURL string) (string, error) {
	paURL, err := url.ParseRequestURI(videoURL)
	if err != nil {
		return "", err
	}
	ext := strings.ToLower(filepath.Ext(paURL.String()))
	filter := map[string]bool{
		".mp4":  true,
		".avi":  true,
		".mov":  true,
		".wmv":  true,
		".flv":  true,
		".webm": true,
	}
	if !filter[ext] || ext == "" {
		return "", fmt.Errorf("file extension not supported")
	}
	var buf bytes.Buffer
	writer := multipart.NewWriter(&buf)
	_ = writer.WriteField("video", paURL.String())
	writer.Close()
	req, err := http.NewRequest(http.MethodPost, baseURL, &buf)
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())
	res, err := de.client.Do(req)
	if err != nil {
		return "", err
	}
	defer res.Body.Close()
	data, err := io.ReadAll(res.Body)
	if err != nil {
		return "", err
	}
	regex := regexp.MustCompile("<pre>(.*)</pre>")
	match := regex.FindStringSubmatch(string(data))
	if len(match) <= 1 {
		return "", fmt.Errorf("no match found due to request failure")
	}
	return match[1], nil
}
