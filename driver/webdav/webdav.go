package webdav

import (
	"io"
	"os"
	"strings"

	"github.com/studio-b12/gowebdav"
	"goftp.io/server/v2"
)

type Driver struct {
	client *gowebdav.Client
}

func NewDriver(url string, user string, password string, useSSL bool) (server.Driver, error) {
	c := gowebdav.NewClient(url, user, password)
	return &Driver{
		client: c,
	}, nil
}

func buildWebDAVPath(p string) string {
	return strings.TrimPrefix(p, "/")
}

func (driver *Driver) DeleteDir(ctx *server.Context, path string) error {
	p := buildWebDAVPath(path)
	return driver.client.Remove(p)
}

func (driver *Driver) DeleteFile(ctx *server.Context, path string) error {
	return driver.client.Remove(buildWebDAVPath(path))
}

func (driver *Driver) GetFile(ctx *server.Context, path string, offset int64) (int64, io.ReadCloser, error) {
	p := buildWebDAVPath(path)
	info, err := driver.client.Stat(p)
	if err != nil {
		return 0, nil, err
	}

	var object io.ReadCloser
	length := info.Size() - offset
	if offset == 0 {
		object, err = driver.client.ReadStream(p)
	} else {
		object, err = driver.client.ReadStreamRange(p, offset, length)
	}
	if err != nil {
		return 0, nil, err
	}
	defer func() {
		if err != nil && object != nil {
			object.Close()
		}
	}()

	return length, object, nil
}

func (driver *Driver) ListDir(ctx *server.Context, path string, callback func(os.FileInfo) error) error {
	p := buildWebDAVPath(path)
	if p == "/" {
		p = ""
	}
	objects, err := driver.client.ReadDir(p)
	if err != nil {
		return err
	}
	for _, object := range objects {
		err := callback(object)
		if err != nil {
			return err
		}
	}
	return nil
}

func (driver *Driver) MakeDir(ctx *server.Context, path string) error {
	return driver.client.Mkdir(path, os.ModePerm)
}

func (driver *Driver) PutFile(ctx *server.Context, destPath string, data io.Reader, offset int64) (int64, error) {
	p := buildWebDAVPath(destPath)
	err := driver.client.WriteStream(p, data, 0644)
	return 0, err
}

func (driver *Driver) Rename(ctx *server.Context, fromPath string, toPath string) error {
	fp := buildWebDAVPath(fromPath)
	tp := buildWebDAVPath(toPath)
	return driver.client.Rename(fp, tp, false)
}

func (driver *Driver) Stat(ctx *server.Context, path string) (os.FileInfo, error) {
	p := buildWebDAVPath(path)
	return driver.client.Stat(p)
}
