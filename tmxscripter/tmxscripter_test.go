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
		return
	}
	f.Write([]byte(input))
	f.Close()
	if f, err = fs.Create(scripter.ScriptPath); err != nil {
		return
	}
	f.Write([]byte(script))
	f.Close()
	if err = scripter.Run(); err != nil {
		return
	}
	if f, err = fs.Open(scripter.OutputPath); err != nil {
		return
	}
	defer f.Close()
	if byteOutput, err = ioutil.ReadAll(f); err != nil {
		return
	}
	output = string(byteOutput)
	return
}

func runTest(js string, layer string, expected []uint32) (err error) {
	var (
		result string
		grid   tmxgo.DataTileGrid
		ids    []uint32
	)
	if result, err = script(TEST_MAP, js); err != nil {
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
