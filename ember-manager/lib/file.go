package lib

import (
	"code.google.com/p/go.exp/fsnotify"
)

type File struct {
	Path    string
	Content []byte
	Event   *fsnotify.FileEvent
}

func (f *File) IsEmpty() bool {
	return len(f.Content) == 0
}
