package fs

type IEntry interface {
	GetName() string
	GetPath() string
	GetSize() int64
	IsDirectory() bool
}

type Entry struct {
	IEntry

	name string
	path string
	size int64
}

func (e *Entry) GetName() string {
	return e.name
}

func (e *Entry) GetPath() string {
	return e.path
}

func (e *Entry) GetSize() int64 {
	return e.size
}

type File struct {
	*Entry
}

func (*File) IsDirectory() bool {
	return false
}

type Directory struct {
	*Entry
}

func (*Directory) IsDirectory() bool {
	return true
}
