package schema

import "github.com/buildkite/pipeline-sdk/pkg/schema_types"

type PipelineSchema struct {
	Version     string
	Name        string
	Types       []schema_types.Field
	Environment []EnvironmentVariable
	Steps       []Step
}
