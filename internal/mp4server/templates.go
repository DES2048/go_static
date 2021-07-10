package mp4server

import (
	"net/http"
	"path"
	"sort"

	"github.com/dustin/go-humanize"
)

func RenderMp4ListTemplate(config *ServerConfig, response http.ResponseWriter, sorting string) error {
	videos, err := ListMp4(config.StaticFolder)

	if err != nil {
		return err
	}

	// sort by modTime
	sort.Sort(
		sort.Reverse(Mp4ByModTime(videos)),
	)

	videosTmpl := make([]VideoData, 0, len(videos))

	for _, v := range videos {
		videosTmpl = append(videosTmpl,
			VideoData{
				Url:      path.Join(config.StaticPath, v.Name()),
				Title:    v.Name(),
				Size:     humanize.Bytes(uint64(v.Size())),
				Time:     humanize.Time(v.ModTime()),
				Duration: v.Duration.String(),
			},
		)
	}

	data := &struct {
		Videos []VideoData
	}{
		Videos: videosTmpl,
	}

	err = config.ListTemplate.Execute(response, data)
	if err != nil {
		return err
	}

	return nil
}
