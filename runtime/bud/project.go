package bud

import (
	"context"
	"io"
	"io/fs"
	"os"
	"path/filepath"

	"gitlab.com/mnm/bud/internal/gobin"
	"gitlab.com/mnm/bud/package/gomod"
	"gitlab.com/mnm/bud/package/overlay"
)

func New(fsys *overlay.FileSystem, module *gomod.Module) *Project {
	return &Project{
		fsys:   fsys,
		module: module,
		Env:    os.Environ(),
		Stderr: os.Stderr,
		Stdout: os.Stdout,
	}
}

type Project struct {
	fsys   *overlay.FileSystem
	module *gomod.Module
	Env    []string
	Stdout io.Writer
	Stderr io.Writer
}

func (c *Project) Compile(ctx context.Context, flag *Flag) (*App, error) {
	// Sync the app
	if err := c.fsys.Sync("bud/.app"); err != nil {
		return nil, err
	}
	// Ensure that main.go exists
	if _, err := fs.Stat(c.module, "bud/.app/main.go"); err != nil {
		return nil, err
	}
	// Build the binary
	if err := gobin.Build(ctx, c.module.Directory(), "bud/.app/main.go", filepath.Join("bud", "app")); err != nil {
		return nil, err
	}
	return &App{
		Module: c.module,
		Env:    c.Env,
		Stderr: c.Stderr,
		Stdout: c.Stdout,
	}, nil
}