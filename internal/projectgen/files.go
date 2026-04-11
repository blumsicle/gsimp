package projectgen

import (
	"fmt"
	"os"
	"path/filepath"
)

func ensureTargetDir(targetPath string) error {
	if _, err := os.Stat(targetPath); err == nil {
		return fmt.Errorf("target path already exists: %s", targetPath)
	} else if !os.IsNotExist(err) {
		return fmt.Errorf("stat target path: %w", err)
	}

	if err := os.MkdirAll(targetPath, 0o755); err != nil {
		return fmt.Errorf("create target directory: %w", err)
	}

	return nil
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
