package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"github.com/ugent-library/dilliver/ulid"
)

// Folder holds the schema definition for the Folder entity.
type Folder struct {
	ent.Schema
}

// Fields of the Folder.
func (Folder) Fields() []ent.Field {
	return []ent.Field{
		field.String("id").
			Unique().
			DefaultFunc(ulid.MustGenerate),
		field.String("space_id"),
		field.String("name"),
		field.Time("created_at").
			Default(time.Now).
			Immutable(),
		field.Time("updated_at").
			Default(time.Now).
			UpdateDefault(time.Now),
		field.Time("expires_at").
			Optional(),
	}
}

// Edges of the Folder.
func (Folder) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("space", Space.Type).
			Ref("folders").
			Unique().
			Required().
			Field("space_id"),
		edge.To("files", File.Type),
	}
}
