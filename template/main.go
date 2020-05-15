// Copyright 2020 Drone.IO Inc. All rights reserved.
// Use of this source code is governed by the Blue Oak Model License
// that can be found in the LICENSE file.

package main

import (
	"{{ Module }}/command"
	_ "github.com/joho/godotenv/autoload"
)

func main() {
	command.Command()
}
