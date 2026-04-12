package projectgen

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
)

func ensureTargetDir(targetPath string, inPlace bool) error {
	if inPlace {
		return ensureInPlaceTargetDir(targetPath)
	}

	if _, err := os.Stat(targetPath); err == nil {
		return fmt.Errorf("target path already exists: %s", targetPath)
	} else if !errors.Is(err, fs.ErrNotExist) {
		return fmt.Errorf("stat target path: %w", err)
	}

	if err := os.MkdirAll(targetPath, 0o755); err != nil {
		return fmt.Errorf("create target directory: %w", err)
	}

	return nil
}

func ensureInPlaceTargetDir(targetPath string) error {
	entries, err := os.ReadDir(targetPath)
	if err != nil {
		return fmt.Errorf("read target directory: %w", err)
	}

	for _, entry := range entries {
		if isIgnorableInPlaceEntry(entry) {
			continue
		}
		return fmt.Errorf("target directory is not empty: %s", targetPath)
	}

	return nil
}

func isIgnorableInPlaceEntry(entry fs.DirEntry) bool {
	switch entry.Name() {
	case ".DS_Store", ".localized":
		return true
	default:
		return false
	}
}

func writeFile(targetPath string, relativePath string, content []byte) error {
	fullPath := filepath.Join(targetPath, relativePath)
	if err := os.MkdirAll(filepath.Dir(fullPath), 0o755); err != nil {
		return fmt.Errorf("create parent directory for %q: %w", relativePath, err)
	}
	if err := os.WriteFile(fullPath, content, 0o644); err != nil {
		return fmt.Errorf("write %q: %w", relativePath, err)
	}

	return nil
}
