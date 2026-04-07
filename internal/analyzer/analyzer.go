package analyzer

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/yudgnahk/tools-decision/pkg/types"
)

// Analyzer analyzes a project directory to extract context
type Analyzer struct {
	detectors []Detector
}

// Detector is the interface for language/framework detectors
type Detector interface {
	// Name returns the detector name
	Name() string
	// Detect analyzes the project and returns partial context
	Detect(ctx context.Context, projectPath string) (*DetectorResult, error)
}

// DetectorResult contains results from a single detector
type DetectorResult struct {
	Languages  []types.Language
	Frameworks []types.Framework
	Tools      []types.Tool
	Services   []types.Service
	Type       types.ProjectType
}

// New creates a new Analyzer with all detectors
func New() *Analyzer {
	return &Analyzer{
		detectors: []Detector{
			NewJavaScriptDetector(),
			NewPythonDetector(),
			NewGoDetector(),
			NewRustDetector(),
			NewJavaDetector(),
		},
	}
}

// Analyze analyzes the project at the given path
func (a *Analyzer) Analyze(ctx context.Context, projectPath string) (*types.ProjectContext, error) {
	// Ensure path exists
	if _, err := os.Stat(projectPath); os.IsNotExist(err) {
		return nil, fmt.Errorf("project path does not exist: %s", projectPath)
	}

	// Get absolute path
	absPath, err := filepath.Abs(projectPath)
	if err != nil {
		return nil, fmt.Errorf("failed to get absolute path: %w", err)
	}

	result := &types.ProjectContext{
		Path: absPath,
		Type: types.ProjectTypeUnknown,
	}

	// Run all detectors
	for _, detector := range a.detectors {
		detResult, err := detector.Detect(ctx, absPath)
		if err != nil {
			// Log warning but continue with other detectors
			continue
		}

		if detResult == nil {
			continue
		}

		// Merge results
		result.Languages = append(result.Languages, detResult.Languages...)
		result.Frameworks = append(result.Frameworks, detResult.Frameworks...)
		result.Tools = append(result.Tools, detResult.Tools...)
		result.Services = append(result.Services, detResult.Services...)

		// Take highest confidence project type
		if detResult.Type != types.ProjectTypeUnknown {
			result.Type = detResult.Type
		}
	}

	// Sort languages/frameworks by confidence
	sortByConfidence(result)

	return result, nil
}

// sortByConfidence sorts languages and frameworks by confidence descending
func sortByConfidence(ctx *types.ProjectContext) {
	// Sort languages
	for i := 0; i < len(ctx.Languages); i++ {
		for j := i + 1; j < len(ctx.Languages); j++ {
			if ctx.Languages[j].Confidence > ctx.Languages[i].Confidence {
				ctx.Languages[i], ctx.Languages[j] = ctx.Languages[j], ctx.Languages[i]
			}
		}
	}

	// Sort frameworks
	for i := 0; i < len(ctx.Frameworks); i++ {
		for j := i + 1; j < len(ctx.Frameworks); j++ {
			if ctx.Frameworks[j].Confidence > ctx.Frameworks[i].Confidence {
				ctx.Frameworks[i], ctx.Frameworks[j] = ctx.Frameworks[j], ctx.Frameworks[i]
			}
		}
	}
}
