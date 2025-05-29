package file

import (
	"crypto/sha1"
	"fmt"
	"io"
	"os"

	"github.com/HugoSmits86/nativewebp"
	"github.com/disintegration/imaging"
)

type ThumbnailProcessor struct {
	size int
}

func NewThumbnailProcessor(size int) ThumbnailProcessor {
	return ThumbnailProcessor{
		size: size,
	}
}

func (t ThumbnailProcessor) Process(db FileDB, m *Meta) error {
	fmt.Println("thumb process", t.size)
	r, err := db.p.Get(m.Hash)
	if err != nil {
		return err
	}
	defer r.Close()

	x, err := imaging.Decode(r)
	if err != nil {
		return err
	}

	x = imaging.Thumbnail(x, t.size, t.size, imaging.Lanczos)
	// x = imaging.CropAnchor(x, t.size, t.size, imaging.Center)

	tmpFile, err := os.CreateTemp("./files/tmp", "myapp-*.txt")
	if err != nil {
		return err
	}
	defer tmpFile.Close()

	h := sha1.New()
	var wc writeCounter
	mw := io.MultiWriter(tmpFile, h, wc)

	// err = imaging.Encode(mw, x, imaging.JPEG, imaging.JPEGQuality(90))
	// if err != nil {
	// 	return err
	// }

	err = nativewebp.Encode(mw, x, nil)
	if err != nil {
		return err
	}

	hash := fmt.Sprintf("%x", h.Sum(nil))
	fmt.Println("moving", tmpFile.Name(), "to", db.Path(hash))
	err = os.Rename(tmpFile.Name(), db.Path(hash))
	if err != nil {
		return err
	}

	db.m.StoreThumbnail(hash, 123, m)

	return nil
}
