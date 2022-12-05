package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"github.com/ugent-library/dilliver/ulid"
)

// Space holds the schema definition for the Space entity.
type Space struct {
	ent.Schema
}

// Fields of the Space.
func (Space) Fields() []ent.Field {
	return []ent.Field{
		field.String("id").
			Unique().
			DefaultFunc(ulid.MustGenerate),
		field.String("name").
			Unique(),
	}
}

// Edges of the Space.
func (Space) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("folders", Folder.Type),
	}
}
