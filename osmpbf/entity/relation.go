package entity

import (
	"github.com/twpayne/go-geom"
)

type Relation struct {
	*Entity

	relationMembers []*RelationMember
}

func NewRelation(id int64) *Relation {
	r := &Relation{
		Entity: NewEntity(id),
	}
	r.SetGeometry(r.getGeometry)

	return r
}

func (r *Relation) getGeometry() geom.T {
	geometries := make(map[int64]geom.T, len(r.relationMembers))
	r.getRelationGeometry(geometries)

	gc := geom.NewGeometryCollection()
	for _, geometry := range geometries {
		gc.MustPush(geometry)
	}

	return gc
}

func (r *Relation) getRelationGeometry(geometries map[int64]geom.T) {
	for _, rm := range r.relationMembers {
		switch rm.Exporter.(type) {
		case *Relation:
			rm.Exporter.(*Relation).getRelationGeometry(geometries)
		default:
			if _, ok := geometries[rm.GetId()]; ok {
				continue
			}
			geometries[rm.GetId()] = rm.GetGeometry()
		}
	}
}

func (r *Relation) GetRelationMembers() []*RelationMember {
	return r.relationMembers
}

func (r *Relation) SetRelationMembers(relationMembers []*RelationMember) {
	r.relationMembers = relationMembers
}
