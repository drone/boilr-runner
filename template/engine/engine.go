// Copyright 2020 Drone.IO Inc. All rights reserved.
// Use of this source code is governed by the Blue Oak Model License
// that can be found in the LICENSE file.

package engine

import (
	"context"
	"io"

	"github.com/drone/runner-go/pipeline/runtime"
)

// Opts configures the Engine.
type Opts struct {
	// TODO replace or remove
	Param1 string
	Param2 string
}

// Engine implements a pipeline engine.
type Engine struct {
	// TODO replace or remove
	Param1 string
	Param2 string
}

// New returns a new engine.
func New(opts Opts) (*Engine, error) {
	return &Engine{
		// TODO replace or remove
		Param1: opts.Param1,
		Param2: opts.Param2,
	}, nil
}

// Setup the pipeline environment.
func (e *Engine) Setup(ctx context.Context, specv runtime.Spec) error {
	// spec := specv.(*Spec)

	// TODO setup any resources for this pipeline to execute.
	// For example, the digitalocean runner creates a virtual
	// machine here.

	return nil
}

// Destroy the pipeline environment.
func (e *Engine) Destroy(ctx context.Context, specv runtime.Spec) error {
	// spec := specv.(*Spec)

	// TODO teardown any resources created for this pipeline.
	// For example, the digitalocean runner terminates the
	// virtual machiune here.

	return nil
}

// Run runs the pipeline step.
func (e *Engine) Run(ctx context.Context, specv runtime.Spec, stepv runtime.Step, output io.Writer) (*runtime.State, error) {
	// spec := specv.(*Spec)
	// step := stepv.(*Step)

	// TODO execute the pipeline step
	// TODO write the pipeine step output to the io.Writer
	// TODO return the pipeline step results

	return &runtime.State{
		ExitCode: 0,
		Exited:   true,
	}, nil
}

// Ping pings the underlying runtime to verify connectivity.
func (e *Engine) Ping(ctx context.Context) error {
	// TODO optionally add code to ping the underlying
	// service to verify credentials or connectivity. For
	// example, the kubernetes runner might ping kubernetes
	// to ensure the client can connect and is authorized
	// to make requests.
	return nil
}
