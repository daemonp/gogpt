// File: cmd/gogpt/mock_exporter.go

package main

import (
	"github.com/daemonp/gogpt/pkg/exporter"
	"github.com/daemonp/gogpt/pkg/types"
)

type TestExporter struct {
	*exporter.Exporter
	ExportFunc func() error
}

func (te *TestExporter) Export() error {
	if te.ExportFunc != nil {
		return te.ExportFunc()
	}
	return te.Exporter.Export()
}

func NewTestExporter(rootDir string, flags *types.Flags) (*exporter.Exporter, error) {
	return &exporter.Exporter{}, nil
}
