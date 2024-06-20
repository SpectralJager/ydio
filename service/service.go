package service

import (
	"archive/zip"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"regexp"

	"github.com/kkdai/youtube/v2"
)

var (
	rgxp = regexp.MustCompile(`[^a-zA-Z0-9\.,;:&\[\]'" ]+`)
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
	data, _ := json.MarshalIndent(meta.Formats, "", "  ")
	file, _ := os.Create("formats.json")
	file.Write(data)
	formats := meta.Formats.WithAudioChannels()
	formats = formats.Select(func(f youtube.Format) bool {
		return f.MimeType == "audio/webm; codecs=\"opus\"" &&
			f.AudioQuality == "AUDIO_QUALITY_MEDIUM"
	})
	meta.Formats = formats
	return meta, nil
}

func (serv *DownloadAudio) GetPlaylistMetadate(url string) (*youtube.Playlist, error) {
	return serv.client.GetPlaylist(url)
}

func (serv *DownloadAudio) DownloadAudio(audio *youtube.Video) error {
	path := fmt.Sprintf("public/audio/%s.mp3", audio.ID)
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()

	return serv.downloadAudio(audio, file)

}

func (serv *DownloadAudio) DownloadPlaylist(playlist *youtube.Playlist, ids []string) error {
	if len(ids) != 0 {
		videos := []*youtube.PlaylistEntry{}
		for _, id := range ids {
			for _, entry := range playlist.Videos {
				if id == entry.ID {
					videos = append(videos, entry)
					break
				}
			}
		}
		playlist.Videos = videos
	}
	zipFile, err := os.Create(fmt.Sprintf("./public/audio/%s.zip", playlist.ID))
	if err != nil {
		return err
	}
	zipWriter := zip.NewWriter(zipFile)
	defer zipWriter.Close()
	counter := 0
	for _, video := range playlist.Videos {
		video, err := serv.GetAudioMetadate(video.ID)
		if err != nil {
			continue
		}
		zipFile, err := zipWriter.Create(fmt.Sprintf("%s.mp3", rgxp.ReplaceAllString(video.Title, " ")))
		if err != nil {
			continue
		}
		err = serv.downloadAudio(video, zipFile)
		if err != nil {
			continue
		}
		log.Println(video.Title)
		counter += 1
	}
	if counter == 0 {
		return fmt.Errorf("can't download audios for playlist %s", playlist.ID)
	}
	return nil
}

func (serv *DownloadAudio) downloadAudio(audio *youtube.Video, writer io.Writer) error {
	if len(audio.Formats) == 0 {
		return fmt.Errorf("no audio format for '%s'", audio.ID)
	}
	stream, _, err := serv.client.GetStream(audio, &audio.Formats[0])
	if err != nil {
		return err
	}
	defer stream.Close()

	_, err = io.Copy(writer, stream)
	if err != nil {
		return err
	}
	return nil
}
