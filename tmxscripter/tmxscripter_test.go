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

package tmxscripter

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

func getGridIds(grid tmxgo.DataTileGrid) (ids []uint32) {
	ids = make([]uint32, grid.Width*grid.Height)
	for y := 0; y < grid.Height; y++ {
		for x := 0; x < grid.Width; x++ {
			ids[y*grid.Width+x] = grid.Tiles[x][y].Id
		}
	}
	return
}

func writeFile(fs fauxfile.Filesystem, path, contents string) (err error) {
	var f fauxfile.File
	if f, err = fs.Create(path); err != nil {
		return
	}
	defer f.Close()
	f.Write([]byte(contents))
	return
}

func readFile(fs fauxfile.Filesystem, path string) (contents string, err error) {
	var (
		byteOutput []byte
		f          fauxfile.File
	)
	if f, err = fs.Open(path); err != nil {
		return
	}
	defer f.Close()
	if byteOutput, err = ioutil.ReadAll(f); err != nil {
		return
	}
	contents = string(byteOutput)
	return
}

func script(input, script, data string) (output string, err error) {
	var (
		fs       = fauxfile.NewMockFilesystem()
		scripter = NewTmxScripter(fs)
	)
	scripter.InputPath = "./foo/map.tmx"
	scripter.OutputPath = "./foo/modified.tmx"
	scripter.ScriptPath = "./foo/script.js"
	if err = fs.MkdirAll("./foo/bar", 0755); err != nil {
		return
	}
	if err = writeFile(fs, scripter.InputPath, input); err != nil {
		return
	}
	if err = writeFile(fs, scripter.ScriptPath, script); err != nil {
		return
	}
	if err = writeFile(fs, "./foo/bar/data.json", data); err != nil {
		return
	}
	if err = scripter.Run(); err != nil {
		return
	}
	if output, err = readFile(fs, scripter.OutputPath); err != nil {
		return
	}
	return
}

func runTest(js, layer string, expected []uint32) (err error) {
	var (
		result string
		grid   tmxgo.DataTileGrid
		ids    []uint32
	)
	if result, err = script(TEST_MAP, js, TEST_DATA); err != nil {
		return
	}
	if grid, err = getGrid(result, layer); err != nil {
		return
	}
	ids = getGridIds(grid)
	if len(ids) != len(expected) {
		err = fmt.Errorf("IDs didn't match! expected: %v, got %v", expected, ids)
		return
	}
	for i := 0; i < len(ids); i++ {
		if ids[i] != expected[i] {
			err = fmt.Errorf("No match at index %v! expected: %v, got %v", i, expected[i], ids[i])
			return
		}
	}
	return
}

const TEST_DATA = `
{
  "aString": "bar",
  "aNumber": 20
}`

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

func TestValidateInputPath(t *testing.T) {
	var (
		fs       = fauxfile.NewMockFilesystem()
		scripter = NewTmxScripter(fs)
		err      error
	)
	scripter.InputPath = "./doesnotexist.tmx"
	scripter.OutputPath = "./modified.tmx"
	scripter.ScriptPath = "./script.js"
	if err = writeFile(fs, scripter.ScriptPath, `console.log("foo");`); err != nil {
		return
	}
	if err = scripter.Validate(); err == nil {
		t.Fatalf("Expected error if input path doesn't exist")
		return
	}
}

func TestValidateScriptPath(t *testing.T) {
	var (
		fs       = fauxfile.NewMockFilesystem()
		scripter = NewTmxScripter(fs)
		err      error
	)
	scripter.InputPath = "./map.tmx"
	scripter.OutputPath = "./modified.tmx"
	scripter.ScriptPath = "./doesnotexist.js"
	if err = writeFile(fs, scripter.InputPath, TEST_MAP); err != nil {
		return
	}
	if err = scripter.Validate(); err == nil {
		t.Fatalf("Expected error if script path doesn't exist")
		return
	}
}

func TestNoop(t *testing.T) {
	if err := runTest(`
		// This script does nothing.
	`, "layer1", []uint32{
		1, 0, 0,
		0, 1, 0,
		0, 0, 1,
	}); err != nil {
		t.Fatal(err)
	}
}

func TestPlusOne(t *testing.T) {
	if err := runTest(`
		// This script adds one to each tile Id.
		addEventListener("map", function(m) {
			var layer = m.GetLayer("layer1"),
			    grid = layer.GetGrid(),
			    tile;
			for (var y = 0; y < grid.Height(); y++) {
				for (var x = 0; x < grid.Width(); x++) {
					tile = grid.TileAt(x, y);
					tile.Id += 1;
				}
			}
			grid.Save();
		});
	`, "layer1", []uint32{
		2, 1, 1,
		1, 2, 1,
		1, 1, 2,
	}); err != nil {
		t.Fatal(err)
	}
}

func TestAddLayer(t *testing.T) {
	if err := runTest(`
		// This script adds a new layer to the map.
		addEventListener("map", function(m) {
			var layer = m.AddLayer("layer2"),
			    grid = layer.GetGrid(),
			    tile;
			grid.TileAt(1,1).Id = 1;
			grid.Save();
		});
	`, "layer2", []uint32{
		0, 0, 0,
		0, 1, 0,
		0, 0, 0,
	}); err != nil {
		t.Fatal(err)
	}
}

func TestReadData(t *testing.T) {
	if err := runTest(`
		// This script reads a data file and uses it to adjust a layer.
		addEventListener("map", function(m) {
			var layer = m.GetLayer("layer1"),
			    grid = layer.GetGrid(),
			    data = JSON.parse(readFile("bar/data.json")),
			    tile;
			for (var y = 0; y < grid.Height(); y++) {
				for (var x = 0; x < grid.Width(); x++) {
					grid.TileAt(x, y).Id = data.aNumber;
				}
			}
			grid.Save();
		});
	`, "layer1", []uint32{
		20, 20, 20,
		20, 20, 20,
		20, 20, 20,
	}); err != nil {
		t.Fatal(err)
	}
}
