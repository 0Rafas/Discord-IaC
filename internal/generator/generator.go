package generator

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"text/template"

	"discord-iac/internal/config"
)

func GenerateBot(cfg *config.BotConfig) error {
	outDir := "out"

	// Create directory structure
	dirs := []string{
		filepath.Join(outDir, "cmd", "bot"),
		filepath.Join(outDir, "internal", "bot"),
		filepath.Join(outDir, "internal", "handlers"),
	}

	if cfg.Database == "postgres" {
		dirs = append(dirs, filepath.Join(outDir, "internal", "database"))
	}
	if cfg.Cache == "redis" {
		dirs = append(dirs, filepath.Join(outDir, "internal", "cache"))
	}
	if cfg.Observability.Prometheus {
		dirs = append(dirs, filepath.Join(outDir, "internal", "metrics"))
	}
	if cfg.Deployment.GithubActions {
		dirs = append(dirs, filepath.Join(outDir, ".github", "workflows"))
	}
	if cfg.Deployment.Kubernetes {
		dirs = append(dirs, filepath.Join(outDir, "k8s"))
	}

	for _, d := range dirs {
		if err := os.MkdirAll(d, 0755); err != nil {
			return fmt.Errorf("failed to create directory %s: %w", d, err)
		}
	}

	// Definition of all files to generate
	files := []struct {
		Path     string
		Template string
	}{
		{filepath.Join(outDir, "cmd", "bot", "main.go"), MainTemplate},
		{filepath.Join(outDir, "internal", "bot", "client.go"), ClientTemplate},
		{filepath.Join(outDir, "internal", "handlers", "handlers.go"), HandlersTemplate},
		{filepath.Join(outDir, "go.mod"), GoModTemplate},
		{filepath.Join(outDir, "Dockerfile"), DockerfileTemplate},
		{filepath.Join(outDir, "docker-compose.yml"), DockerComposeTemplate},
	}

	if cfg.Database == "postgres" {
		files = append(files, struct{ Path, Template string }{
			filepath.Join(outDir, "internal", "database", "db.go"), DatabaseTemplate,
		})
	}
	if cfg.Cache == "redis" {
		files = append(files, struct{ Path, Template string }{
			filepath.Join(outDir, "internal", "cache", "redis.go"), CacheTemplate,
		})
	}
	if cfg.Observability.Prometheus {
		files = append(files, struct{ Path, Template string }{
			filepath.Join(outDir, "internal", "metrics", "metrics.go"), MetricsTemplate,
		})
	}
	if cfg.Deployment.GithubActions {
		files = append(files, struct{ Path, Template string }{
			filepath.Join(outDir, ".github", "workflows", "ci.yml"), GithubActionsTemplate,
		})
	}
	if cfg.Deployment.Kubernetes {
		files = append(files, struct{ Path, Template string }{
			filepath.Join(outDir, "k8s", "deployment.yaml"), K8sDeploymentTemplate,
		})
	}

	for _, f := range files {
		tmpl, err := template.New(filepath.Base(f.Path)).Parse(f.Template)
		if err != nil {
			return fmt.Errorf("failed to parse template for %s: %w", f.Path, err)
		}

		var buf bytes.Buffer
		if err := tmpl.Execute(&buf, cfg); err != nil {
			return fmt.Errorf("failed to execute template for %s: %w", f.Path, err)
		}

		if err := os.WriteFile(f.Path, buf.Bytes(), 0644); err != nil {
			return fmt.Errorf("failed to write file %s: %w", f.Path, err)
		}
	}

	// Try to run go fmt on the generated code
	cmd := exec.Command("C:\\Program Files\\Go\\bin\\go.exe", "fmt", "./...")
	cmd.Dir = outDir
	_ = cmd.Run() // Ignore errors if go fmt fails

	return nil
}
