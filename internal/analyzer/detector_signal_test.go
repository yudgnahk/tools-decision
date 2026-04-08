package analyzer

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/yudgnahk/tools-decision/pkg/types"
)

func TestAnalyze_JavaScriptAutomationSignals(t *testing.T) {
	root := t.TempDir()
	mustWriteFile(t, filepath.Join(root, "package.json"), `{
		"name": "bot-app",
		"dependencies": {
			"playwright": "^1.49.0",
			"openai": "^4.0.0"
		}
	}`)

	a := New()
	ctx, err := a.Analyze(context.Background(), root)
	if err != nil {
		t.Fatalf("analyze failed: %v", err)
	}

	if !hasFramework(ctx.Frameworks, "playwright") {
		t.Fatalf("expected playwright framework signal, got %+v", ctx.Frameworks)
	}
	if !hasService(ctx.Services, "openai") {
		t.Fatalf("expected openai service signal, got %+v", ctx.Services)
	}

	if _, ok := findArchetype(ctx.Archetypes, types.ArchetypeAutomationBot); !ok {
		t.Fatalf("expected automation archetype, got %+v", ctx.Archetypes)
	}
}

func TestAnalyze_GoExcelizeSignalsDataArchetype(t *testing.T) {
	root := t.TempDir()
	mustWriteFile(t, filepath.Join(root, "go.mod"), `module example.com/etl

go 1.22

require github.com/xuri/excelize/v2 v2.9.0
`)

	a := New()
	ctx, err := a.Analyze(context.Background(), root)
	if err != nil {
		t.Fatalf("analyze failed: %v", err)
	}

	if !hasService(ctx.Services, "excelize") {
		t.Fatalf("expected excelize service signal, got %+v", ctx.Services)
	}
	if _, ok := findArchetype(ctx.Archetypes, types.ArchetypeDataProcessing); !ok {
		t.Fatalf("expected data-processing archetype, got %+v", ctx.Archetypes)
	}
}

func TestAnalyze_GoAPIPreferredOverCmdCLIDefault(t *testing.T) {
	root := t.TempDir()
	mustWriteFile(t, filepath.Join(root, "go.mod"), `module example.com/service

go 1.22

require github.com/gin-gonic/gin v1.10.0
`)
	mustWriteFile(t, filepath.Join(root, "cmd", "api", "main.go"), "package main")

	a := New()
	ctx, err := a.Analyze(context.Background(), root)
	if err != nil {
		t.Fatalf("analyze failed: %v", err)
	}

	if ctx.Type != types.ProjectTypeAPI {
		t.Fatalf("expected project type api, got %s", ctx.Type)
	}
	if _, ok := findArchetype(ctx.Archetypes, types.ArchetypeAPIService); !ok {
		t.Fatalf("expected API archetype, got %+v", ctx.Archetypes)
	}
}

func TestAnalyze_PyprojectDependenciesFeedServiceSignals(t *testing.T) {
	root := t.TempDir()
	mustWriteFile(t, filepath.Join(root, "pyproject.toml"), `[project]
name = "media-pipeline"
dependencies = [
  "openai>=1.0.0",
  "moviepy>=1.0.3",
]
`)

	a := New()
	ctx, err := a.Analyze(context.Background(), root)
	if err != nil {
		t.Fatalf("analyze failed: %v", err)
	}

	if !hasService(ctx.Services, "openai") {
		t.Fatalf("expected openai service signal from pyproject, got %+v", ctx.Services)
	}
	if !hasService(ctx.Services, "moviepy") {
		t.Fatalf("expected moviepy service signal from pyproject, got %+v", ctx.Services)
	}
	if _, ok := findArchetype(ctx.Archetypes, types.ArchetypeAIContentPipe); !ok {
		t.Fatalf("expected AI pipeline archetype, got %+v", ctx.Archetypes)
	}
}

func mustWriteFile(t *testing.T, path string, content string) {
	t.Helper()
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		t.Fatalf("mkdir failed: %v", err)
	}
	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		t.Fatalf("write failed: %v", err)
	}
}

func hasService(services []types.Service, want string) bool {
	for _, s := range services {
		if s.Name == want {
			return true
		}
	}
	return false
}

func hasFramework(frameworks []types.Framework, want string) bool {
	for _, f := range frameworks {
		if f.Name == want {
			return true
		}
	}
	return false
}
