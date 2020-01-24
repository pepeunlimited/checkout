package schema

import (
	"github.com/facebookincubator/ent"
	"github.com/facebookincubator/ent/schema/field"
)

type Vault struct {
	ent.Schema
}

func (Vault) Fields() []ent.Field {
	return []ent.Field{
		field.Int64("amount"),
		field.String("reference_number").Unique().MaxLen(36),
	}
}

func (Vault) Edges() []ent.Edge {
	return []ent.Edge{}
}