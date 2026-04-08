package analyzer

import (
	"os"
	"path/filepath"
	"testing"
)

func TestShouldIgnoreDir(t *testing.T) {
	if !shouldIgnoreDir("node_modules") {
		t.Fatalf("expected node_modules to be ignored")
	}
	if !shouldIgnoreDir("chromium-data") {
		t.Fatalf("expected chromium-data to be ignored")
	}
	if shouldIgnoreDir("src") {
		t.Fatalf("expected src to not be ignored")
	}
}

func TestCountSourceFiles_IgnoresNoisyDirs(t *testing.T) {
	root := t.TempDir()

	mustWrite := func(p string) {
		dir := filepath.Dir(p)
		if err := os.MkdirAll(dir, 0755); err != nil {
			t.Fatalf("mkdir failed: %v", err)
		}
		if err := os.WriteFile(p, []byte("x"), 0644); err != nil {
			t.Fatalf("write failed: %v", err)
		}
	}

	mustWrite(filepath.Join(root, "main.go"))
	mustWrite(filepath.Join(root, "app.js"))
	mustWrite(filepath.Join(root, "node_modules", "lib", "x.js"))
	mustWrite(filepath.Join(root, "chromium-data", "cache", "y.js"))

	counts := countSourceFiles(root)
	if counts["go"] != 1 {
		t.Fatalf("expected go count=1, got %d", counts["go"])
	}
	if counts["javascript"] != 1 {
		t.Fatalf("expected javascript count=1, got %d", counts["javascript"])
	}
}
