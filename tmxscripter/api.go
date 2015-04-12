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
	"github.com/kurrik/tmxgo"
)

type ScriptableMap struct {
	*tmxgo.Map
}

func NewScriptableMap(m *tmxgo.Map) *ScriptableMap {
	return &ScriptableMap{
		Map: m,
	}
}

// Returns a layer with the given name if one exists.
func (m *ScriptableMap) GetLayer(name string) *ScriptableLayer {
	if l, err := m.LayerByName(name); err != nil {
		return nil
	} else {
		return NewScriptableLayer(l)
	}
}

// Adds a new layer with the given name to the map. Tile IDs will be 0.
func (m *ScriptableMap) AddLayer(name string) *ScriptableLayer {
	var grid = tmxgo.DataTileGrid{
		Width:  int(m.Width),
		Height: int(m.Height),
		Tiles:  make([][]tmxgo.DataTileGridTile, m.Width),
	}
	for x := 0; x < int(m.Width); x++ {
		grid.Tiles[x] = make([]tmxgo.DataTileGridTile, m.Height)
		for y := 0; y < int(m.Height); y++ {
			grid.Tiles[x][y] = tmxgo.DataTileGridTile{
				Id:    0,
				FlipX: false,
				FlipY: false,
				FlipD: false,
			}
		}
	}
	var layer = &tmxgo.Layer{
		Name:   name,
		Width:  m.Width,
		Height: m.Height,
		Data:   &tmxgo.Data{},
	}
	layer.SetGrid(grid)
	m.Layers = append(m.Layers, layer)
	return NewScriptableLayer(layer)
}

type ScriptableLayer struct {
	*tmxgo.Layer
}

func NewScriptableLayer(l *tmxgo.Layer) *ScriptableLayer {
	return &ScriptableLayer{
		Layer: l,
	}
}

// Returns a scriptable grid for this layer.
func (l *ScriptableLayer) GetGrid() *ScriptableGrid {
	return NewScriptableGrid(l.Layer)
}

type ScriptableGrid struct {
	*tmxgo.DataTileGrid
	*tmxgo.Layer
}

func NewScriptableGrid(l *tmxgo.Layer) *ScriptableGrid {
	if g, err := l.GetGrid(); err != nil {
		return nil
	} else {
		return &ScriptableGrid{
			DataTileGrid: &g,
			Layer:        l,
		}
	}
}

// Returns the width of the grid in tiles.
func (g *ScriptableGrid) Width() int {
	return g.DataTileGrid.Width
}

// Returns the height of the grid in tiles.
func (g *ScriptableGrid) Height() int {
	return g.DataTileGrid.Height
}

// Saves the grid back into the layer.
func (g *ScriptableGrid) Save() {
	g.Layer.SetGrid(*g.DataTileGrid)
}

// Returns the tile at the specified location.
func (g *ScriptableGrid) TileAt(x int, y int) *ScriptableTile {
	return NewScriptableTile(&g.Tiles[x][y])
}

// Returns a linear array of tiles.
func (g *ScriptableGrid) TileList() []*ScriptableTile {
	var tiles = make([]*ScriptableTile, g.DataTileGrid.Width*g.DataTileGrid.Height)
	for y := 0; y < g.DataTileGrid.Height; y++ {
		for x := 0; x < g.DataTileGrid.Width; x++ {
			tiles[y*g.DataTileGrid.Width+x] = NewScriptableTile(&g.Tiles[x][y])
		}
	}
	return tiles
}

// Represents a tile object.  Has Id, FlipX, FlipY and FlipD attributes.
type ScriptableTile struct {
	*tmxgo.DataTileGridTile
}

func NewScriptableTile(t *tmxgo.DataTileGridTile) *ScriptableTile {
	return &ScriptableTile{
		DataTileGridTile: t,
	}
}
