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

const thumbIntermediateSize = 350
const thumbSize = 200
const thumbQuality = 76

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
				Width:  thumbIntermediateSize,
				Height: thumbIntermediateSize,
			},
		},
	)

	if err != nil {
		return Thumbnail{}, err
	}

	dat, err := imageproxy.Transform(
		thumb.Data,
		imageproxy.Options{
			Width:     thumbSize,
			Height:    thumbSize,
			Quality:   thumbQuality,
			SmartCrop: true,
		},
	)

	if err != nil {
		return Thumbnail{}, err
	}

	return buildThumbnail(dat)
}

func buildThumbnail(dat []byte) (Thumbnail, error) {
	b := bytes.NewReader(dat)
	hash, err := calcHash(b)
	if err != nil {
		return Thumbnail{}, err
	}

	b.Seek(0, 0)
	img, _, err := image.DecodeConfig(b)
	if err != nil {
		return Thumbnail{}, err
	}
	b.Seek(0, 0)

	return Thumbnail{
		Hash:   hash,
		Size:   b.Len(),
		Data:   b,
		Width:  img.Width,
		Height: img.Height,
	}, nil
}
