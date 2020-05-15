// Copyright 2020 Drone.IO Inc. All rights reserved.
// Use of this source code is governed by the Blue Oak Model License
// that can be found in the LICENSE file.

package command

import (
	"context"
	"os"

	"{{ Module }}/command/daemon"

	"gopkg.in/alecthomas/kingpin.v2"
)

// program version
var version = "0.0.0"

// empty context
var nocontext = context.Background()

// Command parses the command line arguments and then executes a
// subcommand program.
func Command() {
	app := kingpin.New("drone", "drone {{ Type }} runner")
	registerCompile(app)
	registerExec(app)
	daemon.Register(app)

	kingpin.Version(version)
	kingpin.MustParse(app.Parse(os.Args[1:]))
}
