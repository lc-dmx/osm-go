package entity

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestWay(t *testing.T) {
	node1 := NewNode(1)
	node1.SetLon(121.5136637)
	node1.SetLat(25.049556)

	node2 := NewNode(2)
	node2.SetLon(121.511518)
	node2.SetLat(25.0473407)

	node3 := NewNode(3)
	node3.SetLon(121.50012)
	node3.SetLat(25.03112)

	way := NewWay(11)
	way.SetNodes([]*Node{node1, node2, node3})

	fmt.Println(way.GetGeometry().FlatCoords())
	wkt, err := way.ToWKT()
	assert.Nil(t, err)
	assert.Equal(t, "LINESTRING (121.5136637 25.049556, 121.511518 25.0473407, 121.50012 25.03112)", wkt)
	j, err := way.ToJSON()
	assert.Nil(t, err)
	assert.Equal(t, "{\"type\":\"Feature\",\"id\":\"11\",\"geometry\":{\"type\":\"LineString\",\"coordinates\":[[121.5136637,25.049556],[121.511518,25.0473407],[121.50012,25.03112]]},\"properties\":{\"id\":\"11\"}}", j)
}
