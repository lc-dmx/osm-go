package osmpbf

import (
	"math"
	"sort"

	"github.com/apache/thrift/lib/go/thrift"
	"github.com/lc-dmx/osm-go/osmpbf/entity"
	"github.com/lc-dmx/osm-go/osmpbf/model_pb"
)

const (
	FEATURE_OSM_SCHEMA             = "OsmSchema-V0.6"
	FEATURE_DENSE_NODES            = "DenseNodes"
	FEATURE_HISTORICAL_INFORMATION = "HistoricalInformation"
	FEATURE_HAS_METADATA           = "Has_Metadata"
	OSM_PBF_WRITER                 = "osm_pbf_writer"
)

const (
	GRANULARITY = 100
	NANO        = 1e9
)

type Encoder struct {
	sTable    *StrTable
	nodes     *model_pb.DenseNodes
	ways      []*model_pb.Way
	relations []*model_pb.Relation

	lastNodeId          int64
	lastNodeLon         int64
	lastNodeLat         int64
	lastNodeTimestamp   int64
	lastNodeChangesetId int64
	lastNodeUid         int32
	lastNodeUserSid     int32

	// used for estimating
	nodeTagsLen        int
	wayNodesLen        int
	wayTagsLen         int
	relationMembersLen int
	relationTagsLen    int
}

func NewEncoder() *Encoder {
	return &Encoder{
		sTable: NewStringTable(),
		nodes: &model_pb.DenseNodes{
			Denseinfo: &model_pb.DenseInfo{},
		},
		ways:      make([]*model_pb.Way, 0),
		relations: make([]*model_pb.Relation, 0),
	}
}

func (enc *Encoder) encodeHeader() *model_pb.HeaderBlock {
	return &model_pb.HeaderBlock{
		RequiredFeatures: []string{FEATURE_OSM_SCHEMA, FEATURE_DENSE_NODES, FEATURE_HISTORICAL_INFORMATION},
		OptionalFeatures: []string{FEATURE_HAS_METADATA},
		Writingprogram:   thrift.StringPtr(OSM_PBF_WRITER),
	}
}

func (enc *Encoder) encodeData() *model_pb.PrimitiveBlock {
	if enc.estimateBlockSize() == 0 {
		return nil
	}

	block := &model_pb.PrimitiveBlock{
		Stringtable: &model_pb.StringTable{
			S: enc.sTable.getStrTable(),
		},
	}

	if enc.estimateNodeSize() > 0 {
		block.Primitivegroup = append(block.Primitivegroup, &model_pb.PrimitiveGroup{
			Dense: enc.nodes,
		})
		block.Granularity = thrift.Int32Ptr(GRANULARITY)
		block.LonOffset = thrift.Int64Ptr(0)
		block.LatOffset = thrift.Int64Ptr(0)
	}

	if enc.estimateWaySize() > 0 {
		block.Primitivegroup = append(block.Primitivegroup, &model_pb.PrimitiveGroup{
			Ways: enc.ways,
		})
	}

	if enc.estimateRelationSize() > 0 {
		block.Primitivegroup = append(block.Primitivegroup, &model_pb.PrimitiveGroup{
			Relations: enc.relations,
		})
	}

	return block
}

// remember to modify the estimation logic when there is an addition or deletion of a field.
func (enc *Encoder) estimateBlockSize() int {
	return enc.estimateNodeSize() + enc.estimateWaySize() + enc.estimateRelationSize() + enc.sTable.getStrTableSize()
}

func (enc *Encoder) estimateNodeSize() int {
	return len(enc.nodes.Id)*53 + enc.nodeTagsLen
}

func (enc *Encoder) estimateWaySize() int {
	return len(enc.ways)*37 + enc.wayNodesLen + enc.wayTagsLen
}

func (enc *Encoder) estimateRelationSize() int {
	return len(enc.relations)*37 + enc.relationMembersLen + enc.relationTagsLen
}

func (enc *Encoder) encodeEntity(exp entity.Exporter) {
	switch exp.(type) {
	case *entity.Node:
		enc.encodeNode(exp.(*entity.Node))
	case *entity.Way:
		enc.encodeWay(exp.(*entity.Way))
	case *entity.Relation:
		enc.encodeRelation(exp.(*entity.Relation))
	default:
		panic("no entity type matched")
	}
}

func (enc *Encoder) encodeNode(node *entity.Node) {
	enc.nodes.Id = append(enc.nodes.Id, node.GetId()-enc.lastNodeId)
	enc.lastNodeId = node.GetId()

	latMillis := doubleToNanoScaled(node.GetLat() / GRANULARITY)
	lonMillis := doubleToNanoScaled(node.GetLon() / GRANULARITY)
	enc.nodes.Lat = append(enc.nodes.Lat, latMillis-enc.lastNodeLat)
	enc.nodes.Lon = append(enc.nodes.Lon, lonMillis-enc.lastNodeLon)
	enc.lastNodeLat = latMillis
	enc.lastNodeLon = lonMillis

	enc.nodes.Denseinfo.Version = append(enc.nodes.Denseinfo.Version, node.GetVersion())
	enc.nodes.Denseinfo.Timestamp = append(enc.nodes.Denseinfo.Timestamp, node.GetTimestamp()-enc.lastNodeTimestamp)
	enc.lastNodeTimestamp = node.GetTimestamp()
	enc.nodes.Denseinfo.Changeset = append(enc.nodes.Denseinfo.Changeset, node.GetChangesetId()-enc.lastNodeChangesetId)
	enc.lastNodeChangesetId = node.GetChangesetId()
	enc.nodes.Denseinfo.Uid = append(enc.nodes.Denseinfo.Uid, node.GetUid()-enc.lastNodeUid)
	enc.lastNodeUid = node.GetUid()
	nodeUserSid := int32(enc.sTable.getStrIndex(node.GetUser()))
	enc.nodes.Denseinfo.UserSid = append(enc.nodes.Denseinfo.UserSid, nodeUserSid-enc.lastNodeUserSid)
	enc.lastNodeUserSid = nodeUserSid
	enc.nodes.Denseinfo.Visible = append(enc.nodes.Denseinfo.Visible, node.GetVisible())

	enc.nodeTagsLen += len(node.GetTags())*8 + 4
	sortedTags := make([]string, 0, len(node.GetTags()))
	for k := range node.GetTags() {
		sortedTags = append(sortedTags, k)
	}
	sort.Strings(sortedTags)
	for _, k := range sortedTags {
		enc.nodes.KeysVals = append(enc.nodes.KeysVals,
			int32(enc.sTable.getStrIndex(k)),
			int32(enc.sTable.getStrIndex(node.GetTags()[k].(string))))
	}
	// Index zero means 'end of tags for node'
	// The pattern is pattern is: ((<keyid> <valid>)* '0' )*
	// As an exception, if no node in the current block has any key/value pairs, this array does not contain any delimiters, but is simply empty.
	enc.nodes.KeysVals = append(enc.nodes.KeysVals, 0)
}

func doubleToNanoScaled(value float64) int64 {
	return int64(math.Round(value * NANO))
}

func (enc *Encoder) encodeWay(way *entity.Way) {
	w := &model_pb.Way{
		Id: thrift.Int64Ptr(way.GetId()),
		Info: &model_pb.Info{
			Version:   thrift.Int32Ptr(way.GetVersion()),
			Timestamp: thrift.Int64Ptr(way.GetTimestamp()),
			Changeset: thrift.Int64Ptr(way.GetChangesetId()),
			Uid:       thrift.Int32Ptr(way.GetUid()),
			UserSid:   thrift.Uint32Ptr(uint32(enc.sTable.getStrIndex(way.GetUser()))),
			Visible:   thrift.BoolPtr(way.GetVisible()),
		},
	}

	enc.wayNodesLen += len(way.GetNodes()) * 8
	lastNodeId := int64(0)
	for _, node := range way.GetNodes() {
		w.Refs = append(w.Refs, node.GetId()-lastNodeId)
		lastNodeId = node.GetId()
	}

	enc.wayTagsLen += len(way.GetTags()) * 8
	sortedTags := make([]string, 0, len(way.GetTags()))
	for k := range way.GetTags() {
		sortedTags = append(sortedTags, k)
	}
	sort.Strings(sortedTags)
	for _, k := range sortedTags {
		w.Keys = append(w.Keys, uint32(enc.sTable.getStrIndex(k)))
		w.Vals = append(w.Vals, uint32(enc.sTable.getStrIndex(way.GetTags()[k].(string))))
	}

	enc.ways = append(enc.ways, w)
}

func (enc *Encoder) encodeRelation(relation *entity.Relation) {
	r := &model_pb.Relation{
		Id: thrift.Int64Ptr(relation.GetId()),
		Info: &model_pb.Info{
			Version:   thrift.Int32Ptr(relation.GetVersion()),
			Timestamp: thrift.Int64Ptr(relation.GetTimestamp()),
			Changeset: thrift.Int64Ptr(relation.GetChangesetId()),
			Uid:       thrift.Int32Ptr(relation.GetUid()),
			UserSid:   thrift.Uint32Ptr(uint32(enc.sTable.getStrIndex(relation.GetUser()))),
			Visible:   thrift.BoolPtr(relation.GetVisible()),
		},
	}

	enc.relationMembersLen += len(relation.GetRelationMembers()) * 16
	lastMemberId := int64(0)
	for _, rm := range relation.GetRelationMembers() {
		r.RolesSid = append(r.RolesSid, int32(enc.sTable.getStrIndex(rm.GetRole())))

		r.Memids = append(r.Memids, rm.Exporter.GetId()-lastMemberId)
		lastMemberId = rm.Exporter.GetId()

		switch rm.Exporter.(type) {
		case *entity.Node:
			r.Types = append(r.Types, model_pb.Relation_NODE)
		case *entity.Way:
			r.Types = append(r.Types, model_pb.Relation_WAY)
		case *entity.Relation:
			r.Types = append(r.Types, model_pb.Relation_RELATION)
		default:
			panic("no entity type matched")
		}
	}

	enc.relationTagsLen += len(relation.GetTags()) * 8
	sortedTags := make([]string, 0, len(relation.GetTags()))
	for k := range relation.GetTags() {
		sortedTags = append(sortedTags, k)
	}
	sort.Strings(sortedTags)
	for _, k := range sortedTags {
		r.Keys = append(r.Keys, uint32(enc.sTable.getStrIndex(k)))
		r.Vals = append(r.Vals, uint32(enc.sTable.getStrIndex(relation.GetTags()[k].(string))))
	}

	enc.relations = append(enc.relations, r)
}
