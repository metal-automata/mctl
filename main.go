/*
Copyright © 2022-2024 Equinix Metal <EMAIL ADDRESS>
Copyright 2024 metal-automata https://github.com/metal-automata

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

	http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package main

import (
	"github.com/metal-automata/mctl/cmd"
	_ "github.com/metal-automata/mctl/cmd/bios"
	_ "github.com/metal-automata/mctl/cmd/collect"
	_ "github.com/metal-automata/mctl/cmd/create"
	_ "github.com/metal-automata/mctl/cmd/delete"
	_ "github.com/metal-automata/mctl/cmd/edit"
	_ "github.com/metal-automata/mctl/cmd/generate"
	_ "github.com/metal-automata/mctl/cmd/get"
	_ "github.com/metal-automata/mctl/cmd/install"
	_ "github.com/metal-automata/mctl/cmd/list"
	_ "github.com/metal-automata/mctl/cmd/power"
)

func main() {
	cmd.Execute()
}
