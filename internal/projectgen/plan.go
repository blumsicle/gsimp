package projectgen

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/blumsicle/bcli/internal/poststep"
)

type generationPlan struct {
	rootPath      string
	targetPath    string
	modulePath    string
	templateData  templateData
	postStepInput poststep.PostStepInput
}

func newGenerationPlan(cfg Config) generationPlan {
	rootPath := cfg.RootPath
	if rootPath == "" {
		rootPath = "."
	}

	targetPath := filepath.Join(rootPath, cfg.ProjectDirPrefix+cfg.Name)
	if cfg.InPlace {
		workingDir, err := os.Getwd()
		if err == nil {
			targetPath = workingDir
			rootPath = workingDir
		} else {
			targetPath = "."
			rootPath = "."
		}
	}
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

func modulePath(gitLocation string, name string) string {
	gitLocation = strings.TrimSuffix(gitLocation, "/")
	if gitLocation == "" {
		return name
	}
	return gitLocation + "/" + name
}
