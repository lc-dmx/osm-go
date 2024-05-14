package entity

import (
	"strconv"

	"github.com/twpayne/go-geom"
	"github.com/twpayne/go-geom/encoding/geojson"
	"github.com/twpayne/go-geom/encoding/wkt"
)

type Exporter interface {
	GetVisible() bool
	GetVersion() int32
	GetUid() int32
	GetId() int64
	GetChangesetId() int64
	GetTimestamp() int64
	GetUser() string
	GetTags() map[string]interface{}
	GetGeometry() geom.T

	ToWKT() (string, error)
	ToJSON() (string, error)
}

// unused protection
var (
	_ Exporter = &Entity{}
	_ Exporter = &Node{}
	_ Exporter = &Way{}
	_ Exporter = &Relation{}
)

type Entity struct {
	visible     bool
	version     int32
	uid         int32
	id          int64
	changesetId int64
	timestamp   int64
	user        string
	tags        map[string]interface{}
	getGeometry func() geom.T
}

func NewEntity(id int64) *Entity {
	return &Entity{
		visible: true,
		id:      id,
		tags:    make(map[string]interface{}, 0),
	}
}

func (e *Entity) GetVisible() bool {
	return e.visible
}

func (e *Entity) GetVersion() int32 {
	return e.version
}

func (e *Entity) GetUid() int32 {
	return e.uid
}

func (e *Entity) GetId() int64 {
	return e.id
}

func (e *Entity) GetChangesetId() int64 {
	return e.changesetId
}

func (e *Entity) GetTimestamp() int64 {
	return e.timestamp
}

func (e *Entity) GetUser() string {
	return e.user
}

func (e *Entity) GetTags() map[string]interface{} {
	return e.tags
}

func (e *Entity) GetGeometry() geom.T {
	return e.getGeometry()
}

func (e *Entity) SetVisible(visible bool) {
	e.visible = visible
}

func (e *Entity) SetVersion(version int32) {
	e.version = version
}

func (e *Entity) SetUid(uid int32) {
	e.uid = uid
}

func (e *Entity) SetId(id int64) {
	e.id = id
}

func (e *Entity) SetChangesetId(changesetId int64) {
	e.changesetId = changesetId
}

func (e *Entity) SetTimestamp(timestamp int64) {
	e.timestamp = timestamp
}

func (e *Entity) SetUser(user string) {
	e.user = user
}

func (e *Entity) SetTags(tags map[string]interface{}) {
	for k, v := range tags {
		e.tags[k] = v
	}
}

func (e *Entity) SetGeometry(getGeometry func() geom.T) {
	e.getGeometry = getGeometry
}

func (e *Entity) ToWKT() (string, error) {
	return wkt.Marshal(e.getGeometry())
}

func (e *Entity) ToJSON() (string, error) {
	e.tags["id"] = strconv.FormatInt(e.id, 10)
	f := &geojson.Feature{
		ID:         e.tags["id"].(string),
		Geometry:   e.getGeometry(),
		Properties: e.tags,
	}
	b, err := f.MarshalJSON()

	return string(b), err
}
