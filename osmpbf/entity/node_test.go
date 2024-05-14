package entity

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNode(t *testing.T) {
	node1 := NewNode(1)
	node1.SetLon(121.5136637)
	node1.SetLat(25.049556)
	fmt.Println(node1.GetGeometry().FlatCoords())
	wkt, err := node1.ToWKT()
	assert.Nil(t, err)
	assert.Equal(t, "POINT (121.5136637 25.049556)", wkt)
	j, err := node1.ToJSON()
	assert.Nil(t, err)
	assert.Equal(t, "{\"type\":\"Feature\",\"id\":\"1\",\"geometry\":{\"type\":\"Point\",\"coordinates\":[121.5136637,25.049556]},\"properties\":{\"id\":\"1\"}}", j)

	node2 := NewNode(2)
	node2.SetLon(121.511518)
	node2.SetLat(25.0473407)
	fmt.Println(node2.GetGeometry().FlatCoords())
	wkt, err = node2.ToWKT()
	assert.Nil(t, err)
	assert.Equal(t, "POINT (121.511518 25.0473407)", wkt)
	j, err = node2.ToJSON()
	assert.Nil(t, err)
	assert.Equal(t, "{\"type\":\"Feature\",\"id\":\"2\",\"geometry\":{\"type\":\"Point\",\"coordinates\":[121.511518,25.0473407]},\"properties\":{\"id\":\"2\"}}", j)
}
