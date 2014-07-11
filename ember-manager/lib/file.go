package lib

import (
	"code.google.com/p/go.exp/fsnotify"
)

type File struct {
	Path    string
	Content []byte
	Event   *fsnotify.FileEvent
}
