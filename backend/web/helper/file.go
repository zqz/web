package helper

import (
	"github.com/a-h/templ"
	"github.com/zqz/web/backend/filedb"
)

func StyleFileBG(f *filedb.Meta) templ.SafeCSS {
	if len(f.Thumbnail) == 0 {
		return templ.SafeCSS("")
	}

	url := URLThumbnailData(f)
	return templ.SafeCSS("background-image: url(" + url + ");")
}

func TitleAdminFile(f *filedb.Meta) string {
	return "a | file | " + f.Name
}

func TitleAdminEditFile(f *filedb.Meta) string {
	return "a | edit file | " + f.Name
}
