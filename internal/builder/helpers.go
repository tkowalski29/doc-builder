package builder

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

func ensureTempGitignore(tempDir string) error {
	content := []byte("*\n!.gitignore\n")
	target := filepath.Join(tempDir, ".gitignore")
	if err := os.WriteFile(target, content, 0o644); err != nil {
		return fmt.Errorf("failed to initialise %s: %w", target, err)
	}
	return nil
}

func shouldSkipDirectory(name string) bool {
	switch name {
	case "node_modules", "vendor", ".git", ".hg", ".svn":
		return true
	}
	return false
}

func samePath(a, b string) bool {
	return filepath.Clean(a) == filepath.Clean(b)
}

func copyFile(src, dst string) error {
	input, err := os.ReadFile(src)
	if err != nil {
		return fmt.Errorf("failed to read %s: %w", src, err)
	}
	if err := os.MkdirAll(filepath.Dir(dst), 0o755); err != nil {
		return fmt.Errorf("failed to prepare directory for %s: %w", dst, err)
	}
	if err := os.WriteFile(dst, input, 0o644); err != nil {
		return fmt.Errorf("failed to write %s: %w", dst, err)
	}
	return nil
}

func copyDirectory(src, dst string) error {
	return filepath.WalkDir(src, func(path string, d fs.DirEntry, walkErr error) error {
		if walkErr != nil {
			return walkErr
		}
		rel, err := filepath.Rel(src, path)
		if err != nil {
			return err
		}
		target := filepath.Join(dst, rel)
		if d.IsDir() {
			if err := os.MkdirAll(target, 0o755); err != nil {
				return fmt.Errorf("failed to create directory %s: %w", target, err)
			}
			return nil
		}
		return copyFile(path, target)
	})
}

func escapeQuotes(value string) string {
	return strings.ReplaceAll(value, "'", "\\'")
}
