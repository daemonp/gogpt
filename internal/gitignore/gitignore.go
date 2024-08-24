package gitignore

import (
	"bufio"
	"os"
	"path/filepath"
	"strings"

	"github.com/rs/zerolog/log"
)

type GitIgnore struct {
	patterns []string
	rootDir  string
}

func NewGitIgnore(rootDir string) (*GitIgnore, error) {
	gi := &GitIgnore{rootDir: rootDir}
	err := gi.parseGitIgnoreFile(filepath.Join(rootDir, ".gitignore"))
	if err != nil {
		return nil, err
	}
	return gi, nil
}

func (gi *GitIgnore) parseGitIgnoreFile(path string) error {
	file, err := os.Open(path)
	if err != nil {
		if os.IsNotExist(err) {
			log.Warn().Str("path", path).Msg("No .gitignore file found")
			return nil
		}
		return err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line != "" && !strings.HasPrefix(line, "#") {
			gi.patterns = append(gi.patterns, line)
		}
	}

	return scanner.Err()
}

func (gi *GitIgnore) ShouldIgnore(path string) bool {
	relPath, err := filepath.Rel(gi.rootDir, path)
	if err != nil {
		log.Error().Err(err).Str("path", path).Msg("Failed to get relative path")
		return false
	}

	for _, pattern := range gi.patterns {
		matched, err := filepath.Match(pattern, relPath)
		if err != nil {
			log.Error().Err(err).Str("pattern", pattern).Str("path", relPath).Msg("Failed to match pattern")
			continue
		}
		if matched {
			return true
		}
	}

	return false
}
