package entity

type RelationMember struct {
	Exporter

	role string
}

func NewRelationMember(exp Exporter, r string) *RelationMember {
	return &RelationMember{
		Exporter: exp,
		role:     r,
	}
}

func (rm *RelationMember) GetRole() string {
	return rm.role
}

func (rm *RelationMember) SetRole(role string) {
	rm.role = role
}
