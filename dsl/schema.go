package dsl

import "github.com/hashicorp/hcl/v2"

const slackID = "slack"

// Configuration for block slack
var (
	slackLabels = []string{}

	slackSchema = &hcl.BodySchema{
		Blocks: []hcl.BlockHeaderSchema{},
		Attributes: []hcl.AttributeSchema{
			{Name: "channel", Required: false},
			{Name: "message", Required: true},
		},
	}
)

// Configuration for block python
const pythonID = "python"

var (
	pythonLabels = []string{"script"}
	pythonSchema = &hcl.BodySchema{
		Blocks: []hcl.BlockHeaderSchema{},
		Attributes: []hcl.AttributeSchema{
			{Name: "root_dir", Required: false},
		},
	}
)

// Configuration for block job
const jobID = "job"

var (
	jobLabels = []string{"name", "description"}
	jobSchema = &hcl.BodySchema{
		Blocks: []hcl.BlockHeaderSchema{

			{
				Type:       pythonID,
				LabelNames: pythonLabels,
			},
			{
				Type:       slackID,
				LabelNames: slackLabels,
			},
		},
		Attributes: []hcl.AttributeSchema{},
	}
)

// Configuration for the full schema
var schema = &hcl.BodySchema{
	Attributes: []hcl.AttributeSchema{},
	Blocks: []hcl.BlockHeaderSchema{
		{
			Type:       jobID,
			LabelNames: jobLabels,
		},
	},
}
