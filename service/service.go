package service

import (
	"archive/zip"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/kkdai/youtube/v2"
)

var (
// rgxp = regexp.MustCompile(`[]+`)
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
		return f.MimeType == "audio/webm; codecs=\"opus\"" ||
			f.MimeType == "audio/mp4; codecs=\"mp4a.40.2\""
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
	videos := playlist.Videos[:]
	zipFile, err := os.Create(fmt.Sprintf("./public/audio/%s.zip", playlist.ID))
	if err != nil {
		return err
	}
	zipWriter := zip.NewWriter(zipFile)
	defer zipWriter.Close()
	counter := 0
	for _, id := range ids {
		for i, entry := range videos {
			if id != entry.ID {
				continue
			}
			video, err := serv.GetAudioMetadate(entry.ID)
			if err != nil {
				log.Println("no meta", entry.Title)
				continue
			}
			zipFile, err := zipWriter.Create(fmt.Sprintf("%s.mp3", video.Title))
			if err != nil {
				log.Println("can't create zip entity", entry.Title)
				continue
			}
			err = serv.downloadAudio(video, zipFile)
			if err != nil {
				log.Println("can't download", entry.Title, err.Error())
				continue
			}
			log.Println(video.Title)
			counter += 1
			if i < len(videos) {
				videos = append(videos[:i], videos[i+1:]...)
			}
		}
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
	counter := 0
	var er error
	for counter < 10 {
		stream, _, err := serv.client.GetStream(audio, &audio.Formats[0])
		if err != nil {
			er = err
			counter += 1
			continue
		}
		defer stream.Close()

		_, err = io.Copy(writer, stream)
		if err != nil {
			er = err
			counter += 1
			continue
		}
		return nil
	}
	return er
}
