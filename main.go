// Copyright 2015 Arne Roomann-Kurrik
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"flag"
	"fmt"
	"github.com/kurrik/fauxfile"
	"github.com/kurrik/tmxscripter/tmxscripter"
	"os"
)

func main() {
	var (
		err      error
		scripter = tmxscripter.NewTmxScripter(&fauxfile.RealFilesystem{})
	)
	flag.StringVar(&scripter.InputPath, "input", "", "Input path")
	flag.StringVar(&scripter.OutputPath, "output", "", "Output path")
	flag.StringVar(&scripter.ScriptPath, "script", "", "Script file")
	flag.Parse()
	if err = scripter.Run(); err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}

}
