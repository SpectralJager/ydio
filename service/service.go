package service

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"

	"github.com/kkdai/youtube/v2"
)

type DownloadAudio struct {
	client *youtube.Client
}

func NewDownloadAudioService() *DownloadAudio {
	return &DownloadAudio{
		client: &youtube.Client{
			MaxRoutines: 2,
			HTTPClient:  http.DefaultClient,
		},
	}
}

func (serv *DownloadAudio) GetAudioMetadate(url string) (*youtube.Video, error) {
	meta, err := serv.client.GetVideo(url)
	if err != nil {
		return nil, err
	}
	formats := meta.Formats.WithAudioChannels()
	formats = formats.Select(func(f youtube.Format) bool {
		return strings.Contains(f.MimeType, "audio/webm")
	})
	meta.Formats = formats
	return meta, nil
}

func (serv *DownloadAudio) DownloadAudio(audio *youtube.Video) (string, error) {
	stream, _, err := serv.client.GetStream(audio, &audio.Formats[0])
	if err != nil {
		return "", err
	}
	defer stream.Close()

	path := fmt.Sprintf("public/audio/%s.mp3", audio.ID)

	file, err := os.Create(path)
	if err != nil {
		return "", err
	}
	defer file.Close()

	_, err = io.Copy(file, stream)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("/static/audio/%s.mp3", audio.ID), nil
}
