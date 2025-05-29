package helper

import (
	"github.com/a-h/templ"
	"github.com/zqz/web/backend/internal/domain/file"
)

func StyleFileBG(f *file.File) templ.SafeCSS {
	if len(f.Thumbnail) == 0 {
		return templ.SafeCSS("")
	}

	url := URLThumbnailData(f)
	return templ.SafeCSS("background-image: url(" + url + ");")
}

func TitleAdminFile(f *file.File) string {
	return "a | file | " + f.Name
}

func TitleAdminEditFile(f *file.File) string {
	return "a | edit file | " + f.Name
}
