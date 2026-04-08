package analyzer

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/yudgnahk/tools-decision/pkg/types"
)

func TestDetectArchetypes_DocumentAuthoring(t *testing.T) {
	root := t.TempDir()
	mustWrite(t, filepath.Join(root, "main.tex"))
	mustWrite(t, filepath.Join(root, "sections", "intro.tex"))

	ctx := &types.ProjectContext{}
	arch := detectArchetypes(ctx, root)

	signal, ok := findArchetype(arch, types.ArchetypeDocumentAuthor)
	if !ok {
		t.Fatalf("expected %s archetype, got %+v", types.ArchetypeDocumentAuthor, arch)
	}
	if signal.Confidence < 0.8 {
		t.Fatalf("expected high confidence for document archetype, got %.2f", signal.Confidence)
	}
}

func TestDetectArchetypes_DesktopApp(t *testing.T) {
	root := t.TempDir()
	mustWrite(t, filepath.Join(root, "src", "main.rs"))

	ctx := &types.ProjectContext{
		Type: types.ProjectTypeDesktop,
		Frameworks: []types.Framework{
			{Name: "tauri", Confidence: 0.95},
		},
	}

	arch := detectArchetypes(ctx, root)
	if _, ok := findArchetype(arch, types.ArchetypeDesktopApp); !ok {
		t.Fatalf("expected %s archetype, got %+v", types.ArchetypeDesktopApp, arch)
	}
}

func TestDetectArchetypes_AutomationAndAIPipeline(t *testing.T) {
	root := t.TempDir()
	mustWrite(t, filepath.Join(root, "automation_bot.py"))
	mustWrite(t, filepath.Join(root, "assets", "episode.mp4"))

	ctx := &types.ProjectContext{
		Services: []types.Service{
			{Name: "playwright", Confidence: 0.9},
			{Name: "openai", Confidence: 0.9},
		},
	}

	arch := detectArchetypes(ctx, root)
	if _, ok := findArchetype(arch, types.ArchetypeAutomationBot); !ok {
		t.Fatalf("expected %s archetype, got %+v", types.ArchetypeAutomationBot, arch)
	}
	if _, ok := findArchetype(arch, types.ArchetypeAIContentPipe); !ok {
		t.Fatalf("expected %s archetype, got %+v", types.ArchetypeAIContentPipe, arch)
	}
}

func TestDetectArchetypes_DataProcessing(t *testing.T) {
	root := t.TempDir()
	mustWrite(t, filepath.Join(root, "input", "students.csv"))
	mustWrite(t, filepath.Join(root, "input", "scores.xlsx"))

	ctx := &types.ProjectContext{
		Services: []types.Service{
			{Name: "pandas", Confidence: 0.9},
		},
	}

	arch := detectArchetypes(ctx, root)
	if _, ok := findArchetype(arch, types.ArchetypeDataProcessing); !ok {
		t.Fatalf("expected %s archetype, got %+v", types.ArchetypeDataProcessing, arch)
	}
}

func mustWrite(t *testing.T, p string) {
	t.Helper()
	dir := filepath.Dir(p)
	if err := os.MkdirAll(dir, 0755); err != nil {
		t.Fatalf("mkdir failed: %v", err)
	}
	if err := os.WriteFile(p, []byte("x"), 0644); err != nil {
		t.Fatalf("write failed: %v", err)
	}
}

func findArchetype(values []types.ArchetypeSignal, want types.Archetype) (types.ArchetypeSignal, bool) {
	for _, a := range values {
		if a.Name == want {
			return a, true
		}
	}
	return types.ArchetypeSignal{}, false
}
