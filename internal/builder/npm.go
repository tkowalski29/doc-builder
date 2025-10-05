package builder

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

func (b *Builder) installDependencies(ctx context.Context, env environment) error {
	if b.cfg.Verbose {
		fmt.Println("[5/7] Preparing Node.js dependencies")
	}

	filesToCopy := []string{"package.json", "package-lock.json", "yarn.lock", "pnpm-lock.yaml"}
	for _, name := range filesToCopy {
		src := filepath.Join(env.docDir, name)
		if _, err := os.Stat(src); err == nil {
			dst := filepath.Join(env.tempDir, name)
			if err := copyFile(src, dst); err != nil {
				return err
			}
		}
	}

	nodeModules := filepath.Join(env.tempDir, "node_modules")
	if _, err := os.Stat(nodeModules); err == nil {
		if b.cfg.Verbose {
			fmt.Println("  node_modules already present, skipping npm install")
		}
		return nil
	}

	if b.cfg.Verbose {
		fmt.Println("  running npm install (this may take a while)")
	}

	//nolint:gosec // npm install is safe in controlled environment
	cmd := exec.CommandContext(ctx, "npm", "install")
	cmd.Dir = env.tempDir
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("npm install failed: %w", err)
	}
	return nil
}

func (b *Builder) buildSite(ctx context.Context, env environment) error {
	if b.cfg.Verbose {
		fmt.Println("[6/7] Running npm run docs:build")
	}

	//nolint:gosec // npm run is safe in controlled environment
	cmd := exec.CommandContext(ctx, "npm", "run", "docs:build")
	cmd.Dir = env.tempDir
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("npm run docs:build failed: %w", err)
	}
	return nil
}

func (b *Builder) publishDist(env environment) error {
	if b.cfg.Verbose {
		fmt.Println("[7/7] Publishing VitePress dist output")
	}

	if _, err := os.Stat(env.distSrc); err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return fmt.Errorf("expected directory not found: %s", env.distSrc)
		}
		return fmt.Errorf("failed to access dist directory: %w", err)
	}

	if err := os.RemoveAll(env.distDst); err != nil {
		return fmt.Errorf("failed to clean %s: %w", env.distDst, err)
	}

	if err := copyDirectory(env.distSrc, env.distDst); err != nil {
		return err
	}
	return nil
}

func (b *Builder) printSummary(env environment, prefCount, existingCount, menuCount int) {
	fmt.Println("Build complete.")
	fmt.Printf("  Found %d prefixed markdown files\n", prefCount)
	fmt.Printf("  Merged %d existing documentation files\n", existingCount)
	fmt.Printf("  Sidebar entries: %d\n", menuCount)
	fmt.Printf("  Output directory: %s\n", env.distDst)
}
