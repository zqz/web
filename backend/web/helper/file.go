package helper

import (
	"strconv"

	"github.com/a-h/templ"
	"github.com/zqz/web/backend/filedb"
	"github.com/zqz/web/backend/userdb"
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

func URLEditFile(f *filedb.Meta) templ.SafeURL {
	return templ.URL("/admin/files/" + f.Slug + "/edit")
}

func URLViewFile(f *filedb.Meta) templ.SafeURL {
	return templ.URL("/files/" + f.Slug)
}

func URLFileData(f *filedb.Meta) templ.SafeURL {
	return templ.URL("/api/file/by-hash/" + f.Hash)
}

func URLThumbnailData(f *filedb.Meta) templ.SafeURL {
	return templ.URL("/api/file/by-slug/" + f.Slug + "/thumbnail")
}

func URLProcessFile(f *filedb.Meta) templ.SafeURL {
	return templ.URL("/admin/files/" + f.Slug + "/process")
}

func URLViewUser(u *userdb.User) templ.SafeURL {
	return templ.URL("/admin/users/" + strconv.Itoa(u.ID))
}
