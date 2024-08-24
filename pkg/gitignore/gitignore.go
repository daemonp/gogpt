package gitignore

import (
	"path/filepath"

	"github.com/denormal/go-gitignore"
	"github.com/rs/zerolog/log"
)

type GitIgnore struct {
	ignoreMatcher gitignore.GitIgnore
	rootDir       string
}

func NewGitIgnore(rootDir string) (*GitIgnore, error) {
	ignoreMatcher, err := gitignore.NewRepository(rootDir)
	if err != nil {
		log.Error().Err(err).Msg("Failed to parse .gitignore files")
		return nil, err
	}

	return &GitIgnore{
		ignoreMatcher: ignoreMatcher,
		rootDir:       rootDir,
	}, nil
}

func (g *GitIgnore) ShouldIgnore(path string) bool {
	if g.ignoreMatcher == nil {
		return false
	}

	relPath, err := filepath.Rel(g.rootDir, path)
	if err != nil {
		log.Error().Err(err).Str("path", path).Msg("Failed to get relative path")
		return false
	}

	match := g.ignoreMatcher.Match(relPath)
	return match != nil && match.Ignore()
}
