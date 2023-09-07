package storage

import (
	"context"
	"os"
	"path"
)

type Local string

func (d Local) Mkdir(ctx context.Context, name string, perm os.FileMode) error {
	realName := path.Join(string(d), name)
	return os.Mkdir(realName, perm)
}

func (d Local) OpenFile(ctx context.Context, name string, flag int, perm os.FileMode) (File, error) {
	realName := path.Join(string(d), name)
	return os.OpenFile(realName, flag, perm)
}

func (d Local) RemoveAll(ctx context.Context, name string) error {
	realName := path.Join(string(d), name)
	return os.RemoveAll(realName)
}

func (d Local) Rename(ctx context.Context, oldName, newName string) error {
	realOldName := path.Join(string(d), oldName)
	realNewName := path.Join(string(d), newName)
	return os.Rename(realOldName, realNewName)
}

func (d Local) Stat(ctx context.Context, name string) (os.FileInfo, error) {
	realName := path.Join(string(d), name)
	return os.Stat(realName)
}
