package helper

import (
	"github.com/a-h/templ"
	"github.com/zqz/web/backend/filedb"
)

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
	return templ.URL("/admin/files/" + f.Slug)
}
