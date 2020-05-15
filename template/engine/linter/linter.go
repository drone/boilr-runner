// Copyright 2020 Drone.IO Inc. All rights reserved.
// Use of this source code is governed by the Blue Oak Model License
// that can be found in the LICENSE file.

package linter

import (
	"errors"

	"{{ Module }}/engine/resource"
	"github.com/drone/drone-go/drone"
	"github.com/drone/runner-go/manifest"
)


// Linter evaluates the pipeline against a set of
// rules and returns an error if one or more of the
// rules are broken.
type Linter struct{}

// New returns a new Linter.
func New() *Linter {
	return new(Linter)
}

// Lint executes the linting rules for the pipeline
// configuration.
func (l *Linter) Lint(pipeline manifest.Resource, repo *drone.Repo) error {
	return checkPipeline(pipeline.(*resource.Pipeline), repo.Trusted)
}

func checkPipeline(pipeline *resource.Pipeline, trusted bool) error {
	if err := checkSteps(pipeline, trusted); err != nil {
		return err
	}
	return nil
}

func checkSteps(pipeline *resource.Pipeline, trusted bool) error {
	for _, step := range pipeline.Steps {
		if step == nil {
			return errors.New("Linter: nil step")
		}
		if err := checkStep(step, trusted); err != nil {
			return err
		}
	}
	return nil
}

func checkStep(step *resource.Step, trusted bool) error {
	// TODO optionally add per-step pipeline rules or remove
	return nil
}
