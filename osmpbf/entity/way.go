package entity

import (
	"github.com/twpayne/go-geom"
)

type Way struct {
	*Entity

	nodes []*Node
}

func NewWay(id int64) *Way {
	w := &Way{
		Entity: NewEntity(id),
	}
	w.SetGeometry(w.getGeometry)

	return w
}

func (w *Way) getGeometry() geom.T {
	var coords []geom.Coord
	for _, node := range w.nodes {
		coords = append(coords, geom.Coord{node.lon, node.lat})
	}

	return geom.NewLineString(geom.XY).MustSetCoords(coords)
}

func (w *Way) GetNodes() []*Node {
	return w.nodes
}

func (w *Way) SetNodes(nodes []*Node) {
	w.nodes = nodes
}
