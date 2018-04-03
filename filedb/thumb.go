package filedb

import (
	"bytes"
	"image"
	"io"

	"github.com/bakape/thumbnailer"

	"willnorris.com/go/imageproxy"
)

type Thumbnail struct {
	MetaID int
	Hash   string
	Width  int
	Height int
	Size   int

	Data io.Reader
}

func generateThumbnail(r io.Reader) (Thumbnail, error) {
	b := new(bytes.Buffer)
	_, err := b.ReadFrom(r)
	if err != nil {
		return Thumbnail{}, err
	}

	_, thumb, err := thumbnailer.ProcessBuffer(
		b.Bytes(),
		thumbnailer.Options{
			JPEGQuality: 100,
			ThumbDims: thumbnailer.Dims{
				Width:  350,
				Height: 350,
			},
		},
	)

	if err != nil {
		return Thumbnail{}, err
	}

	dat, err := imageproxy.Transform(
		thumb.Data,
		// b.Bytes(),
		imageproxy.Options{
			Width:     200,
			Height:    200,
			Quality:   76,
			SmartCrop: true,
		},
	)

	if err != nil {
		return Thumbnail{}, err
	}

	b2 := bytes.NewBuffer(dat)
	hash, err := calcHash(b2)
	if err != nil {
		return Thumbnail{}, err
	}

	b2 = bytes.NewBuffer(dat)
	img, _, err := image.DecodeConfig(b2)

	b2 = bytes.NewBuffer(dat)

	if err != nil {
		return Thumbnail{}, err
	}

	t := Thumbnail{
		Hash:   hash,
		Size:   len(dat),
		Data:   b2,
		Width:  img.Width,
		Height: img.Height,
	}

	return t, nil
}
