package builder

import (
	"bufio"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

type menuRecord struct {
	CategoryPath string
	Slug         string
	Title        string
}

func (b *Builder) collectPrefixedDocs(env environment, recordSet map[string]struct{}) ([]menuRecord, int, error) {
	if b.cfg.Verbose {
		fmt.Printf("[2/7] Scanning %s for files starting with %s\n", env.searchRoot, b.cfg.Prefix)
	}

	var menuRecords []menuRecord
	count := 0

	err := filepath.WalkDir(env.searchRoot, func(path string, d fs.DirEntry, walkErr error) error {
		if walkErr != nil {
			return walkErr
		}

		if d.IsDir() {
			if shouldSkipDirectory(d.Name()) {
				return filepath.SkipDir
			}
			if samePath(path, filepath.Join(env.docDir, b.cfg.TempDirName)) {
				return filepath.SkipDir
			}
			return nil
		}

		if filepath.Ext(path) != ".md" {
			return nil
		}

		base := filepath.Base(path)
		if !strings.HasPrefix(base, b.cfg.Prefix) {
			return nil
		}

		//nolint:gosec // file path is validated and safe
		data, err := os.ReadFile(path)
		if err != nil {
			return fmt.Errorf("failed to read %s: %w", path, err)
		}

		count++
		fm := parseFrontMatter(data)
		category := strings.TrimSpace(fm["category"])
		if category == "" {
			category = "guides"
		}
		categoryPath := normalizeCategoryPath(category)
		slug := buildSlug(base, b.cfg.Prefix)
		title := deriveTitle(data, fm["title"], slug, b.cfg.Prefix)

		targetDir := filepath.Join(env.tempDir, filepath.FromSlash(categoryPath))
		if err := os.MkdirAll(targetDir, 0o755); err != nil {
			return fmt.Errorf("failed to create directory %s: %w", targetDir, err)
		}
		targetFile := filepath.Join(targetDir, slug+".md")
		if err := os.WriteFile(targetFile, data, 0o644); err != nil {
			return fmt.Errorf("failed to copy %s to %s: %w", path, targetFile, err)
		}

		key := menuKey(categoryPath, slug)
		if _, exists := recordSet[key]; !exists {
			recordSet[key] = struct{}{}
			menuRecords = append(menuRecords, menuRecord{CategoryPath: categoryPath, Slug: slug, Title: title})
		}

		if b.cfg.Verbose {
			fmt.Printf("  collected %s -> %s\n", path, targetFile)
		}
		return nil
	})
	if err != nil {
		return nil, 0, err
	}

	return menuRecords, count, nil
}

func (b *Builder) collectExistingDocs(env environment, recordSet map[string]struct{}) ([]menuRecord, int, error) {
	if b.cfg.Verbose {
		fmt.Printf("[3/7] Merging existing documentation from %s\n", env.docDir)
	}

	var menuRecords []menuRecord
	count := 0

	err := filepath.WalkDir(env.docDir, func(path string, d fs.DirEntry, walkErr error) error {
		if walkErr != nil {
			return walkErr
		}

		rel, err := filepath.Rel(env.docDir, path)
		if err != nil {
			return err
		}

		if rel == "." {
			return nil
		}

		if d.IsDir() {
			base := d.Name()
			if base == b.cfg.TempDirName || base == ".vitepress" || base == "node_modules" {
				return filepath.SkipDir
			}
			return nil
		}

		if filepath.Ext(rel) != ".md" {
			return nil
		}

		if rel == "index.md" || rel == "DOC_BUILD_README.md" {
			return nil
		}

		//nolint:gosec // file path is validated and safe
		data, err := os.ReadFile(path)
		if err != nil {
			return fmt.Errorf("failed to read %s: %w", path, err)
		}

		count++
		category := normalizeCategoryPath(filepath.ToSlash(filepath.Dir(rel)))
		slug := strings.TrimSuffix(filepath.Base(rel), ".md")
		title := deriveTitle(data, "", slug, "")
		if slug == "index" && title != "" {
			title = fmt.Sprintf("%s (overview)", title)
		}

		targetPath := filepath.Join(env.tempDir, filepath.FromSlash(rel))
		if err := os.MkdirAll(filepath.Dir(targetPath), 0o755); err != nil {
			return fmt.Errorf("failed to create directory %s: %w", filepath.Dir(targetPath), err)
		}
		if err := os.WriteFile(targetPath, data, 0o644); err != nil {
			return fmt.Errorf("failed to copy %s to %s: %w", path, targetPath, err)
		}

		key := menuKey(category, slug)
		if _, exists := recordSet[key]; !exists {
			recordSet[key] = struct{}{}
			menuRecords = append(menuRecords, menuRecord{CategoryPath: category, Slug: slug, Title: title})
		}

		if b.cfg.Verbose {
			fmt.Printf("  merged %s\n", rel)
		}
		return nil
	})
	if err != nil {
		return nil, 0, err
	}

	indexPath := filepath.Join(env.docDir, "index.md")
	if _, err := os.Stat(indexPath); err == nil {
		data, readErr := os.ReadFile(indexPath)
		if readErr != nil {
			return nil, 0, fmt.Errorf("failed to read %s: %w", indexPath, readErr)
		}
		target := filepath.Join(env.tempDir, "index.md")
		if writeErr := os.WriteFile(target, data, 0o644); writeErr != nil {
			return nil, 0, fmt.Errorf("failed to copy %s to %s: %w", indexPath, target, writeErr)
		}
	}

	return menuRecords, count, nil
}

func (b *Builder) writeMenuIndex(env environment, records []menuRecord) error {
	menuPath := filepath.Join(env.tempDir, ".menu-items.txt")
	file, err := os.Create(menuPath)
	if err != nil {
		return fmt.Errorf("failed to create %s: %w", menuPath, err)
	}
	defer file.Close()

	writer := bufio.NewWriter(file)
	for _, rec := range records {
		line := fmt.Sprintf("%s|%s|%s\n", rec.CategoryPath, rec.Slug, rec.Title)
		if _, err := writer.WriteString(line); err != nil {
			return fmt.Errorf("failed to write menu entry to %s: %w", menuPath, err)
		}
	}
	if err := writer.Flush(); err != nil {
		return fmt.Errorf("failed to flush %s: %w", menuPath, err)
	}
	return nil
}
