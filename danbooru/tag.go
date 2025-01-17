package danbooru

import (
	"bytes"
	"encoding/json"
	"fmt"
	"image"
	"net/url"
	"sort"

	"github.com/FloatTech/zbputils/file"
	"github.com/FloatTech/zbputils/txt2img" // jpg png gif
	"github.com/FloatTech/zbputils/web"
	"github.com/fogleman/gg"
	_ "golang.org/x/image/webp"
)

const api = "https://sayuri.fumiama.top/file?path="

type sorttags struct {
	tags map[string]float64
	tseq []string
}

func newsorttags(tags map[string]float64) (s *sorttags) {
	t := make([]string, 0, len(tags))
	for k := range tags {
		t = append(t, k)
	}
	return &sorttags{tags: tags, tseq: t}
}

func (s *sorttags) Len() int {
	return len(s.tags)
}

func (s *sorttags) Less(i, j int) bool {
	v1 := s.tseq[i]
	v2 := s.tseq[j]
	return s.tags[v1] >= s.tags[v2]
}

// Swap swaps the elements with indexes i and j.
func (s *sorttags) Swap(i, j int) {
	s.tseq[j], s.tseq[i] = s.tseq[i], s.tseq[j]
}

func TagURL(name, u string) (t txt2img.TxtCanvas, err error) {
	ch := make(chan []byte, 1)
	go func() {
		var data []byte
		data, err = web.GetData(u)
		ch <- data
	}()

	data, err := web.GetData(api + url.QueryEscape(u))
	if err != nil {
		return
	}
	tags := make(map[string]float64)
	err = json.Unmarshal(data, &tags)
	if err != nil {
		return
	}

	longestlen := 0
	for k := range tags {
		if len(k) > longestlen {
			longestlen = len(k)
		}
	}
	longestlen++

	st := newsorttags(tags)
	sort.Sort(st)

	_, err = file.GetLazyData(txt2img.BoldFontFile, false, true)
	if err != nil {
		return
	}
	_, err = file.GetLazyData(txt2img.ConsolasFontFile, false, true)
	if err != nil {
		return
	}

	data = <-ch
	if err != nil {
		return
	}
	img, _, err := image.Decode(bytes.NewReader(data))
	if err != nil {
		return
	}

	canvas := gg.NewContext(img.Bounds().Size().X, img.Bounds().Size().Y+int(float64(img.Bounds().Size().X)*0.2)+len(tags)*img.Bounds().Size().X/25)
	canvas.SetRGB(1, 1, 1)
	canvas.Clear()
	canvas.DrawImage(img, 0, 0)
	if err = canvas.LoadFontFace(txt2img.BoldFontFile, float64(img.Bounds().Size().X)*0.1); err != nil {
		return
	}
	canvas.SetRGB(0, 0, 0)
	canvas.DrawString(name, float64(img.Bounds().Size().X)*0.02, float64(img.Bounds().Size().Y)+float64(img.Bounds().Size().X)*0.1)
	i := float64(img.Bounds().Size().Y) + float64(img.Bounds().Size().X)*0.2
	if err = canvas.LoadFontFace(txt2img.ConsolasFontFile, float64(img.Bounds().Size().X)*0.04); err != nil {
		return
	}
	rate := float64(img.Bounds().Size().X) * 0.04
	for _, k := range st.tseq {
		canvas.DrawString(fmt.Sprintf("* %-*s -%.3f-", longestlen, k, tags[k]), float64(img.Bounds().Size().X)*0.04, i)
		i += rate
	}
	t.Canvas = canvas
	return
}
