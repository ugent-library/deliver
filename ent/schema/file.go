package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"github.com/ugent-library/dilliver/ulid"
)

// File holds the schema definition for the File entity.
type File struct {
	ent.Schema
}

// Fields of the File.
func (File) Fields() []ent.Field {
	return []ent.Field{
		field.String("id").
			Unique().
			DefaultFunc(ulid.MustGenerate),
		field.String("folder_id"),
		field.String("md5"),
		field.String("name"),
		field.Int64("size"),
		field.String("content_type"),
		field.Time("created_at").
			Default(time.Now).
			Immutable(),
		field.Time("updated_at").
			Default(time.Now).
			UpdateDefault(time.Now),
	}
}

// Edges of the File.
func (File) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("folder", Folder.Type).
			Ref("files").
			Unique().
			Required().
			Field("folder_id"),
	}
}
