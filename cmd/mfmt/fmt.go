package mfmt

import (
	"github.com/YeHeng/go-web-api/pkg/cmd/mfmt"

	"github.com/spf13/cobra"
)

func NewFmtCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "mfmt",
		Short: "format go source files",
		Annotations: map[string]string{
			"IsActions": "true",
		},
		Run: func(cmd *cobra.Command, args []string) {
			mfmt.RunFmt()
		},
	}

	return cmd
}
