package projectgen

import (
	"bytes"
	"embed"
	"fmt"
	"io/fs"
	"path/filepath"
	"sort"
	"strings"
	"text/template"
)

//go:embed all:templates
var templateFS embed.FS

type templateData struct {
	Name        string
	Description string
	ModulePath  string
	GoVersion   string
}

func (g *Generator) renderTemplates(targetPath string, data templateData) error {
	templates, err := g.templatePaths()
	if err != nil {
		return err
	}
	g.log.Debug().Int("count", len(templates)).Msg("discovered templates")

	for _, templatePath := range templates {
		relativePath := outputPath(templatePath, data)
		g.log.Debug().
			Str("template", templatePath).
			Str("output_path", relativePath).
			Msg("rendering template")
		content, err := g.renderTemplate(templatePath, data)
		if err != nil {
			return fmt.Errorf("render file %q: %w", relativePath, err)
		}

		if err := writeFile(targetPath, relativePath, []byte(content)); err != nil {
			return err
		}
	}

	return nil
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

	tmpl, err := template.New(filepath.Base(path)).
		Funcs(template.FuncMap{
			"upper": strings.ToUpper,
		}).
		Parse(string(raw))
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
