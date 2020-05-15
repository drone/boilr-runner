// Copyright 2020 Drone.IO Inc. All rights reserved.
// Use of this source code is governed by the Blue Oak Model License
// that can be found in the LICENSE file.

package compiler

import (
	"context"
	"fmt"

	"{{ Module }}/engine"
	"{{ Module }}/engine/resource"

	"github.com/drone/runner-go/clone"
	"github.com/drone/runner-go/environ"
	"github.com/drone/runner-go/environ/provider"
	"github.com/drone/runner-go/manifest"
	"github.com/drone/runner-go/pipeline/runtime"
	"github.com/drone/runner-go/secret"

	"github.com/dchest/uniuri"
	"github.com/gosimple/slug"
)

// random generator function
var random = func() string {
	return "drone-" + uniuri.NewLen(20)
}

// Settings defines default settings.
type Settings struct {
	// TODO replace or remove
	Param1 string
	Param2 string
}

// Compiler compiles the Yaml configuration file to an
// intermediate representation optimized for simple execution.
type Compiler struct {
	// Environ provides a set of environment variables that
	// should be added to each pipeline step by default.
	Environ provider.Provider

	// Secret returns a named secret value that can be injected
	// into the pipeline step.
	Secret secret.Provider

	// Settings provides global settings that apply to
	// all pipelines.
	Settings Settings
}

// Compile compiles the configuration file.
func (c *Compiler) Compile(ctx context.Context, args runtime.CompilerArgs) runtime.Spec {
	pipeline := args.Pipeline.(*resource.Pipeline)
	os := pipeline.Platform.OS

	spec := &engine.Spec{
		Platform: engine.Platform{
			OS:      pipeline.Platform.OS,
			Arch:    pipeline.Platform.Arch,
			Variant: pipeline.Platform.Variant,
			Version: pipeline.Platform.Version,
		},
		Settings: engine.Settings{
			// TODO replace or remove
			Param1: c.Settings.Param1,
			Param2: c.Settings.Param2,
		},
	}

	// IMPORTANT:
	// this pipeline starter project is optimized for pipelines
	// that execute all steps on the same host. It is not optimized
	// for pipelines that execute each step inside a different
	// linux container. If you are building a container-based pipeline
	// please see the Docker runner source for inspiration.

	// create the root directory
	spec.Root = tempdir(os)

	// creates a home directory in the root.
	// note: mkdirall fails on windows so we need to create all
	// directories in the tree.
	homedir := join(os, spec.Root, "home", "drone")
	spec.Files = append(spec.Files, &engine.File{
		Path:  join(os, spec.Root, "home"),
		Mode:  0700,
		IsDir: true,
	})
	spec.Files = append(spec.Files, &engine.File{
		Path:  homedir,
		Mode:  0700,
		IsDir: true,
	})

	// creates a source directory in the root.
	// note: mkdirall fails on windows so we need to create all
	// directories in the tree.
	sourcedir := join(os, spec.Root, "drone", "src")
	spec.Files = append(spec.Files, &engine.File{
		Path:  join(os, spec.Root, "drone"),
		Mode:  0700,
		IsDir: true,
	})
	spec.Files = append(spec.Files, &engine.File{
		Path:  sourcedir,
		Mode:  0700,
		IsDir: true,
	})

	// creates the opt directory to hold all scripts.
	spec.Files = append(spec.Files, &engine.File{
		Path:  join(os, spec.Root, "opt"),
		Mode:  0700,
		IsDir: true,
	})

	// creates the netrc file
	if args.Netrc != nil && args.Netrc.Password != "" {
		netrcfile := getNetrc(os)
		netrcpath := join(os, homedir, netrcfile)
		netrcdata := fmt.Sprintf(
			"machine %s login %s password %s",
			args.Netrc.Machine,
			args.Netrc.Login,
			args.Netrc.Password,
		)
		spec.Files = append(spec.Files, &engine.File{
			Path: netrcpath,
			Mode: 0600,
			Data: []byte(netrcdata),
		})
	}

	// list the global environment variables
	globals, _ := c.Environ.List(ctx, &provider.Request{
		Build: args.Build,
		Repo:  args.Repo,
	})

	// create the default environment variables.
	envs := environ.Combine(
		provider.ToMap(
			provider.FilterUnmasked(globals),
		),
		args.Build.Params,
		pipeline.Environment,
		environ.Proxy(),
		environ.System(args.System),
		environ.Repo(args.Repo),
		environ.Build(args.Build),
		environ.Stage(args.Stage),
		environ.Link(args.Repo, args.Build, args.System),
		clone.Environ(clone.Config{
			SkipVerify: pipeline.Clone.SkipVerify,
			Trace:      pipeline.Clone.Trace,
			User: clone.User{
				Name:  args.Build.AuthorName,
				Email: args.Build.AuthorEmail,
			},
		}),
		map[string]string{
			"HOME":                homedir,
			"HOMEPATH":            homedir, // for windows
			"USERPROFILE":         homedir, // for windows
			"DRONE_HOME":          sourcedir,
			"DRONE_WORKSPACE":     sourcedir,
			"GIT_TERMINAL_PROMPT": "0",
		},
	)

	match := manifest.Match{
		Action:   args.Build.Action,
		Cron:     args.Build.Cron,
		Ref:      args.Build.Ref,
		Repo:     args.Repo.Slug,
		Instance: args.System.Host,
		Target:   args.Build.Deploy,
		Event:    args.Build.Event,
		Branch:   args.Build.Target,
	}

	// create the clone step, maybe
	if pipeline.Clone.Disable == false {
		clonepath := join(os, spec.Root, "opt", getExt(os, "clone"))
		clonefile := genScript(os,
			clone.Commands(
				clone.Args{
					Branch: args.Build.Target,
					Commit: args.Build.After,
					Ref:    args.Build.Ref,
					Remote: args.Repo.HTTPURL,
				},
			),
		)

		cmd, args := getCommand(os, clonepath)
		spec.Steps = append(spec.Steps, &engine.Step{
			Name:      "clone",
			Args:      args,
			Command:   cmd,
			Envs:      envs,
			RunPolicy: runtime.RunAlways,
			Files: []*engine.File{
				{
					Path: clonepath,
					Mode: 0700,
					Data: []byte(clonefile),
				},
			},
			Secrets:    []*engine.Secret{},
			WorkingDir: sourcedir,
		})
	}

	// create steps
	for _, src := range pipeline.Steps {
		buildslug := slug.Make(src.Name)
		buildpath := join(os, spec.Root, "opt", getExt(os, buildslug))
		buildfile := genScript(os, src.Commands)

		cmd, args := getCommand(os, buildpath)
		dst := &engine.Step{
			Name:      src.Name,
			Args:      args,
			Command:   cmd,
			Detach:    src.Detach,
			DependsOn: src.DependsOn,
			Envs: environ.Combine(envs,
				environ.Expand(
					convertStaticEnv(src.Environment),
				),
			),
			RunPolicy: runtime.RunOnSuccess,
			Files: []*engine.File{
				{
					Path: buildpath,
					Mode: 0700,
					Data: []byte(buildfile),
				},
			},
			Secrets:    convertSecretEnv(src.Environment),
			WorkingDir: sourcedir,
		}
		spec.Steps = append(spec.Steps, dst)

		// set the pipeline step run policy. steps run on
		// success by default, but may be optionally configured
		// to run on failure.
		if isRunAlways(src) {
			dst.RunPolicy = runtime.RunAlways
		} else if isRunOnFailure(src) {
			dst.RunPolicy = runtime.RunOnFailure
		}

		// if the pipeline step has unmet conditions the step is
		// automatically skipped.
		if !src.When.Match(match) {
			dst.RunPolicy = runtime.RunNever
		}
	}

	if isGraph(spec) == false {
		configureSerial(spec)
	} else if pipeline.Clone.Disable == false {
		configureCloneDeps(spec)
	} else if pipeline.Clone.Disable == true {
		removeCloneDeps(spec)
	}

	for _, step := range spec.Steps {
		for _, s := range step.Secrets {
			secret, ok := c.findSecret(ctx, args, s.Name)
			if ok {
				s.Data = []byte(secret)
			}
		}
	}

	return spec
}

// helper function attempts to find and return the named secret.
// from the secret provider.
func (c *Compiler) findSecret(ctx context.Context, args runtime.CompilerArgs, name string) (s string, ok bool) {
	if name == "" {
		return
	}
	// source secrets from the global secret provider
	// and the repository secret provider.
	provider := secret.Combine(
		args.Secret,
		c.Secret,
	)
	// fine the secret from the provider. please note we
	// currently ignore errors if the secret is not found,
	// which is something that we'll need to address in the
	// next major (breaking) release.
	found, _ := provider.Find(ctx, &secret.Request{
		Name:  name,
		Build: args.Build,
		Repo:  args.Repo,
		Conf:  args.Manifest,
	})
	if found == nil {
		return
	}
	return found.Data, true
}
