package osmpbf

import (
	"testing"
	"time"

	"github.com/lc-dmx/osm-go/osmpbf/entity"
	"github.com/stretchr/testify/assert"
)

func TestEncodeEntity(t *testing.T) {
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
	way2.SetChangesetId(2)
	way2.SetTimestamp(time.Now().Unix())
	way2.SetUser("bb")
	way2.SetTags(tags)
	way2.SetNodes([]*entity.Node{node1, node3})

	relation := entity.NewRelation(6)
	relation.SetVisible(true)
	relation.SetVersion(8)
	relation.SetUid(10003)
	relation.SetChangesetId(4)
	relation.SetTimestamp(time.Now().Unix())
	relation.SetUser("dd")
	relation.SetTags(tags)
	relation.SetRelationMembers([]*entity.RelationMember{
		entity.NewRelationMember(node1, "node"),
		entity.NewRelationMember(way1, "way"),
	})

	enc := NewEncoder()
	enc.encodeEntity(node1)
	assert.Equal(t, 109, enc.estimateBlockSize())
	assert.Equal(t, node1.GetId(), enc.nodes.GetId()[0])
	assert.Equal(t, node1.GetVersion(), enc.nodes.GetDenseinfo().GetVersion()[0])
	assert.Equal(t, node1.GetChangesetId(), enc.nodes.GetDenseinfo().GetChangeset()[0])
	assert.Equal(t, node1.GetTimestamp(), enc.nodes.GetDenseinfo().GetTimestamp()[0])
	assert.Equal(t, node1.GetUid(), enc.nodes.GetDenseinfo().GetUid()[0])
	assert.Equal(t, int32(enc.sTable.getStrIndex(node1.GetUser())), enc.nodes.GetDenseinfo().GetUserSid()[0])
	assert.Equal(t, node1.GetVisible(), enc.nodes.GetDenseinfo().GetVisible()[0])

	enc.encodeEntity(node2)
	assert.Equal(t, 192, enc.estimateBlockSize())
	assert.Equal(t, node2.GetId()-node1.GetId(), enc.nodes.GetId()[1])
	assert.Equal(t, node2.GetVersion(), enc.nodes.GetDenseinfo().GetVersion()[1])
	assert.Equal(t, node2.GetChangesetId()-node1.GetChangesetId(), enc.nodes.GetDenseinfo().GetChangeset()[1])
	assert.Equal(t, node2.GetTimestamp()-node1.GetTimestamp(), enc.nodes.GetDenseinfo().GetTimestamp()[1])
	assert.Equal(t, node2.GetUid()-node1.GetUid(), enc.nodes.GetDenseinfo().GetUid()[1])
	assert.Equal(t, int32(enc.sTable.getStrIndex(node2.GetUser())-enc.sTable.getStrIndex(node1.GetUser())), enc.nodes.GetDenseinfo().GetUserSid()[1])
	assert.Equal(t, node2.GetVisible(), enc.nodes.GetDenseinfo().GetVisible()[1])

	enc.encodeEntity(node3)
	assert.Equal(t, 275, enc.estimateBlockSize())
	assert.Equal(t, node3.GetId()-node2.GetId(), enc.nodes.GetId()[2])
	assert.Equal(t, node3.GetVersion(), enc.nodes.GetDenseinfo().GetVersion()[2])
	assert.Equal(t, node3.GetChangesetId()-node2.GetChangesetId(), enc.nodes.GetDenseinfo().GetChangeset()[2])
	assert.Equal(t, node3.GetTimestamp()-node2.GetTimestamp(), enc.nodes.GetDenseinfo().GetTimestamp()[2])
	assert.Equal(t, node3.GetUid()-node2.GetUid(), enc.nodes.GetDenseinfo().GetUid()[2])
	assert.Equal(t, int32(enc.sTable.getStrIndex(node3.GetUser())-enc.sTable.getStrIndex(node2.GetUser())), enc.nodes.GetDenseinfo().GetUserSid()[2])
	assert.Equal(t, node3.GetVisible(), enc.nodes.GetDenseinfo().GetVisible()[2])

	enc.encodeEntity(way1)
	assert.Equal(t, 352, enc.estimateBlockSize())
	assert.Equal(t, way1.GetId(), enc.ways[0].GetId())
	assert.Equal(t, way1.GetVersion(), enc.ways[0].GetInfo().GetVersion())
	assert.Equal(t, way1.GetChangesetId(), enc.ways[0].GetInfo().GetChangeset())
	assert.Equal(t, way1.GetTimestamp(), enc.ways[0].GetInfo().GetTimestamp())
	assert.Equal(t, way1.GetUid(), enc.ways[0].GetInfo().GetUid())
	assert.Equal(t, uint32(enc.sTable.getStrIndex(way1.GetUser())), enc.ways[0].GetInfo().GetUserSid())
	assert.Equal(t, way1.GetVisible(), enc.ways[0].GetInfo().GetVisible())
	assert.Equal(t, 2, len(enc.ways[0].GetRefs()))
	assert.Equal(t, node1.GetId(), enc.ways[0].GetRefs()[0])
	assert.Equal(t, node2.GetId()-node1.GetId(), enc.ways[0].GetRefs()[1])

	enc.encodeEntity(way2)
	assert.Equal(t, 429, enc.estimateBlockSize())
	assert.Equal(t, way2.GetId(), enc.ways[1].GetId())
	assert.Equal(t, way2.GetVersion(), enc.ways[1].GetInfo().GetVersion())
	assert.Equal(t, way2.GetChangesetId(), enc.ways[1].GetInfo().GetChangeset())
	assert.Equal(t, way2.GetTimestamp(), enc.ways[1].GetInfo().GetTimestamp())
	assert.Equal(t, way2.GetUid(), enc.ways[1].GetInfo().GetUid())
	assert.Equal(t, uint32(enc.sTable.getStrIndex(way2.GetUser())), enc.ways[1].GetInfo().GetUserSid())
	assert.Equal(t, way2.GetVisible(), enc.ways[1].GetInfo().GetVisible())
	assert.Equal(t, 2, len(enc.ways[1].GetRefs()))
	assert.Equal(t, node1.GetId(), enc.ways[1].GetRefs()[0])
	assert.Equal(t, node3.GetId()-node1.GetId(), enc.ways[1].GetRefs()[1])

	enc.encodeEntity(relation)
	assert.Equal(t, 531, enc.estimateBlockSize())
	assert.Equal(t, relation.GetId(), enc.relations[0].GetId())
	assert.Equal(t, relation.GetVersion(), enc.relations[0].GetInfo().GetVersion())
	assert.Equal(t, relation.GetChangesetId(), enc.relations[0].GetInfo().GetChangeset())
	assert.Equal(t, relation.GetTimestamp(), enc.relations[0].GetInfo().GetTimestamp())
	assert.Equal(t, relation.GetUid(), enc.relations[0].GetInfo().GetUid())
	assert.Equal(t, uint32(enc.sTable.getStrIndex(relation.GetUser())), enc.relations[0].GetInfo().GetUserSid())
	assert.Equal(t, relation.GetVisible(), enc.relations[0].GetInfo().GetVisible())
	assert.Equal(t, 2, len(enc.relations[0].GetMemids()))
	assert.Equal(t, node1.GetId(), enc.relations[0].GetMemids()[0])
	assert.Equal(t, way1.GetId()-node1.GetId(), enc.relations[0].GetMemids()[1])
}
