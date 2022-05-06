package main

import (
	"fmt"

	"github.com/spf13/cobra"
)

var (
	checkPaths       []string
	checkLevelFilter string
	checkRuleFilter  []string
	checkThreshold   int
)

var checkCmd = &cobra.Command{
	Use:   "check",
	Short: "Query the results of a SARIF",
	RunE: func(cmd *cobra.Command, args []string) error {
		var filtered []*DenormalizedResult
		for _, f := range checkPaths {
			partial, err := query(f, checkRuleFilter, checkLevelFilter)
			if err != nil {
				return err
			}
			filtered = append(filtered, partial...)
		}

		if len(filtered) > checkThreshold {
			return fmt.Errorf("Found %d results", len(filtered))
		}

		return nil
	},
}

func init() {
	checkCmd.Flags().StringArrayVarP(&checkPaths, "file", "f", nil, "File path(s) to the SARIF")
	checkCmd.MarkFlagRequired("file")

	checkCmd.Flags().IntVarP(&checkThreshold, "threshold", "t", 0, "Only fail the check if len(results) > threshold")
	checkCmd.Flags().StringArrayVarP(&checkRuleFilter, "rule", "r", nil, "Filter results by rule(s). By default it selects all.")
	checkCmd.Flags().StringVarP(&checkLevelFilter, "level", "l", "", "Filter results by severity level")
	// TODO: add other filtering, e.g. tool
}
