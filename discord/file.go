package discord

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path"
)

// File represents a file that can be sent with Send and the WithFiles option.
type File struct {
	Name   string
	Reader io.ReadCloser
}

// FileFromReadCloser returns a File given a ReadCloser and a name.
// If the name ends with a valid extension recognized by Discord's
// client applications, the file will be displayed inline in the channel
// instead of asking users to manually download it.
func FileFromReadCloser(r io.ReadCloser, name string) *File {
	return &File{
		Name:   name,
		Reader: r,
	}
}

// FileFromDisk returns a File from a local, on disk file.
// If name is left empty, it will default to the name of the
// file on disk.
// Note that since Send is responsible for closing files opened
// by FileFromDisk, calling this function and *not* calling Send
// after can lead to resource leaks.
func FileFromDisk(filepath, name string) (*File, error) {
	f, err := os.Open(filepath)
	if err != nil {
		return nil, err
	}

	if name == "" {
		name = path.Base(filepath)
	}

	return FileFromReadCloser(f, name), nil
}

// FileFromURL returns a File from a remote HTTP resource.
// If the name is left empty, it will default to the name of
// the file specified in the URL.
// Note that since Send is responsible for closing files opened
// by FileFromURL, calling this function and *not* calling Send
// after can lead to resource leaks.
func FileFromURL(url, name string) (*File, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}

	ct := resp.Header.Get("Content-Type")
	if !stringsContains(supportedContentTypes, ct) {
		return nil, fmt.Errorf("unsupported Content-Type: %q", ct)
	}

	if name == "" {
		name = path.Base(url)
	}

	return FileFromReadCloser(resp.Body, name), nil
}

var supportedContentTypes = []string{
	"image/png",
	"image/gif",
	"image/jpeg",
	"video/mp4",
}

func stringsContains(stack []string, needle string) bool {
	for _, hay := range stack {
		if needle == hay {
			return true
		}
	}
	return false
}
