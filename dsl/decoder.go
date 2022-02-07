package dsl

import (
	"errors"
	"fmt"
	"path/filepath"

	"github.com/hashicorp/hcl/v2"
	"github.com/zclconf/go-cty/cty"
	"github.com/zclconf/go-cty/cty/function"
)

func Decode(body hcl.Body) (*Pipeline, error) {
	ctx := &hcl.EvalContext{
		Functions: map[string]function.Function{
			RandomID: Random,
		},
	}
	pipeline := &Pipeline{
		jobs: make([]Job, 0),
		vars: make(map[string]cty.Value),
	}
	attributes, _ := body.JustAttributes()
	for name, value := range attributes {
		v, d := value.Expr.Value(ctx)
		if d.HasErrors() {
			fmt.Println(d.Errs()[0])
			continue
		}
		pipeline.vars[name] = v
	}
	ctx.Variables = pipeline.vars
	bc, _ := body.Content(schema)

	if len(bc.Blocks) == 0 {
		return nil, errors.New("at least one pipeline must be provided")
	}
	blocks := bc.Blocks.ByType()
	for blockName := range blocks {
		switch blockName {
		case jobID:
			for _, b := range blocks[blockName] {
				job := new(Job)
				err := job.FromHCLBlock(b, ctx)
				if err != nil {
					return nil, err
				}
				pipeline.jobs = append(pipeline.jobs, *job)
			}
		}
	}
	return pipeline, nil
}

func (j *Job) FromHCLBlock(block *hcl.Block, ctx *hcl.EvalContext) error {
	bc, d := block.Body.Content(jobSchema)
	if d.HasErrors() {
		return d.Errs()[0]
	}
	j.steps = make([]Step, 0)
	for _, subBlock := range bc.Blocks {
		switch subBlock.Type {
		case pythonID:
			step := new(Python)
			if err := step.FromHCLBlock(subBlock, ctx); err != nil {
				return err
			}
			j.steps = append(j.steps, step)
		case slackID:
			step := new(Slack)
			if err := step.FromHCLBlock(subBlock, ctx); err != nil {
				return err
			}
			j.steps = append(j.steps, step)
		}
	}
	j.name = block.Labels[0]
	j.description = block.Labels[1]
	return nil
}

func (p *Python) FromHCLBlock(block *hcl.Block, ctx *hcl.EvalContext) error {
	p.path = block.Labels[0]
	bc, d := block.Body.Content(pythonSchema)
	if d.HasErrors() {
		return d.Errs()[0]
	}
	if attr, ok := bc.Attributes["root_dir"]; ok {
		rootDir, d := attr.Expr.Value(ctx)
		if d.HasErrors() {
			return d.Errs()[0]
		}
		p.path = filepath.Join(rootDir.AsString(), p.path)
	}
	return nil
}

func (s *Slack) FromHCLBlock(block *hcl.Block, ctx *hcl.EvalContext) error {
	bc, d := block.Body.Content(slackSchema)
	if d.HasErrors() {
		return d.Errs()[0]
	}

	for name, attr := range bc.Attributes {

		switch name {
		case "channel":
			value, d := attr.Expr.Value(ctx)
			if d.HasErrors() {
				continue
			}
			s.channelID = value.AsString()
		case "message":
			value, d := attr.Expr.Value(ctx)
			if d.HasErrors() {
				continue
			}
			s.message = value.AsString()
		}
	}
	return nil
}
