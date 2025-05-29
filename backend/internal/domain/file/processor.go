package file

type processor interface {
	Process(DB, *File) error
}
