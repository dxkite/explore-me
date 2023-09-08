package core

import (
	"context"

	"dxkite.cn/explorer/src/core/config"
	"dxkite.cn/explorer/src/core/scan"
	"dxkite.cn/explorer/src/core/storage"
)

func CreateIndex(cfg *config.Config) error {
	s := scan.NewScanner(cfg.DataRoot)
	ctx := context.TODO()
	ctx = context.WithValue(ctx, scan.DirConfigKey, cfg.DirConfig)
	return s.Scan(ctx, storage.Local(cfg.SrcRoot))
}

func CreateIndexForStorage(ctx context.Context, fs storage.FileSystem, output string) error {
	s := scan.NewScanner(output)
	return s.Scan(ctx, fs)
}
