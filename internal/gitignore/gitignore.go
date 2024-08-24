// internal/gitignore/gitignore.go
package gitignore

import (
	"path/filepath"

	"github.com/denormal/go-gitignore"
	"github.com/rs/zerolog/log"
)

type GitIgnore struct {
	ignoreMatcher gitignore.GitIgnore // Use the interface type
	rootDir       string
}

// NewGitIgnore initializes a GitIgnore structure that considers all .gitignore
// files from the root directory upwards in the directory hierarchy.
func NewGitIgnore(rootDir string) (*GitIgnore, error) {
	// Attempt to create a new GitIgnore repository matcher for the entire repository
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

// ShouldIgnore checks if a given path should be ignored based on .gitignore patterns
func (g *GitIgnore) ShouldIgnore(path string) bool {
	if g.ignoreMatcher == nil {
		return false
	}

	relPath, err := filepath.Rel(g.rootDir, path)
	if err != nil {
		log.Error().Err(err).Str("path", path).Msg("Failed to get relative path")
		return false
	}

	// Use the Matcher method to check if the path should be ignored
	match := g.ignoreMatcher.Match(relPath)
	if match != nil && match.Ignore() {
		return true
	}

	return false
}
