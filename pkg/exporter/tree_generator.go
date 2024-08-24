// File: pkg/exporter/tree_generator.go

package exporter

import (
	"bytes"
	"fmt"
	"os"
	"strings"

	"github.com/ddddddO/gtree"
)

type TreeGenerator struct{}

func NewTreeGenerator() *TreeGenerator {
	return &TreeGenerator{}
}

func (tg *TreeGenerator) Generate(files []FileInfo) (string, error) {
	if len(files) == 0 {
		return "## Repository Structure\n\n(empty)", nil
	}

	var (
		root *gtree.Node
		node *gtree.Node
	)

	for _, file := range files {
		splited := strings.Split(file.Path, string(os.PathSeparator))

		for i, s := range splited {
			if root == nil {
				root = gtree.NewRoot(s)
				node = root
				continue
			}
			if i == 0 {
				continue
			}

			tmp := node.Add(s)
			node = tmp
		}
		node = root
	}

	var buffer bytes.Buffer
	buffer.WriteString("## Repository Structure\n\n")
	if err := gtree.OutputProgrammably(&buffer, root); err != nil {
		return "", fmt.Errorf("failed to generate tree structure: %w", err)
	}

	return buffer.String(), nil
}
