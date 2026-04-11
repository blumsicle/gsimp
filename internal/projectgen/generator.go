// Package projectgen renders starter project scaffolds from embedded templates.
package projectgen

import (
	"bytes"
	"context"
	"embed"
	"fmt"
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"sort"
	"strings"
	"text/template"

	scaffoldsrc "github.com/blumsicle/bcli"
	cliutil "github.com/blumsicle/bcli/internal/cli"
	"github.com/blumsicle/bcli/internal/poststep"
	"github.com/rs/zerolog"
)

//go:embed all:templates
var templateFS embed.FS

// Config defines the inputs used to generate a new project scaffold.
type Config struct {
	Name             string
	Description      string
	GitLocation      string
	ProjectDirPrefix string
	RootPath         string
}

type templateData struct {
	Name        string
	Description string
	ModulePath  string
	GoVersion   string
}

type generationPlan struct {
	rootPath      string
	targetPath    string
	modulePath    string
	templateData  templateData
	postStepInput poststep.PostStepInput
}

type staticFile struct {
	sourcePath string
	outputPath string
}

var staticFiles = []staticFile{
	{
		sourcePath: "internal/appconfig/load.go",
		outputPath: "internal/appconfig/load.go",
	},
	{
		sourcePath: "internal/cli/buildinfo.go",
		outputPath: "internal/cli/buildinfo.go",
	},
	{
		sourcePath: "internal/cli/runner.go",
		outputPath: "internal/cli/runner.go",
	},
}

// Generator renders embedded project templates and runs post-generation steps.
type Generator struct {
	templateFS fs.FS
	postSteps  []poststep.PostStep
	log        zerolog.Logger
}

// New constructs a generator with the embedded template filesystem.
func New(log zerolog.Logger) *Generator {
	return &Generator{
		templateFS: templateFS,
		postSteps:  []poststep.PostStep{},
		log:        cliutil.SubLogger(log, "projectgen"),
	}
}

// AddPostStep appends a post-generation step to be run after rendering completes.
func (g *Generator) AddPostStep(step poststep.PostStep) {
	g.log.Debug().Str("step", step.Name()).Msg("registered post step")
	g.postSteps = append(g.postSteps, step)
}

// Generate renders the scaffold into the target directory and runs post steps.
func (g *Generator) Generate(ctx context.Context, cfg Config) (string, error) {
	if err := validateConfig(cfg); err != nil {
		return "", err
	}

	plan := newGenerationPlan(cfg)
	g.log.Info().
		Str("name", cfg.Name).
		Str("project_dir_prefix", cfg.ProjectDirPrefix).
		Str("root_path", filepath.Clean(plan.rootPath)).
		Str("target_path", filepath.Clean(plan.targetPath)).
		Msg("starting project generation")

	g.log.Debug().
		Str("name", cfg.Name).
		Str("module_path", plan.modulePath).
		Str("go_version", plan.templateData.GoVersion).
		Msg("resolved module path")

	if err := ensureTargetDir(plan.targetPath); err != nil {
		return "", err
	}
	if err := g.renderTemplates(plan.targetPath, plan.templateData); err != nil {
		return "", err
	}
	if err := g.copyStaticFiles(plan.targetPath); err != nil {
		return "", err
	}
	if err := g.runPostSteps(ctx, plan.postStepInput); err != nil {
		return "", err
	}

	g.log.Info().
		Str("name", cfg.Name).
		Str("target_path", filepath.Clean(plan.targetPath)).
		Msg("finished project generation")

	return plan.targetPath, nil
}

func validateConfig(cfg Config) error {
	if cfg.Name == "" {
		return fmt.Errorf("name is required")
	}
	if cfg.Description == "" {
		return fmt.Errorf("description is required")
	}

	return nil
}

func newGenerationPlan(cfg Config) generationPlan {
	rootPath := cfg.RootPath
	if rootPath == "" {
		rootPath = "."
	}

	targetPath := filepath.Join(rootPath, cfg.ProjectDirPrefix+cfg.Name)
	modulePath := modulePath(cfg.GitLocation, cfg.Name)
	data := templateData{
		Name:        cfg.Name,
		Description: cfg.Description,
		ModulePath:  modulePath,
		GoVersion:   currentGoVersion(),
	}

	return generationPlan{
		rootPath:     rootPath,
		targetPath:   targetPath,
		modulePath:   modulePath,
		templateData: data,
		postStepInput: poststep.PostStepInput{
			ProjectPath: targetPath,
			Name:        cfg.Name,
			ModulePath:  modulePath,
		},
	}
}

func ensureTargetDir(targetPath string) error {
	if _, err := os.Stat(targetPath); err == nil {
		return fmt.Errorf("target path already exists: %s", targetPath)
	} else if !os.IsNotExist(err) {
		return fmt.Errorf("stat target path: %w", err)
	}

	if err := os.MkdirAll(targetPath, 0o755); err != nil {
		return fmt.Errorf("create target directory: %w", err)
	}

	return nil
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

func (g *Generator) copyStaticFiles(targetPath string) error {
	for _, file := range staticFiles {
		g.log.Debug().
			Str("source_path", file.sourcePath).
			Str("output_path", file.outputPath).
			Msg("copying static file")
		content, err := fs.ReadFile(scaffoldsrc.ScaffoldSourceFS, file.sourcePath)
		if err != nil {
			return fmt.Errorf("read static file %q: %w", file.sourcePath, err)
		}

		if err := writeFile(targetPath, file.outputPath, content); err != nil {
			return err
		}
	}

	return nil
}

func (g *Generator) runPostSteps(ctx context.Context, input poststep.PostStepInput) error {
	for _, step := range g.postSteps {
		g.log.Info().Str("step", step.Name()).Msg("running post step")
		if err := step.Run(ctx, input); err != nil {
			return fmt.Errorf("run post step %q: %w", step.Name(), err)
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

func writeFile(targetPath string, relativePath string, content []byte) error {
	fullPath := filepath.Join(targetPath, relativePath)
	if err := os.MkdirAll(filepath.Dir(fullPath), 0o755); err != nil {
		return fmt.Errorf("create parent directory for %q: %w", relativePath, err)
	}
	if err := os.WriteFile(fullPath, content, 0o644); err != nil {
		return fmt.Errorf("write %q: %w", relativePath, err)
	}

	return nil
}

func modulePath(gitLocation string, name string) string {
	gitLocation = strings.TrimSuffix(gitLocation, "/")
	if gitLocation == "" {
		return name
	}
	return gitLocation + "/" + name
}

func currentGoVersion() string {
	version, err := exec.Command("go", "env", "GOVERSION").Output()
	if err == nil {
		goVersion := strings.TrimSpace(string(version))
		goVersion = strings.TrimPrefix(goVersion, "go")
		if goVersion != "" {
			return goVersion
		}
	}

	return strings.TrimPrefix(runtime.Version(), "go")
}
