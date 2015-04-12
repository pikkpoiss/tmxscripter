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
	"fmt"
	"github.com/kurrik/fauxfile"
	"github.com/kurrik/tmxgo"
	"io/ioutil"
	"testing"
)

func getGrid(input string, layer string) (output tmxgo.DataTileGrid, err error) {
	var (
		m *tmxgo.Map
		l *tmxgo.Layer
	)
	if m, err = tmxgo.ParseMapString(input); err != nil {
		return
	}
	if l, err = m.LayerByName(layer); err != nil {
		return
	}
	output, err = l.GetGrid()
	return
}

func script(input, script string) (output string, err error) {
	var (
		byteOutput []byte
		f          fauxfile.File
		fs         = fauxfile.NewMockFilesystem()
		scripter   = NewTmxScripter(fs)
	)
	scripter.InputPath = "./map.tmx"
	scripter.OutputPath = "./modified.tmx"
	scripter.ScriptPath = "./script.js"
	if f, err = fs.Create(scripter.InputPath); err != nil {
		err = fmt.Errorf("Could not create input file: %v", err)
		return
	}
	f.Write([]byte(input))
	f.Close()
	if f, err = fs.Create(scripter.ScriptPath); err != nil {
		err = fmt.Errorf("Could not create script file: %v", err)
		return
	}
	f.Write([]byte(script))
	f.Close()
	if err = scripter.Run(); err != nil {
		err = fmt.Errorf("Error when running scripter: %s", err)
		return
	}
	if f, err = fs.Open(scripter.OutputPath); err != nil {
		err = fmt.Errorf("Could not open output file: %v", err)
		return
	}
	defer f.Close()
	if byteOutput, err = ioutil.ReadAll(f); err != nil {
		err = fmt.Errorf("Could not read output file: %v", err)
		return
	}
	output = string(byteOutput)
	return
}

const TEST_MAP = `
<?xml version="1.0" encoding="UTF-8"?>
<map version="1.0" orientation="orthogonal" width="3" height="3" tilewidth="32" tileheight="32">
 <tileset firstgid="1" name="sprites32" tilewidth="32" tileheight="32">
  <image source="sprites.png" width="512" height="512"/>
 </tileset>
 <layer name="layer1" width="3" height="3">
  <data>
   <tile gid="1" />
   <tile gid="0" />
   <tile gid="0" />

   <tile gid="0" />
   <tile gid="1" />
   <tile gid="0" />

   <tile gid="0" />
   <tile gid="0" />
   <tile gid="1" />
  </data>
 </layer>
</map>
`

func TestRun(t *testing.T) {
	var (
		result string
		grid   tmxgo.DataTileGrid
		err    error
	)
	if result, err = script(TEST_MAP, `
		console.log('foo');
	`); err != nil {
		t.Fatalf("Problem running script: %v", err)
	}
	if grid, err = getGrid(result, "layer1"); err != nil {
		t.Fatalf("Problem getting grid: %v", err)
	}
	fmt.Printf("Got grid: %v\n", grid)
}
