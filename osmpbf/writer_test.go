package osmpbf

import (
	"bytes"
	"context"
	"testing"
	"time"

	"github.com/lc-dmx/osm-go/osmpbf/entity"
	"github.com/paulmach/osm"
	"github.com/paulmach/osm/osmpbf"
	"github.com/stretchr/testify/assert"
)

func TestWriteHeader(t *testing.T) {
	ctx := context.Background()
	var in bytes.Buffer
	w, err := NewWriter(ctx, &in)
	assert.Nil(t, err)

	err = w.Close()
	assert.Nil(t, err)
	assert.Equal(t, 115, in.Len())
}

func TestWriteEntity(t *testing.T) {
	tags := map[string]interface{}{
		"first":  "hehe",
		"second": "haha",
		"third":  "oo",
	}

	node1 := entity.NewNode(1)
	node1.SetVisible(false)
	node1.SetVersion(2)
	node1.SetUid(10000)
	node1.SetChangesetId(1)
	node1.SetTimestamp(time.Now().Unix())
	node1.SetUser("aa")
	node1.SetTags(tags)
	node1.SetLon(116.385822)
	node1.SetLat(39.9066105)

	node2 := entity.NewNode(2)
	node2.SetVisible(true)
	node2.SetVersion(4)
	node2.SetUid(10001)
	node2.SetChangesetId(2)
	node2.SetTimestamp(time.Now().Unix())
	node2.SetUser("bb")
	node2.SetTags(tags)
	node2.SetLon(39.9066105)
	node2.SetLat(39.9066105)

	node3 := entity.NewNode(3)
	node3.SetVisible(false)
	node3.SetVersion(5)
	node3.SetUid(10002)
	node3.SetChangesetId(3)
	node3.SetTimestamp(time.Now().Unix())
	node3.SetUser("cc")
	node3.SetTags(tags)
	node3.SetLon(39.9066105)
	node3.SetLat(39.9066105)

	way1 := entity.NewWay(4)
	way1.SetVisible(true)
	way1.SetVersion(3)
	way1.SetUid(10000)
	way1.SetChangesetId(1)
	way1.SetTimestamp(time.Now().Unix())
	way1.SetUser("aa")
	way1.SetTags(tags)
	way1.SetNodes([]*entity.Node{node1, node2})

	way2 := entity.NewWay(5)
	way2.SetVisible(false)
	way2.SetVersion(6)
	way2.SetUid(10001)
	way2.SetChangesetId(1)
	way2.SetTimestamp(time.Now().Unix())
	way2.SetUser("bb")
	way2.SetTags(tags)
	way2.SetNodes([]*entity.Node{node1, node3})

	relation := entity.NewRelation(6)
	relation.SetVisible(true)
	relation.SetVersion(8)
	relation.SetUid(10003)
	relation.SetChangesetId(1)
	relation.SetTimestamp(time.Now().Unix())
	relation.SetUser("dd")
	relation.SetTags(tags)
	relation.SetRelationMembers([]*entity.RelationMember{
		entity.NewRelationMember(node1, "node"),
		entity.NewRelationMember(way1, "way"),
	})

	ctx := context.Background()
	var in bytes.Buffer
	w, err := NewWriter(ctx, &in)
	assert.Nil(t, err)

	err = w.WriteEntity(node1)
	assert.Nil(t, err)
	err = w.WriteEntity(node2)
	assert.Nil(t, err)
	err = w.WriteEntity(node3)
	assert.Nil(t, err)
	err = w.WriteEntity(way1)
	assert.Nil(t, err)
	err = w.WriteEntity(way2)
	assert.Nil(t, err)
	err = w.WriteEntity(relation)
	assert.Nil(t, err)

	err = w.Close()
	assert.Nil(t, err)

	osmNodes := make([]*osm.Node, 0, 3)
	osmWays := make([]*osm.Way, 0, 2)
	osmRelations := make([]*osm.Relation, 0, 1)

	scanner := osmpbf.New(ctx, &in, 3)
	defer scanner.Close()
	for scanner.Scan() {
		switch scanner.Object().(type) {
		case *osm.Node:
			osmNodes = append(osmNodes, scanner.Object().(*osm.Node))
		case *osm.Way:
			osmWays = append(osmWays, scanner.Object().(*osm.Way))
		case *osm.Relation:
			osmRelations = append(osmRelations, scanner.Object().(*osm.Relation))
		}
	}
	err = scanner.Err()
	assert.Nil(t, err)

	assert.Equal(t, node1.GetVisible(), osmNodes[0].Visible)
	assert.Equal(t, osm.NodeID(node1.GetId()), osmNodes[0].ID)
	assert.Equal(t, int(node1.GetVersion()), osmNodes[0].Version)
	assert.Equal(t, osm.ChangesetID(node1.GetChangesetId()), osmNodes[0].ChangesetID)
	assert.Equal(t, node1.GetTimestamp(), osmNodes[0].Timestamp.Unix())
	assert.Equal(t, osm.UserID(node1.GetUid()), osmNodes[0].UserID)
	assert.Equal(t, node1.GetUser(), osmNodes[0].User)
	assert.Equal(t, len(node1.GetTags()), len(osmNodes[0].Tags))

	assert.Equal(t, node2.GetVisible(), osmNodes[1].Visible)
	assert.Equal(t, osm.NodeID(node2.GetId()), osmNodes[1].ID)
	assert.Equal(t, int(node2.GetVersion()), osmNodes[1].Version)
	assert.Equal(t, osm.ChangesetID(node2.GetChangesetId()), osmNodes[1].ChangesetID)
	assert.Equal(t, node2.GetTimestamp(), osmNodes[1].Timestamp.Unix())
	assert.Equal(t, osm.UserID(node2.GetUid()), osmNodes[1].UserID)
	assert.Equal(t, node2.GetUser(), osmNodes[1].User)
	assert.Equal(t, len(node2.GetTags()), len(osmNodes[1].Tags))

	assert.Equal(t, node3.GetVisible(), osmNodes[2].Visible)
	assert.Equal(t, osm.NodeID(node3.GetId()), osmNodes[2].ID)
	assert.Equal(t, int(node3.GetVersion()), osmNodes[2].Version)
	assert.Equal(t, osm.ChangesetID(node3.GetChangesetId()), osmNodes[2].ChangesetID)
	assert.Equal(t, node3.GetTimestamp(), osmNodes[2].Timestamp.Unix())
	assert.Equal(t, osm.UserID(node3.GetUid()), osmNodes[2].UserID)
	assert.Equal(t, node3.GetUser(), osmNodes[2].User)
	assert.Equal(t, len(node3.GetTags()), len(osmNodes[2].Tags))

	assert.Equal(t, way1.GetVisible(), osmWays[0].Visible)
	assert.Equal(t, osm.WayID(way1.GetId()), osmWays[0].ID)
	assert.Equal(t, int(way1.GetVersion()), osmWays[0].Version)
	assert.Equal(t, osm.ChangesetID(way1.GetChangesetId()), osmWays[0].ChangesetID)
	assert.Equal(t, way1.GetTimestamp(), osmWays[0].Timestamp.Unix())
	assert.Equal(t, osm.UserID(way1.GetUid()), osmWays[0].UserID)
	assert.Equal(t, way1.GetUser(), osmWays[0].User)
	assert.Equal(t, len(way1.GetTags()), len(osmWays[0].Tags))
	assert.Equal(t, len(way1.GetNodes()), len(osmWays[0].Nodes))

	assert.Equal(t, way2.GetVisible(), osmWays[1].Visible)
	assert.Equal(t, osm.WayID(way2.GetId()), osmWays[1].ID)
	assert.Equal(t, int(way2.GetVersion()), osmWays[1].Version)
	assert.Equal(t, osm.ChangesetID(way2.GetChangesetId()), osmWays[1].ChangesetID)
	assert.Equal(t, way2.GetTimestamp(), osmWays[1].Timestamp.Unix())
	assert.Equal(t, osm.UserID(way2.GetUid()), osmWays[1].UserID)
	assert.Equal(t, way2.GetUser(), osmWays[1].User)
	assert.Equal(t, len(way2.GetTags()), len(osmWays[1].Tags))
	assert.Equal(t, len(way2.GetNodes()), len(osmWays[1].Nodes))

	assert.Equal(t, relation.GetVisible(), osmRelations[0].Visible)
	assert.Equal(t, osm.RelationID(relation.GetId()), osmRelations[0].ID)
	assert.Equal(t, int(relation.GetVersion()), osmRelations[0].Version)
	assert.Equal(t, osm.ChangesetID(relation.GetChangesetId()), osmRelations[0].ChangesetID)
	assert.Equal(t, relation.GetTimestamp(), osmRelations[0].Timestamp.Unix())
	assert.Equal(t, osm.UserID(relation.GetUid()), osmRelations[0].UserID)
	assert.Equal(t, relation.GetUser(), osmRelations[0].User)
	assert.Equal(t, len(relation.GetTags()), len(osmRelations[0].Tags))
	assert.Equal(t, len(relation.GetRelationMembers()), len(osmRelations[0].Members))
}
