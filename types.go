package main

type Envelope struct {
	Version string `json:"version"`
	Schema  string `json:"$schema,omitempty"`
	Runs    []*Run `json:"runs"`
}

type DenormalizedResult struct {
	DriverName string
	Rule       *CoreRule
	Result     *CoreResult
}

func (e *Envelope) Results() []*DenormalizedResult {
	var rs []*DenormalizedResult
	for _, run := range e.Runs {
		for _, result := range run.Results {
			r := run.Tool.Driver.Rules[result.RuleIndex]
			rs = append(rs, &DenormalizedResult{
				DriverName: run.Tool.Driver.Name,
				Rule:       r,
				Result:     result,
			})
		}
	}
	return rs
}

type Run struct {
	Tool    Tool          `json:"tool"`
	Results []*CoreResult `json:"results"`
}

type CoreResult struct {
	RuleID              string                 `json:"ruleId,omitempty"`
	RuleIndex           uint                   `json:"ruleIndex,omitempty"`
	Message             Message                `json:"message"`
	Level               string                 `json:"level,omitempty"`
	PartialFingerprints map[string]interface{} `json:"partialFingerprints,omitempty"`
}

type Tool struct {
	Driver *CoreDriver `json:"driver"`
}

type CoreDriver struct {
	Name            string      `json:"name"`
	InformationURL  string      `json:"informationUri,omitempty"`
	SemanticVersion string      `json:"semanticVersion,omitempty"`
	Rules           []*CoreRule `json:"rules"`
}

type CoreRule struct {
	ID                   string                  `json:"id"`
	Name                 string                  `json:"name,omitempty"`
	ShortDescription     *MultiformatText        `json:"shortDescription"`
	FullDescription      *MultiformatText        `json:"fullDescription,omitempty"`
	DefaultConfiguration *ReportingConfiguration `json:"defaultConfiguration,omitempty"`
	HelpURI              string                  `json:"helpUri,omitempty"`
	Help                 *MultiformatText        `json:"help,omitempty"`
	Properties           map[string]interface{}  `json:"properties,omitempty"`
}

type MultiformatText struct {
	Text     string `json:"text,omitempty"`
	Markdown string `json:"markdown,omitempty"`
}

type ReportingConfiguration struct {
	Enabled bool    `json:"enabled,omitempty"`
	Level   string  `json:"level,omitempty"`
	Rank    float64 `json:"rank,omitempty"`
}

type Message struct {
	Text      string   `json:"text,omitempty"`
	Markdown  string   `json:"markdown,omitempty"`
	ID        string   `json:"id,omitempty"`
	Arguments []string `json:"arguments,omitempty"`
}
