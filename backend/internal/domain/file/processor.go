package file

type processor interface {
	Process(FileDB, *File) error
}
