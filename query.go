package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"

	"github.com/spf13/cobra"
)

const (
	defaultLevel = "warning"
)

var (
	queryVerbosity   string
	queryPath        string
	queryLevelFilter string
	queryRuleFilter  []string
)

var queryCmd = &cobra.Command{
	Use:   "query",
	Short: "Query the results of a SARIF",
	RunE: func(cmd *cobra.Command, args []string) error {
		filtered, err := query(queryPath, queryRuleFilter, queryLevelFilter)
		if err != nil {
			return err
		}

		switch queryVerbosity {
		case "default":
			for _, res := range filtered {
				// TODO: use pretty table output
				fmt.Fprintf(cmd.OutOrStdout(), "Driver: %s, Rule: %s, Message: %s\n\n", res.DriverName, res.Rule.ID, res.Result.Message.Text)
			}
		case "count":
			fmt.Fprintf(cmd.OutOrStderr(), "%d", len(filtered))
		}

		return nil
	},
}

func query(f string, rules []string, level string) ([]*DenormalizedResult, error) {
	b, err := ioutil.ReadFile(f)
	if err != nil {
		return nil, fmt.Errorf("failed to read SARIF file: %w", err)
	}

	var report Envelope
	if err := json.Unmarshal(b, &report); err != nil {
		return nil, fmt.Errorf("failed to unmarshal SARIF file: %w", err)
	}

	results := report.Results()
	var filtered []*DenormalizedResult
	for _, res := range results {
		if !ruleMatch(res, rules) || !levelMatch(res, level) {
			continue
		}
		filtered = append(filtered, res)
	}

	return filtered, nil
}

func ruleMatch(res *DenormalizedResult, rules []string) bool {
	if len(rules) == 0 {
		return true
	}

	for _, r := range rules {
		if res.Rule.ID == r {
			return true
		}
	}
	return false
}

func levelMatch(res *DenormalizedResult, level string) bool {
	if level == "" {
		return true
	}
	l := defaultLevel
	if res.Rule.DefaultConfiguration != nil && res.Rule.DefaultConfiguration.Level != "" {
		l = res.Rule.DefaultConfiguration.Level
	}
	if res.Result.Level != "" {
		l = res.Result.Level
	}
	return l == level
}

func init() {
	queryCmd.Flags().StringVarP(&queryVerbosity, "verbosity", "v", "default", "Query result verbosity level. Allowed values are [default, count]")
	queryCmd.Flags().StringVarP(&queryPath, "file", "f", "", "File path to the SARIF")
	queryCmd.MarkFlagRequired("file")

	queryCmd.Flags().StringArrayVarP(&queryRuleFilter, "rule", "r", nil, "Filter results by rule(s). By default it selects all.")
	queryCmd.Flags().StringVarP(&queryLevelFilter, "level", "l", "", "Filter results by severity level")
	// TODO: add other filtering, e.g. tool
}
