package filedb

import (
	"bytes"
	"io"

	"github.com/bakape/thumbnailer"

	"willnorris.com/go/imageproxy"
)

type Thumbnail struct {
	MetaHash string
	Hash     string
	Width    int
	Height   int
	Size     int

	Data io.Reader
}

func GenThumbnail(r io.Reader, h string) (Thumbnail, error) {

	b := new(bytes.Buffer)
	b.ReadFrom(r)

	// disabled because it panics with weird cgo issues.
	// uncomment at later date, it handles thumbnailing
	// of videos, pdfs and more text formats.
	_, thumb, err := thumbnailer.ProcessBuffer(
		b.Bytes(),
		thumbnailer.Options{
			JPEGQuality: 100,
			ThumbDims: thumbnailer.Dims{
				Width:  150,
				Height: 150,
			},
		},
	)

	if err != nil {
		return Thumbnail{}, err
	}

	dat, err := imageproxy.Transform(
		thumb.Data,
		imageproxy.Options{
			Width:     150,
			Height:    150,
			Quality:   70,
			SmartCrop: true,
		},
	)

	if err != nil {
		return Thumbnail{}, err
	}

	b = bytes.NewBuffer(dat)
	hash, err := calcHash(b)
	b = bytes.NewBuffer(dat)

	if err != nil {
		return Thumbnail{}, err
	}

	t := Thumbnail{
		MetaHash: h,
		Hash:     hash,
		Size:     len(dat),
		Data:     b,
	}

	return t, nil
}
