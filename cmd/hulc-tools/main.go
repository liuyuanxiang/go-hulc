package main

import (
	"log"

	"github.com/liuyuanxiang/go-hulc/cmd/hulc-tools/v0.1/internal/project"
	"github.com/liuyuanxiang/go-hulc/cmd/hulc-tools/v0.1/internal/upgrade"

	"github.com/spf13/cobra"
)

var (
	version string = "v2.0.0-rc1"

	rootCmd = &cobra.Command{
		Use:     "hulk-tools",
		Short:   "hulk-tools: An elegant toolkit for Go microservices.",
		Long:    `hulk-tools: An elegant toolkit for Go microservices.`,
		Version: version,
	}
)

func init() {
	rootCmd.AddCommand(project.CmdNew)
	rootCmd.AddCommand(upgrade.CmdUpgrade)
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}
