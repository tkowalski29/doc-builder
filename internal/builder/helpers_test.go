package builder

import (
	"os"
	"path/filepath"
	"testing"
)

func TestEnsureTempGitignore(t *testing.T) {
	dir := t.TempDir()
	if err := ensureTempGitignore(dir); err != nil {
		t.Fatalf("ensureTempGitignore returned error: %v", err)
	}

	data, err := os.ReadFile(filepath.Join(dir, ".gitignore"))
	if err != nil {
		t.Fatalf("expected .gitignore file, read failed: %v", err)
	}

	expected := "*\n!.gitignore\n"
	if string(data) != expected {
		t.Fatalf("unexpected gitignore content: %q", data)
	}
}

func TestCopyFileAndDirectory(t *testing.T) {
	srcDir := t.TempDir()
	dstDir := t.TempDir()

	srcFile := filepath.Join(srcDir, "nested", "file.txt")
	if err := os.MkdirAll(filepath.Dir(srcFile), 0o755); err != nil {
		t.Fatalf("mkdir failed: %v", err)
	}
	if err := os.WriteFile(srcFile, []byte("hello"), 0o644); err != nil {
		t.Fatalf("write failed: %v", err)
	}

	dstFile := filepath.Join(dstDir, "output", "file.txt")
	if err := copyFile(srcFile, dstFile); err != nil {
		t.Fatalf("copyFile returned error: %v", err)
	}
	data, err := os.ReadFile(dstFile)
	if err != nil || string(data) != "hello" {
		t.Fatalf("expected copied content 'hello', got %q (err=%v)", data, err)
	}

	srcStructure := filepath.Join(srcDir, "dir")
	if err := os.MkdirAll(srcStructure, 0o755); err != nil {
		t.Fatalf("mkdir for directory copy failed: %v", err)
	}
	if err := os.WriteFile(filepath.Join(srcStructure, "inside.txt"), []byte("world"), 0o644); err != nil {
		t.Fatalf("write nested failed: %v", err)
	}
	dstStructure := filepath.Join(dstDir, "copied")
	if err := copyDirectory(srcStructure, dstStructure); err != nil {
		t.Fatalf("copyDirectory returned error: %v", err)
	}
	if data, err := os.ReadFile(filepath.Join(dstStructure, "inside.txt")); err != nil || string(data) != "world" {
		t.Fatalf("expected directory copy to include file, got %q (err=%v)", data, err)
	}
}

func TestDirectoryHelpers(t *testing.T) {
	skipped := []string{"node_modules", "vendor", ".git"}
	for _, name := range skipped {
		if !shouldSkipDirectory(name) {
			t.Fatalf("expected %s to be skipped", name)
		}
	}

	if shouldSkipDirectory("content") {
		t.Fatalf("did not expect generic directory to be skipped")
	}

	if !samePath("/tmp/../tmp/docs", "/tmp/docs") {
		t.Fatalf("expected samePath to normalize entries")
	}
}
