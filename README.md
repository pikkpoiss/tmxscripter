# tmxscripter
A utility for running scripts against a [Tiled Map Editor](http://www.mapeditor.org/) map.

The purpose of tmxscripter is to allow for procedural modifications to a Tiled map as part of a development workflow.  Possible uses include:

  * Automatically creating collision layers based off of data in the map.
  * Automatically inserting boundary tiles according to a heuristic (imagine automatically adding a sand boundary around bodies of water, for example)
  * Procedurally generating maps which have valid paths from one end to the other.

## Installation

    go get -u github.com/kurrik/autoslice

## Running

    $GOPATH/bin/tmxscripter -input=foo -output=bar -script=baz

Where `foo` is the path to a TMX file, `bar` is the file to be written, and `baz` is a JavaScript file which will modify the contents of `foo`.

## JavaScript API

The script file must be written in JavaScript.  When run, it must set up any event listeners which will handle map processing.  This can be done using the `addEventListener` method.

Here is a simple script which adds 1 to the ID of any tile in the map:

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

| Name        | Arguments           | Description  |
| ------------- |:-------------:| -----:|
| col 3 is      | right-aligned | $1600 |
| col 2 is      | centered      |   $12 |
| zebra stripes | are neat      |    $1 |
