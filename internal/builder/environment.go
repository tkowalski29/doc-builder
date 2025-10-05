package builder

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

type environment struct {
	docDir       string
	searchRoot   string
	tempDir      string
	baseConfig   string
	outputConfig string
	distSrc      string
	distDst      string
}

func (b *Builder) validateConfig() error {
	if b.cfg.Prefix == "" {
		return errors.New("prefix cannot be empty")
	}
	if b.cfg.Engine == "" {
		return errors.New("engine cannot be empty")
	}
	if !strings.EqualFold(b.cfg.Engine, "vitepress") {
		return fmt.Errorf("unsupported engine '%s': only 'vitepress' is currently available", b.cfg.Engine)
	}
	if b.cfg.SearchPath == "" {
		return errors.New("search path cannot be empty")
	}
	if b.cfg.DocDir == "" {
		return errors.New("documentation directory cannot be empty")
	}
	if b.cfg.TempDirName == "" {
		return errors.New("temporary directory name cannot be empty")
	}
	return nil
}

func (b *Builder) prepareEnvironment() (environment, error) {
	docDir, err := filepath.Abs(b.cfg.DocDir)
	if err != nil {
		return environment{}, fmt.Errorf("failed to resolve documentation directory: %w", err)
	}
	searchRoot, err := filepath.Abs(b.cfg.SearchPath)
	if err != nil {
		return environment{}, fmt.Errorf("failed to resolve search path: %w", err)
	}

	if _, err := os.Stat(docDir); err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return environment{}, fmt.Errorf("documentation directory not found: %s", docDir)
		}
		return environment{}, fmt.Errorf("failed to access documentation directory: %w", err)
	}

	if _, err := os.Stat(searchRoot); err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return environment{}, fmt.Errorf("search path not found: %s", searchRoot)
		}
		return environment{}, fmt.Errorf("failed to access search path: %w", err)
	}

	tempDir := filepath.Join(docDir, b.cfg.TempDirName)
	baseConfig := filepath.Join(docDir, ".vitepress", "base.config.js")
	outputConfig := filepath.Join(docDir, ".vitepress", "config.js")
	distSrc := filepath.Join(tempDir, ".vitepress", "dist")
	distDst := filepath.Join(docDir, ".vitepress", "dist")

	if _, err := os.Stat(baseConfig); err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return environment{}, fmt.Errorf("expected file not found: %s", baseConfig)
		}
		return environment{}, fmt.Errorf("failed to access base config: %w", err)
	}

	packageJSON := filepath.Join(docDir, "package.json")
	if _, err := os.Stat(packageJSON); err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return environment{}, fmt.Errorf("expected file not found: %s", packageJSON)
		}
		return environment{}, fmt.Errorf("failed to access %s: %w", packageJSON, err)
	}

	return environment{
		docDir:       docDir,
		searchRoot:   searchRoot,
		tempDir:      tempDir,
		baseConfig:   baseConfig,
		outputConfig: outputConfig,
		distSrc:      distSrc,
		distDst:      distDst,
	}, nil
}

func (b *Builder) prepareWorkspace(env environment) error {
	if b.cfg.Verbose {
		fmt.Printf("[1/7] Cleaning temporary directory %s\n", env.tempDir)
	}
	if err := os.RemoveAll(env.tempDir); err != nil {
		return fmt.Errorf("failed to clean temporary directory %s: %w", env.tempDir, err)
	}
	if err := os.MkdirAll(env.tempDir, 0o755); err != nil {
		return fmt.Errorf("failed to create temporary directory %s: %w", env.tempDir, err)
	}
	if err := ensureTempGitignore(env.tempDir); err != nil {
		return err
	}
	return nil
}
