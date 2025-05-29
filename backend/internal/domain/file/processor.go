package file

type processor interface {
	Process(FileDB, *Meta) error
}
