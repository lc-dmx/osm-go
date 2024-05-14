package entity

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRelation(t *testing.T) {
	node1 := NewNode(1)
	node1.SetLon(101.51366)
	node1.SetLat(25.049556)

	node2 := NewNode(2)
	node2.SetLon(121.511518)
	node2.SetLat(15.0473407)

	node3 := NewNode(3)
	node3.SetLon(110.511)
	node3.SetLat(35.0473407)

	way1 := NewWay(11)
	way1.SetNodes([]*Node{node1, node2})

	way2 := NewWay(12)
	way2.SetNodes([]*Node{node2, node3})

	relation1 := NewRelation(101)
	relation1.SetRelationMembers([]*RelationMember{
		NewRelationMember(way1, "way"),
		NewRelationMember(node1, "node"),
	})

	relation2 := NewRelation(102)
	relation2.SetRelationMembers([]*RelationMember{
		NewRelationMember(relation1, "relation"),
		NewRelationMember(way1, "way"),
	})

	relation3 := NewRelation(103)
	relation3.SetRelationMembers([]*RelationMember{
		NewRelationMember(relation1, "relation"),
		NewRelationMember(relation2, "relation"),
		NewRelationMember(node2, "node"),
	})

	relation4 := NewRelation(104)
	relation4.SetRelationMembers([]*RelationMember{
		NewRelationMember(relation1, "relation"),
		NewRelationMember(relation2, "relation"),
		NewRelationMember(relation3, "relation"),
		NewRelationMember(node3, "node"),
	})

	wkt, err := relation4.ToWKT()
	assert.Nil(t, err)
	fmt.Println(wkt)
	j, err := relation4.ToJSON()
	assert.Nil(t, err)
	fmt.Println(j)
}
