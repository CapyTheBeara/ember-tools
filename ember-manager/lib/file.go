package lib

import (
	"code.google.com/p/go.exp/fsnotify"
	"io/ioutil"
)

type File struct {
	Path    string
	Content []byte
	Event   *fsnotify.FileEvent
}

func (f *File) IsEmpty() bool {
	return len(f.Content) == 0
}

func (f *File) SetContent() {
	file, err := ioutil.ReadFile(f.Path)
	if err != nil {
		f.Content = []byte{}
	} else {
		f.Content = file
	}
}
