package dsl

import (
	"fmt"

	"github.com/hashicorp/hcl/v2"
	"github.com/zclconf/go-cty/cty"
	"github.com/zclconf/go-cty/cty/function"
)

/// Slack
type Slack struct {
	channelID string
	message   string
}

func (s *Slack) Run(ctx *hcl.EvalContext) error {
	fmt.Printf(">> sendnotification: '%s' => channel: %s \n",
		s.message, s.channelID)
	return nil
}

// Pyrhon
type Python struct {
	path string
}

func (p *Python) Run(ctx *hcl.EvalContext) error {
	fmt.Printf(">> python %s \n", p.path)
	return nil
}

// Job
type Job struct {
	name        string
	description string
	steps       []Step
}

type Step interface {
	Run(ctx *hcl.EvalContext) error
}

func (j *Job) run(ctx *hcl.EvalContext) error {
	fmt.Printf("----------\njob: %s, %s\n\n", j.name, j.description)
	for _, step := range j.steps {
		step.Run(ctx)
	}
	fmt.Println("----------")
	return nil
}

//
type Pipeline struct {
	jobs []Job
	vars map[string]cty.Value
}

func (p *Pipeline) Run() {
	ctx := &hcl.EvalContext{
		Functions: map[string]function.Function{
			RandomID: Random,
		},
		Variables: p.vars,
	}
	for _, job := range p.jobs {
		job.run(ctx)
	}
}
