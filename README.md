# tmxscripter
A utility for running scripts against a [Tiled Map Editor](http://www.mapeditor.org/) map.

The purpose of tmxscripter is to allow for procedural modifications to a Tiled map as part of a development workflow.  Possible uses include:

  * Automatically creating collision layers based off of data in the map.
  * Automatically inserting boundary tiles according to a heuristic (imagine automatically adding a sand boundary around bodies of water, for example)
  * Procedurally generating maps which have valid paths from one end to the other.

## Installation

    go get -u github.com/kurrik/tmxscripter

## Running

    $GOPATH/bin/tmxscripter -input=foo -output=bar -script=baz

Where `foo` is the path to a TMX file, `bar` is the file to be written, and `baz` is a JavaScript file which will modify the contents of `foo`.

## Scripting

The script file must be written in JavaScript.  When run, it must set up any event listeners which will handle map processing.  This can be done using the `addEventListener` method.

Here is a simple script which adds `1` to the ID of any tile in the map:

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

### Events

Name | Arguments | Description
---- | --------- | -----------
map  | [Map](https://godoc.org/github.com/kurrik/tmxscripter/tmxscripter#ScriptableMap) | Called when the map is loaded and ready to be modified

### Javascript API

A few convenience methods have been added to the global scope.

 * `addEventListener(name, callback)` Calls the `callback` function when an event with name `name` is fired. Example:
        
        addEventListener("foo", function(data) {...});

 * `readFile(path)` Attempts to read the contents of the file at `path` synchronously and return it as a string.  `path` is relative to the script file's location.  Example:
        
        data = JSON.parse(readFile("relative/path/to/data.json"))

 * `writeFile(path, data)` Attempts to write `data` string into a file located at `path` synchronously.  `path` is relative to the script file's location.  Example:
        
        writeFile("relative/path/to/output.json", JSON.stringify(data))

[Underscore](http://underscorejs.org/) has been included automatically.

Most Tiled entities are wrapped for easier scripting.  Methods available on wrapped objects are available at the [class documentation](https://godoc.org/github.com/kurrik/tmxscripter/tmxscripter).  Following are some of the core classes and methods and their intended use.  Consult the godoc for the most up-to-date information.

A [Map](https://godoc.org/github.com/kurrik/tmxscripter/tmxscripter#ScriptableMap) represents an entire Tiled map.
 * `GetLayer(name)` Returns the layer matching the given name.
 * `AddLayer(name)` Creates a new layer with the given name, initialized to 0s.

A [Layer](https://godoc.org/github.com/kurrik/tmxscripter/tmxscripter#ScriptableLayer) represents a single layer.
  * `GetGrid()` Returns an object representing the 2D tile array.

A [Grid](https://godoc.org/github.com/kurrik/tmxscripter/tmxscripter#ScriptableGrid) represents a 2D tile array.
  * `Height()` Returns the height of the grid in tiles.
  * `Width()` Returns the width of the grid in tiles.
  * `TileAt(x, y)` Returns a tile at the given coordinates.
  * `TileList()` Returns a linear array of tiles, convenient for using underscore's map function.
  * `Save()` Must be called to persist changes to the grid back to the layer.

A [Tile](https://godoc.org/github.com/kurrik/tmxscripter/tmxscripter#ScriptableTile) represents a single tile entity
  * `Id` The global tileset index for this tile
  * `FlipX` Whether the tile is flipped horizontally
  * `FlipY` Whether the tile is flipped vertically
  * `FlipD` Whether the tile is flipped along the diagonal axis
