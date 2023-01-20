package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"github.com/oklog/ulid/v2"
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
			DefaultFunc(func() string {
				return ulid.Make().String()
			}),
		field.String("name").
			Unique(),
		field.Strings("admins").
			Optional(),
		field.Time("created_at").
			Default(time.Now).
			Immutable(),
		field.Time("updated_at").
			Default(time.Now).
			UpdateDefault(time.Now),
	}
}

// Edges of the Space.
func (Space) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("folders", Folder.Type).
			Annotations(entsql.Annotation{
				OnDelete: entsql.Cascade,
			}),
	}
}
