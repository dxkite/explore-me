package core

import (
	"context"

	"dxkite.cn/explorer/src/core/scan"
	"dxkite.cn/explorer/src/core/storage"
)

func CreateIndex(root, output string) error {
	s := scan.NewScanner(output)
	return s.Scan(context.TODO(), storage.Local(root))
}

func CreateIndexForStorage(ctx context.Context, fs storage.FileSystem, output string) error {
	s := scan.NewScanner(output)
	return s.Scan(ctx, fs)
}
