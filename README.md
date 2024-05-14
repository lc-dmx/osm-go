# osm-go

This package is a general purpose library for reading, writing and working
with [OpenStreetMap](https://osm.org) data in Go (golang). Currently, it has the ability to:

- write [OSM PBF](https://wiki.openstreetmap.org/wiki/PBF_Format) data to file.

Made available by the package are the following types:

- Node
- Way
- Relation

## List of sub-package utilities

- [`osmpbf`](osmpbf) - stream processing of `*.osm.pbf` files
