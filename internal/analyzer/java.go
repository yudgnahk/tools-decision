package analyzer

import (
	"bufio"
	"context"
	"os"
	"path/filepath"
	"strings"

	"github.com/yudgnahk/tools-decision/pkg/types"
)

// JavaDetector detects Java/Kotlin projects
type JavaDetector struct{}

// NewJavaDetector creates a new Java detector
func NewJavaDetector() *JavaDetector {
	return &JavaDetector{}
}

// Name returns the detector name
func (d *JavaDetector) Name() string {
	return "java"
}

// Detect analyzes the project for Java/Kotlin
func (d *JavaDetector) Detect(ctx context.Context, projectPath string) (*DetectorResult, error) {
	result := &DetectorResult{}

	// Check for build files
	pomPath := filepath.Join(projectPath, "pom.xml")
	gradlePath := filepath.Join(projectPath, "build.gradle")
	gradleKtsPath := filepath.Join(projectPath, "build.gradle.kts")

	var buildFileContent string
	var hasMaven, hasGradle bool

	if data, err := os.ReadFile(pomPath); err == nil {
		hasMaven = true
		buildFileContent = string(data)
	}
	if data, err := os.ReadFile(gradlePath); err == nil {
		hasGradle = true
		buildFileContent += string(data)
	}
	if data, err := os.ReadFile(gradleKtsPath); err == nil {
		hasGradle = true
		buildFileContent += string(data)
	}

	if !hasMaven && !hasGradle {
		return nil, nil
	}

	// Check for Kotlin
	hasKotlin := strings.Contains(buildFileContent, "kotlin") ||
		strings.Contains(buildFileContent, "org.jetbrains.kotlin")

	if hasKotlin {
		result.Languages = append(result.Languages, types.Language{
			Name:       "kotlin",
			Confidence: 0.95,
		})
	}

	result.Languages = append(result.Languages, types.Language{
		Name:       "java",
		Confidence: 0.98,
	})

	// Add build tool
	if hasMaven {
		result.Tools = append(result.Tools, types.Tool{
			Name:       "maven",
			ConfigFile: "pom.xml",
		})
	}
	if hasGradle {
		result.Tools = append(result.Tools, types.Tool{
			Name:       "gradle",
			ConfigFile: "build.gradle",
		})
	}

	// Detect frameworks
	frameworkPatterns := map[string]string{
		"spring-boot":           "spring-boot",
		"org.springframework":   "spring",
		"io.quarkus":            "quarkus",
		"io.micronaut":          "micronaut",
		"io.vertx":              "vertx",
		"io.dropwizard":         "dropwizard",
		"org.apache.struts":     "struts",
		"javax.servlet":         "servlet",
		"jakarta.servlet":       "servlet",
		"io.grpc":               "grpc",
		"org.apache.kafka":      "kafka",
		"org.hibernate":         "hibernate",
		"io.ktor":               "ktor", // Kotlin
		"org.jetbrains.exposed": "exposed",
	}

	for pattern, framework := range frameworkPatterns {
		if strings.Contains(buildFileContent, pattern) {
			result.Frameworks = append(result.Frameworks, types.Framework{
				Name:       framework,
				Confidence: 0.9,
			})
		}
	}

	// Detect services
	servicePatterns := map[string]string{
		"postgresql":             "postgresql",
		"mysql":                  "mysql",
		"mariadb":                "mysql",
		"mongodb":                "mongodb",
		"redis":                  "redis",
		"elasticsearch":          "elasticsearch",
		"kafka":                  "kafka",
		"rabbitmq":               "rabbitmq",
		"aws-java-sdk":           "aws",
		"software.amazon.awssdk": "aws",
		"com.google.cloud":       "gcp",
		"com.azure":              "azure",
		"org.apache.cassandra":   "cassandra",
		"com.datastax":           "cassandra",
	}

	for pattern, service := range servicePatterns {
		if strings.Contains(strings.ToLower(buildFileContent), pattern) {
			result.Services = append(result.Services, types.Service{
				Name:       service,
				Confidence: 0.85,
			})
		}
	}

	// Detect project type
	if strings.Contains(buildFileContent, "spring-boot") ||
		strings.Contains(buildFileContent, "quarkus") ||
		strings.Contains(buildFileContent, "micronaut") {
		result.Type = types.ProjectTypeAPI
	} else if strings.Contains(buildFileContent, "spring-webflux") ||
		strings.Contains(buildFileContent, "spring-web") {
		result.Type = types.ProjectTypeWebApp
	}

	// Check for Android
	if strings.Contains(buildFileContent, "com.android") ||
		strings.Contains(buildFileContent, "android.application") {
		result.Type = types.ProjectTypeMobile
	}

	return result, nil
}

// parseGradleDeps parses dependencies from gradle build files
func parseGradleDeps(content string) map[string]bool {
	deps := make(map[string]bool)
	scanner := bufio.NewScanner(strings.NewReader(content))

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())

		// Look for dependency declarations
		// implementation 'group:artifact:version'
		// implementation("group:artifact:version")
		if strings.Contains(line, "implementation") ||
			strings.Contains(line, "compile") ||
			strings.Contains(line, "api") ||
			strings.Contains(line, "runtimeOnly") {

			// Extract dependency string
			for _, quote := range []string{"'", "\""} {
				start := strings.Index(line, quote)
				if start >= 0 {
					end := strings.Index(line[start+1:], quote)
					if end > 0 {
						dep := line[start+1 : start+1+end]
						deps[dep] = true
					}
				}
			}
		}
	}

	return deps
}
