package projectgen

import (
	"bytes"
	"embed"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"text/template"
)

//go:embed all:templates
var templateFS embed.FS

type Config struct {
	Name        string
	Description string
	GitLocation string
	RootPath    string
}

type templateData struct {
	Name        string
	Description string
	ModulePath  string
}

type Generator struct {
	templateFS fs.FS
}

func New() *Generator {
	return &Generator{templateFS: templateFS}
}

func (g *Generator) Generate(cfg Config) (string, error) {
	if cfg.Name == "" {
		return "", fmt.Errorf("name is required")
	}
	if cfg.Description == "" {
		return "", fmt.Errorf("description is required")
	}

	rootPath := cfg.RootPath
	if rootPath == "" {
		rootPath = "."
	}

	targetPath := filepath.Join(rootPath, cfg.Name)
	if _, err := os.Stat(targetPath); err == nil {
		return "", fmt.Errorf("target path already exists: %s", targetPath)
	} else if !os.IsNotExist(err) {
		return "", fmt.Errorf("stat target path: %w", err)
	}

	if err := os.MkdirAll(targetPath, 0o755); err != nil {
		return "", fmt.Errorf("create target directory: %w", err)
	}

	data := templateData{
		Name:        cfg.Name,
		Description: cfg.Description,
		ModulePath:  modulePath(cfg.GitLocation, cfg.Name),
	}

	templates, err := g.templatePaths()
	if err != nil {
		return "", err
	}

	for _, templatePath := range templates {
		relativePath := outputPath(templatePath, data)
		content, err := g.renderTemplate(templatePath, data)
		if err != nil {
			return "", fmt.Errorf("render file %q: %w", relativePath, err)
		}

		fullPath := filepath.Join(targetPath, relativePath)
		if err := os.MkdirAll(filepath.Dir(fullPath), 0o755); err != nil {
			return "", fmt.Errorf("create parent directory for %q: %w", relativePath, err)
		}
		if err := os.WriteFile(fullPath, []byte(content), 0o644); err != nil {
			return "", fmt.Errorf("write %q: %w", relativePath, err)
		}
	}

	return targetPath, nil
}

func (g *Generator) templatePaths() ([]string, error) {
	var paths []string

	err := fs.WalkDir(g.templateFS, "templates", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			return nil
		}
		if !strings.HasSuffix(path, ".tmpl") {
			return nil
		}

		paths = append(paths, path)
		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("list templates: %w", err)
	}

	sort.Strings(paths)
	return paths, nil
}

func (g *Generator) renderTemplate(path string, data templateData) (string, error) {
	raw, err := fs.ReadFile(g.templateFS, path)
	if err != nil {
		return "", err
	}

	tmpl, err := template.New(filepath.Base(path)).Parse(string(raw))
	if err != nil {
		return "", err
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		return "", err
	}

	return buf.String(), nil
}

func outputPath(templatePath string, data templateData) string {
	path := strings.TrimPrefix(templatePath, "templates/")
	path = strings.TrimSuffix(path, ".tmpl")
	path = strings.ReplaceAll(path, "__NAME__", data.Name)
	return path
}

func modulePath(gitLocation string, name string) string {
	gitLocation = strings.TrimSuffix(gitLocation, "/")
	if gitLocation == "" {
		return name
	}
	return gitLocation + "/" + name
}
