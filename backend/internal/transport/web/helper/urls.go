package helper

import (
	"strconv"

	"github.com/a-h/templ"
	"github.com/zqz/web/backend/filedb"
	"github.com/zqz/web/backend/userdb"
)

func URLFileEdit(f *filedb.Meta) templ.SafeURL {
	return templ.URL("/admin/files/" + f.Slug + "/edit")
}

func URLFileView(f *filedb.Meta) templ.SafeURL {
	return templ.URL("/files/" + f.Slug)
}

func URLFileData(f *filedb.Meta) templ.SafeURL {
	return templ.URL("/api/file/by-hash/" + f.Hash)
}

func URLFileProcess(f *filedb.Meta) templ.SafeURL {
	return templ.URL("/admin/files/" + f.Slug + "/process")
}

func URLThumbnailData(f *filedb.Meta) templ.SafeURL {
	return templ.URL("/api/file/by-slug/" + f.Slug + "/thumbnail")
}

func URLUserView(u *userdb.User) templ.SafeURL {
	return templ.URL("/admin/users/" + strconv.Itoa(u.ID))
}
