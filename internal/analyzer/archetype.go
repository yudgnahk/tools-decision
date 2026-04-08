package analyzer

import (
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/yudgnahk/tools-decision/pkg/types"
)

type archetypeScore struct {
	value    float64
	evidence []string
}

type archetypeFileSignals struct {
	texFiles       int
	tabularFiles   int
	automationHits int
	mediaHits      int
}

func detectArchetypes(ctx *types.ProjectContext, projectPath string) []types.ArchetypeSignal {
	scores := map[types.Archetype]*archetypeScore{}

	add := func(name types.Archetype, delta float64, evidence string) {
		if delta <= 0 {
			return
		}
		entry, ok := scores[name]
		if !ok {
			entry = &archetypeScore{}
			scores[name] = entry
		}
		entry.value += delta
		if evidence != "" {
			entry.evidence = append(entry.evidence, evidence)
		}
	}

	frameworks := make(map[string]bool)
	for _, f := range ctx.Frameworks {
		frameworks[strings.ToLower(f.Name)] = true
	}

	services := make(map[string]bool)
	for _, s := range ctx.Services {
		services[strings.ToLower(s.Name)] = true
	}

	if ctx.Type == types.ProjectTypeAPI {
		add(types.ArchetypeAPIService, 0.8, "project type=api")
	}
	if ctx.Type == types.ProjectTypeCLI {
		add(types.ArchetypeCLITool, 0.8, "project type=cli")
	}
	if ctx.Type == types.ProjectTypeDesktop {
		add(types.ArchetypeDesktopApp, 0.8, "project type=desktop")
	}
	if ctx.Type == types.ProjectTypeLibrary {
		add(types.ArchetypeLibrary, 0.8, "project type=library")
	}

	if hasAny(frameworks, "express", "nestjs", "gin", "echo", "fiber", "chi", "fastapi", "django", "flask", "actix", "axum", "rocket", "spring", "spring-boot", "grpc", "tonic") {
		add(types.ArchetypeAPIService, 0.45, "api framework detected")
	}

	if hasAny(frameworks, "cobra", "clap", "structopt", "typer", "click", "cli") {
		add(types.ArchetypeCLITool, 0.35, "cli framework detected")
	}

	if hasAny(frameworks, "tauri", "druid", "gtk", "egui", "iced", "cocoa", "dioxus", "yew", "leptos") {
		add(types.ArchetypeDesktopApp, 0.6, "desktop framework detected")
	}

	if hasAny(frameworks, "playwright", "puppeteer", "selenium") || hasAny(services, "playwright", "puppeteer", "selenium") {
		add(types.ArchetypeAutomationBot, 0.7, "browser automation stack detected")
	}

	if hasAny(services, "openai", "anthropic", "langchain") {
		add(types.ArchetypeAIContentPipe, 0.6, "genai service detected")
	}

	if hasAny(services, "pandas", "openpyxl", "excel", "csv") {
		add(types.ArchetypeDataProcessing, 0.6, "data processing dependency detected")
	}

	fileSignals := collectArchetypeFileSignals(projectPath)
	if fileSignals.texFiles >= 2 {
		add(types.ArchetypeDocumentAuthor, 0.95, "multiple .tex files")
	} else if fileSignals.texFiles == 1 {
		add(types.ArchetypeDocumentAuthor, 0.65, ".tex files present")
	}

	if fileSignals.tabularFiles >= 2 {
		add(types.ArchetypeDataProcessing, 0.55, "tabular input files detected")
	}

	if fileSignals.automationHits > 0 {
		add(types.ArchetypeAutomationBot, 0.35, "automation script naming detected")
	}

	if fileSignals.mediaHits > 0 {
		add(types.ArchetypeAIContentPipe, 0.35, "media pipeline files detected")
	}

	if ctx.Type == types.ProjectTypeUnknown && len(ctx.Frameworks) == 0 && len(ctx.Services) == 0 && len(ctx.Languages) > 0 {
		add(types.ArchetypeLibrary, 0.4, "codebase without app/runtime framework")
	}

	hasRuntimeSignal := scores[types.ArchetypeAPIService] != nil || scores[types.ArchetypeAutomationBot] != nil || scores[types.ArchetypeDesktopApp] != nil
	if hasRuntimeSignal {
		if doc, ok := scores[types.ArchetypeDocumentAuthor]; ok {
			doc.value -= 0.35
		}
	}

	var out []types.ArchetypeSignal
	for name, score := range scores {
		if score.value < 0.5 {
			continue
		}
		if score.value > 1.0 {
			score.value = 1.0
		}
		out = append(out, types.ArchetypeSignal{
			Name:       name,
			Confidence: score.value,
			Evidence:   uniqueStrings(score.evidence),
		})
	}

	sort.Slice(out, func(i, j int) bool {
		if out[i].Confidence != out[j].Confidence {
			return out[i].Confidence > out[j].Confidence
		}
		return out[i].Name < out[j].Name
	})

	return out
}

func collectArchetypeFileSignals(projectPath string) archetypeFileSignals {
	signals := archetypeFileSignals{}

	_ = filepath.WalkDir(projectPath, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return nil
		}

		if d.IsDir() {
			if shouldIgnoreDir(d.Name()) {
				return filepath.SkipDir
			}
			return nil
		}

		name := strings.ToLower(d.Name())
		ext := strings.ToLower(filepath.Ext(name))

		switch ext {
		case ".tex":
			signals.texFiles++
		case ".csv", ".tsv", ".xlsx", ".xls":
			signals.tabularFiles++
		case ".mp3", ".wav", ".m4a", ".mp4", ".mov", ".mkv", ".srt":
			signals.mediaHits++
		}

		if strings.Contains(name, "automation") || strings.Contains(name, "bot") || strings.Contains(name, "crawl") || strings.Contains(name, "scrape") {
			signals.automationHits++
		}

		if strings.Contains(name, "ffmpeg") || strings.Contains(name, "moviepy") || strings.Contains(name, "pydub") {
			signals.mediaHits++
		}

		return nil
	})

	return signals
}

func hasAny(values map[string]bool, keys ...string) bool {
	for _, k := range keys {
		if values[strings.ToLower(k)] {
			return true
		}
	}
	return false
}

func uniqueStrings(values []string) []string {
	if len(values) == 0 {
		return nil
	}
	seen := make(map[string]bool)
	out := make([]string, 0, len(values))
	for _, v := range values {
		if seen[v] {
			continue
		}
		seen[v] = true
		out = append(out, v)
	}
	return out
}
