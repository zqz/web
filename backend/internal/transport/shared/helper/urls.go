package helper

import (
	"strconv"

	"github.com/a-h/templ"
	"github.com/zqz/web/backend/internal/domain/file"
	"github.com/zqz/web/backend/internal/domain/user"
)

func URLFileEdit(f *file.Meta) templ.SafeURL {
	return templ.URL("/admin/files/" + f.Slug + "/edit")
}

func URLFileView(f *file.Meta) templ.SafeURL {
	return templ.URL("/files/" + f.Slug)
}

func URLFileData(f *file.Meta) templ.SafeURL {
	return templ.URL("/api/file/by-hash/" + f.Hash)
}

func URLFileProcess(f *file.Meta) templ.SafeURL {
	return templ.URL("/admin/files/" + f.Slug + "/process")
}

func URLThumbnailData(f *file.Meta) templ.SafeURL {
	return templ.URL("/api/file/by-slug/" + f.Slug + "/thumbnail")
}

func URLUserView(u *user.User) templ.SafeURL {
	return templ.URL("/admin/users/" + strconv.Itoa(u.ID))
}
