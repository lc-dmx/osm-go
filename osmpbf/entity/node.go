package entity

import (
	"github.com/twpayne/go-geom"
)

type Node struct {
	*Entity

	lon float64
	lat float64
}

func NewNode(id int64) *Node {
	n := &Node{
		Entity: NewEntity(id),
	}
	n.SetGeometry(n.getGeometry)

	return n
}

func (n *Node) getGeometry() geom.T {
	return geom.NewPoint(geom.XY).MustSetCoords(geom.Coord{n.lon, n.lat})
}

func (n *Node) GetLon() float64 {
	return n.lon
}

func (n *Node) GetLat() float64 {
	return n.lat
}

func (n *Node) SetLon(lon float64) {
	n.lon = lon
}

func (n *Node) SetLat(lat float64) {
	n.lat = lat
}
