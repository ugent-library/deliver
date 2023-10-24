// Code generated by ent, DO NOT EDIT.

package ent

import (
	"context"
	"errors"
	"fmt"
	"time"

	"entgo.io/ent/dialect"
	"entgo.io/ent/dialect/sql"
	"entgo.io/ent/dialect/sql/sqlgraph"
	"entgo.io/ent/schema/field"
	"github.com/ugent-library/deliver/ent/file"
	"github.com/ugent-library/deliver/ent/folder"
	"github.com/ugent-library/deliver/ent/space"
)

// FolderCreate is the builder for creating a Folder entity.
type FolderCreate struct {
	config
	mutation *FolderMutation
	hooks    []Hook
	conflict []sql.ConflictOption
}

// SetSpaceID sets the "space_id" field.
func (fc *FolderCreate) SetSpaceID(s string) *FolderCreate {
	fc.mutation.SetSpaceID(s)
	return fc
}

// SetName sets the "name" field.
func (fc *FolderCreate) SetName(s string) *FolderCreate {
	fc.mutation.SetName(s)
	return fc
}

// SetCreatedAt sets the "created_at" field.
func (fc *FolderCreate) SetCreatedAt(t time.Time) *FolderCreate {
	fc.mutation.SetCreatedAt(t)
	return fc
}

// SetNillableCreatedAt sets the "created_at" field if the given value is not nil.
func (fc *FolderCreate) SetNillableCreatedAt(t *time.Time) *FolderCreate {
	if t != nil {
		fc.SetCreatedAt(*t)
	}
	return fc
}

// SetUpdatedAt sets the "updated_at" field.
func (fc *FolderCreate) SetUpdatedAt(t time.Time) *FolderCreate {
	fc.mutation.SetUpdatedAt(t)
	return fc
}

// SetNillableUpdatedAt sets the "updated_at" field if the given value is not nil.
func (fc *FolderCreate) SetNillableUpdatedAt(t *time.Time) *FolderCreate {
	if t != nil {
		fc.SetUpdatedAt(*t)
	}
	return fc
}

// SetExpiresAt sets the "expires_at" field.
func (fc *FolderCreate) SetExpiresAt(t time.Time) *FolderCreate {
	fc.mutation.SetExpiresAt(t)
	return fc
}

// SetNillableExpiresAt sets the "expires_at" field if the given value is not nil.
func (fc *FolderCreate) SetNillableExpiresAt(t *time.Time) *FolderCreate {
	if t != nil {
		fc.SetExpiresAt(*t)
	}
	return fc
}

// SetID sets the "id" field.
func (fc *FolderCreate) SetID(s string) *FolderCreate {
	fc.mutation.SetID(s)
	return fc
}

// SetNillableID sets the "id" field if the given value is not nil.
func (fc *FolderCreate) SetNillableID(s *string) *FolderCreate {
	if s != nil {
		fc.SetID(*s)
	}
	return fc
}

// SetSpace sets the "space" edge to the Space entity.
func (fc *FolderCreate) SetSpace(s *Space) *FolderCreate {
	return fc.SetSpaceID(s.ID)
}

// AddFileIDs adds the "files" edge to the File entity by IDs.
func (fc *FolderCreate) AddFileIDs(ids ...string) *FolderCreate {
	fc.mutation.AddFileIDs(ids...)
	return fc
}

// AddFiles adds the "files" edges to the File entity.
func (fc *FolderCreate) AddFiles(f ...*File) *FolderCreate {
	ids := make([]string, len(f))
	for i := range f {
		ids[i] = f[i].ID
	}
	return fc.AddFileIDs(ids...)
}

// Mutation returns the FolderMutation object of the builder.
func (fc *FolderCreate) Mutation() *FolderMutation {
	return fc.mutation
}

// Save creates the Folder in the database.
func (fc *FolderCreate) Save(ctx context.Context) (*Folder, error) {
	fc.defaults()
	return withHooks(ctx, fc.sqlSave, fc.mutation, fc.hooks)
}

// SaveX calls Save and panics if Save returns an error.
func (fc *FolderCreate) SaveX(ctx context.Context) *Folder {
	v, err := fc.Save(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Exec executes the query.
func (fc *FolderCreate) Exec(ctx context.Context) error {
	_, err := fc.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (fc *FolderCreate) ExecX(ctx context.Context) {
	if err := fc.Exec(ctx); err != nil {
		panic(err)
	}
}

// defaults sets the default values of the builder before save.
func (fc *FolderCreate) defaults() {
	if _, ok := fc.mutation.CreatedAt(); !ok {
		v := folder.DefaultCreatedAt()
		fc.mutation.SetCreatedAt(v)
	}
	if _, ok := fc.mutation.UpdatedAt(); !ok {
		v := folder.DefaultUpdatedAt()
		fc.mutation.SetUpdatedAt(v)
	}
	if _, ok := fc.mutation.ID(); !ok {
		v := folder.DefaultID()
		fc.mutation.SetID(v)
	}
}

// check runs all checks and user-defined validators on the builder.
func (fc *FolderCreate) check() error {
	if _, ok := fc.mutation.SpaceID(); !ok {
		return &ValidationError{Name: "space_id", err: errors.New(`ent: missing required field "Folder.space_id"`)}
	}
	if _, ok := fc.mutation.Name(); !ok {
		return &ValidationError{Name: "name", err: errors.New(`ent: missing required field "Folder.name"`)}
	}
	if _, ok := fc.mutation.CreatedAt(); !ok {
		return &ValidationError{Name: "created_at", err: errors.New(`ent: missing required field "Folder.created_at"`)}
	}
	if _, ok := fc.mutation.UpdatedAt(); !ok {
		return &ValidationError{Name: "updated_at", err: errors.New(`ent: missing required field "Folder.updated_at"`)}
	}
	if v, ok := fc.mutation.ID(); ok {
		if err := folder.IDValidator(v); err != nil {
			return &ValidationError{Name: "id", err: fmt.Errorf(`ent: validator failed for field "Folder.id": %w`, err)}
		}
	}
	if _, ok := fc.mutation.SpaceID(); !ok {
		return &ValidationError{Name: "space", err: errors.New(`ent: missing required edge "Folder.space"`)}
	}
	return nil
}

func (fc *FolderCreate) sqlSave(ctx context.Context) (*Folder, error) {
	if err := fc.check(); err != nil {
		return nil, err
	}
	_node, _spec := fc.createSpec()
	if err := sqlgraph.CreateNode(ctx, fc.driver, _spec); err != nil {
		if sqlgraph.IsConstraintError(err) {
			err = &ConstraintError{msg: err.Error(), wrap: err}
		}
		return nil, err
	}
	if _spec.ID.Value != nil {
		if id, ok := _spec.ID.Value.(string); ok {
			_node.ID = id
		} else {
			return nil, fmt.Errorf("unexpected Folder.ID type: %T", _spec.ID.Value)
		}
	}
	fc.mutation.id = &_node.ID
	fc.mutation.done = true
	return _node, nil
}

func (fc *FolderCreate) createSpec() (*Folder, *sqlgraph.CreateSpec) {
	var (
		_node = &Folder{config: fc.config}
		_spec = sqlgraph.NewCreateSpec(folder.Table, sqlgraph.NewFieldSpec(folder.FieldID, field.TypeString))
	)
	_spec.OnConflict = fc.conflict
	if id, ok := fc.mutation.ID(); ok {
		_node.ID = id
		_spec.ID.Value = id
	}
	if value, ok := fc.mutation.Name(); ok {
		_spec.SetField(folder.FieldName, field.TypeString, value)
		_node.Name = value
	}
	if value, ok := fc.mutation.CreatedAt(); ok {
		_spec.SetField(folder.FieldCreatedAt, field.TypeTime, value)
		_node.CreatedAt = value
	}
	if value, ok := fc.mutation.UpdatedAt(); ok {
		_spec.SetField(folder.FieldUpdatedAt, field.TypeTime, value)
		_node.UpdatedAt = value
	}
	if value, ok := fc.mutation.ExpiresAt(); ok {
		_spec.SetField(folder.FieldExpiresAt, field.TypeTime, value)
		_node.ExpiresAt = value
	}
	if nodes := fc.mutation.SpaceIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: true,
			Table:   folder.SpaceTable,
			Columns: []string{folder.SpaceColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(space.FieldID, field.TypeString),
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_node.SpaceID = nodes[0]
		_spec.Edges = append(_spec.Edges, edge)
	}
	if nodes := fc.mutation.FilesIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   folder.FilesTable,
			Columns: []string{folder.FilesColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(file.FieldID, field.TypeString),
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges = append(_spec.Edges, edge)
	}
	return _node, _spec
}

// OnConflict allows configuring the `ON CONFLICT` / `ON DUPLICATE KEY` clause
// of the `INSERT` statement. For example:
//
//	client.Folder.Create().
//		SetSpaceID(v).
//		OnConflict(
//			// Update the row with the new values
//			// the was proposed for insertion.
//			sql.ResolveWithNewValues(),
//		).
//		// Override some of the fields with custom
//		// update values.
//		Update(func(u *ent.FolderUpsert) {
//			SetSpaceID(v+v).
//		}).
//		Exec(ctx)
func (fc *FolderCreate) OnConflict(opts ...sql.ConflictOption) *FolderUpsertOne {
	fc.conflict = opts
	return &FolderUpsertOne{
		create: fc,
	}
}

// OnConflictColumns calls `OnConflict` and configures the columns
// as conflict target. Using this option is equivalent to using:
//
//	client.Folder.Create().
//		OnConflict(sql.ConflictColumns(columns...)).
//		Exec(ctx)
func (fc *FolderCreate) OnConflictColumns(columns ...string) *FolderUpsertOne {
	fc.conflict = append(fc.conflict, sql.ConflictColumns(columns...))
	return &FolderUpsertOne{
		create: fc,
	}
}

type (
	// FolderUpsertOne is the builder for "upsert"-ing
	//  one Folder node.
	FolderUpsertOne struct {
		create *FolderCreate
	}

	// FolderUpsert is the "OnConflict" setter.
	FolderUpsert struct {
		*sql.UpdateSet
	}
)

// SetSpaceID sets the "space_id" field.
func (u *FolderUpsert) SetSpaceID(v string) *FolderUpsert {
	u.Set(folder.FieldSpaceID, v)
	return u
}

// UpdateSpaceID sets the "space_id" field to the value that was provided on create.
func (u *FolderUpsert) UpdateSpaceID() *FolderUpsert {
	u.SetExcluded(folder.FieldSpaceID)
	return u
}

// SetName sets the "name" field.
func (u *FolderUpsert) SetName(v string) *FolderUpsert {
	u.Set(folder.FieldName, v)
	return u
}

// UpdateName sets the "name" field to the value that was provided on create.
func (u *FolderUpsert) UpdateName() *FolderUpsert {
	u.SetExcluded(folder.FieldName)
	return u
}

// SetUpdatedAt sets the "updated_at" field.
func (u *FolderUpsert) SetUpdatedAt(v time.Time) *FolderUpsert {
	u.Set(folder.FieldUpdatedAt, v)
	return u
}

// UpdateUpdatedAt sets the "updated_at" field to the value that was provided on create.
func (u *FolderUpsert) UpdateUpdatedAt() *FolderUpsert {
	u.SetExcluded(folder.FieldUpdatedAt)
	return u
}

// SetExpiresAt sets the "expires_at" field.
func (u *FolderUpsert) SetExpiresAt(v time.Time) *FolderUpsert {
	u.Set(folder.FieldExpiresAt, v)
	return u
}

// UpdateExpiresAt sets the "expires_at" field to the value that was provided on create.
func (u *FolderUpsert) UpdateExpiresAt() *FolderUpsert {
	u.SetExcluded(folder.FieldExpiresAt)
	return u
}

// ClearExpiresAt clears the value of the "expires_at" field.
func (u *FolderUpsert) ClearExpiresAt() *FolderUpsert {
	u.SetNull(folder.FieldExpiresAt)
	return u
}

// UpdateNewValues updates the mutable fields using the new values that were set on create except the ID field.
// Using this option is equivalent to using:
//
//	client.Folder.Create().
//		OnConflict(
//			sql.ResolveWithNewValues(),
//			sql.ResolveWith(func(u *sql.UpdateSet) {
//				u.SetIgnore(folder.FieldID)
//			}),
//		).
//		Exec(ctx)
func (u *FolderUpsertOne) UpdateNewValues() *FolderUpsertOne {
	u.create.conflict = append(u.create.conflict, sql.ResolveWithNewValues())
	u.create.conflict = append(u.create.conflict, sql.ResolveWith(func(s *sql.UpdateSet) {
		if _, exists := u.create.mutation.ID(); exists {
			s.SetIgnore(folder.FieldID)
		}
		if _, exists := u.create.mutation.CreatedAt(); exists {
			s.SetIgnore(folder.FieldCreatedAt)
		}
	}))
	return u
}

// Ignore sets each column to itself in case of conflict.
// Using this option is equivalent to using:
//
//	client.Folder.Create().
//	    OnConflict(sql.ResolveWithIgnore()).
//	    Exec(ctx)
func (u *FolderUpsertOne) Ignore() *FolderUpsertOne {
	u.create.conflict = append(u.create.conflict, sql.ResolveWithIgnore())
	return u
}

// DoNothing configures the conflict_action to `DO NOTHING`.
// Supported only by SQLite and PostgreSQL.
func (u *FolderUpsertOne) DoNothing() *FolderUpsertOne {
	u.create.conflict = append(u.create.conflict, sql.DoNothing())
	return u
}

// Update allows overriding fields `UPDATE` values. See the FolderCreate.OnConflict
// documentation for more info.
func (u *FolderUpsertOne) Update(set func(*FolderUpsert)) *FolderUpsertOne {
	u.create.conflict = append(u.create.conflict, sql.ResolveWith(func(update *sql.UpdateSet) {
		set(&FolderUpsert{UpdateSet: update})
	}))
	return u
}

// SetSpaceID sets the "space_id" field.
func (u *FolderUpsertOne) SetSpaceID(v string) *FolderUpsertOne {
	return u.Update(func(s *FolderUpsert) {
		s.SetSpaceID(v)
	})
}

// UpdateSpaceID sets the "space_id" field to the value that was provided on create.
func (u *FolderUpsertOne) UpdateSpaceID() *FolderUpsertOne {
	return u.Update(func(s *FolderUpsert) {
		s.UpdateSpaceID()
	})
}

// SetName sets the "name" field.
func (u *FolderUpsertOne) SetName(v string) *FolderUpsertOne {
	return u.Update(func(s *FolderUpsert) {
		s.SetName(v)
	})
}

// UpdateName sets the "name" field to the value that was provided on create.
func (u *FolderUpsertOne) UpdateName() *FolderUpsertOne {
	return u.Update(func(s *FolderUpsert) {
		s.UpdateName()
	})
}

// SetUpdatedAt sets the "updated_at" field.
func (u *FolderUpsertOne) SetUpdatedAt(v time.Time) *FolderUpsertOne {
	return u.Update(func(s *FolderUpsert) {
		s.SetUpdatedAt(v)
	})
}

// UpdateUpdatedAt sets the "updated_at" field to the value that was provided on create.
func (u *FolderUpsertOne) UpdateUpdatedAt() *FolderUpsertOne {
	return u.Update(func(s *FolderUpsert) {
		s.UpdateUpdatedAt()
	})
}

// SetExpiresAt sets the "expires_at" field.
func (u *FolderUpsertOne) SetExpiresAt(v time.Time) *FolderUpsertOne {
	return u.Update(func(s *FolderUpsert) {
		s.SetExpiresAt(v)
	})
}

// UpdateExpiresAt sets the "expires_at" field to the value that was provided on create.
func (u *FolderUpsertOne) UpdateExpiresAt() *FolderUpsertOne {
	return u.Update(func(s *FolderUpsert) {
		s.UpdateExpiresAt()
	})
}

// ClearExpiresAt clears the value of the "expires_at" field.
func (u *FolderUpsertOne) ClearExpiresAt() *FolderUpsertOne {
	return u.Update(func(s *FolderUpsert) {
		s.ClearExpiresAt()
	})
}

// Exec executes the query.
func (u *FolderUpsertOne) Exec(ctx context.Context) error {
	if len(u.create.conflict) == 0 {
		return errors.New("ent: missing options for FolderCreate.OnConflict")
	}
	return u.create.Exec(ctx)
}

// ExecX is like Exec, but panics if an error occurs.
func (u *FolderUpsertOne) ExecX(ctx context.Context) {
	if err := u.create.Exec(ctx); err != nil {
		panic(err)
	}
}

// Exec executes the UPSERT query and returns the inserted/updated ID.
func (u *FolderUpsertOne) ID(ctx context.Context) (id string, err error) {
	if u.create.driver.Dialect() == dialect.MySQL {
		// In case of "ON CONFLICT", there is no way to get back non-numeric ID
		// fields from the database since MySQL does not support the RETURNING clause.
		return id, errors.New("ent: FolderUpsertOne.ID is not supported by MySQL driver. Use FolderUpsertOne.Exec instead")
	}
	node, err := u.create.Save(ctx)
	if err != nil {
		return id, err
	}
	return node.ID, nil
}

// IDX is like ID, but panics if an error occurs.
func (u *FolderUpsertOne) IDX(ctx context.Context) string {
	id, err := u.ID(ctx)
	if err != nil {
		panic(err)
	}
	return id
}

// FolderCreateBulk is the builder for creating many Folder entities in bulk.
type FolderCreateBulk struct {
	config
	builders []*FolderCreate
	conflict []sql.ConflictOption
}

// Save creates the Folder entities in the database.
func (fcb *FolderCreateBulk) Save(ctx context.Context) ([]*Folder, error) {
	specs := make([]*sqlgraph.CreateSpec, len(fcb.builders))
	nodes := make([]*Folder, len(fcb.builders))
	mutators := make([]Mutator, len(fcb.builders))
	for i := range fcb.builders {
		func(i int, root context.Context) {
			builder := fcb.builders[i]
			builder.defaults()
			var mut Mutator = MutateFunc(func(ctx context.Context, m Mutation) (Value, error) {
				mutation, ok := m.(*FolderMutation)
				if !ok {
					return nil, fmt.Errorf("unexpected mutation type %T", m)
				}
				if err := builder.check(); err != nil {
					return nil, err
				}
				builder.mutation = mutation
				var err error
				nodes[i], specs[i] = builder.createSpec()
				if i < len(mutators)-1 {
					_, err = mutators[i+1].Mutate(root, fcb.builders[i+1].mutation)
				} else {
					spec := &sqlgraph.BatchCreateSpec{Nodes: specs}
					spec.OnConflict = fcb.conflict
					// Invoke the actual operation on the latest mutation in the chain.
					if err = sqlgraph.BatchCreate(ctx, fcb.driver, spec); err != nil {
						if sqlgraph.IsConstraintError(err) {
							err = &ConstraintError{msg: err.Error(), wrap: err}
						}
					}
				}
				if err != nil {
					return nil, err
				}
				mutation.id = &nodes[i].ID
				mutation.done = true
				return nodes[i], nil
			})
			for i := len(builder.hooks) - 1; i >= 0; i-- {
				mut = builder.hooks[i](mut)
			}
			mutators[i] = mut
		}(i, ctx)
	}
	if len(mutators) > 0 {
		if _, err := mutators[0].Mutate(ctx, fcb.builders[0].mutation); err != nil {
			return nil, err
		}
	}
	return nodes, nil
}

// SaveX is like Save, but panics if an error occurs.
func (fcb *FolderCreateBulk) SaveX(ctx context.Context) []*Folder {
	v, err := fcb.Save(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Exec executes the query.
func (fcb *FolderCreateBulk) Exec(ctx context.Context) error {
	_, err := fcb.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (fcb *FolderCreateBulk) ExecX(ctx context.Context) {
	if err := fcb.Exec(ctx); err != nil {
		panic(err)
	}
}

// OnConflict allows configuring the `ON CONFLICT` / `ON DUPLICATE KEY` clause
// of the `INSERT` statement. For example:
//
//	client.Folder.CreateBulk(builders...).
//		OnConflict(
//			// Update the row with the new values
//			// the was proposed for insertion.
//			sql.ResolveWithNewValues(),
//		).
//		// Override some of the fields with custom
//		// update values.
//		Update(func(u *ent.FolderUpsert) {
//			SetSpaceID(v+v).
//		}).
//		Exec(ctx)
func (fcb *FolderCreateBulk) OnConflict(opts ...sql.ConflictOption) *FolderUpsertBulk {
	fcb.conflict = opts
	return &FolderUpsertBulk{
		create: fcb,
	}
}

// OnConflictColumns calls `OnConflict` and configures the columns
// as conflict target. Using this option is equivalent to using:
//
//	client.Folder.Create().
//		OnConflict(sql.ConflictColumns(columns...)).
//		Exec(ctx)
func (fcb *FolderCreateBulk) OnConflictColumns(columns ...string) *FolderUpsertBulk {
	fcb.conflict = append(fcb.conflict, sql.ConflictColumns(columns...))
	return &FolderUpsertBulk{
		create: fcb,
	}
}

// FolderUpsertBulk is the builder for "upsert"-ing
// a bulk of Folder nodes.
type FolderUpsertBulk struct {
	create *FolderCreateBulk
}

// UpdateNewValues updates the mutable fields using the new values that
// were set on create. Using this option is equivalent to using:
//
//	client.Folder.Create().
//		OnConflict(
//			sql.ResolveWithNewValues(),
//			sql.ResolveWith(func(u *sql.UpdateSet) {
//				u.SetIgnore(folder.FieldID)
//			}),
//		).
//		Exec(ctx)
func (u *FolderUpsertBulk) UpdateNewValues() *FolderUpsertBulk {
	u.create.conflict = append(u.create.conflict, sql.ResolveWithNewValues())
	u.create.conflict = append(u.create.conflict, sql.ResolveWith(func(s *sql.UpdateSet) {
		for _, b := range u.create.builders {
			if _, exists := b.mutation.ID(); exists {
				s.SetIgnore(folder.FieldID)
			}
			if _, exists := b.mutation.CreatedAt(); exists {
				s.SetIgnore(folder.FieldCreatedAt)
			}
		}
	}))
	return u
}

// Ignore sets each column to itself in case of conflict.
// Using this option is equivalent to using:
//
//	client.Folder.Create().
//		OnConflict(sql.ResolveWithIgnore()).
//		Exec(ctx)
func (u *FolderUpsertBulk) Ignore() *FolderUpsertBulk {
	u.create.conflict = append(u.create.conflict, sql.ResolveWithIgnore())
	return u
}

// DoNothing configures the conflict_action to `DO NOTHING`.
// Supported only by SQLite and PostgreSQL.
func (u *FolderUpsertBulk) DoNothing() *FolderUpsertBulk {
	u.create.conflict = append(u.create.conflict, sql.DoNothing())
	return u
}

// Update allows overriding fields `UPDATE` values. See the FolderCreateBulk.OnConflict
// documentation for more info.
func (u *FolderUpsertBulk) Update(set func(*FolderUpsert)) *FolderUpsertBulk {
	u.create.conflict = append(u.create.conflict, sql.ResolveWith(func(update *sql.UpdateSet) {
		set(&FolderUpsert{UpdateSet: update})
	}))
	return u
}

// SetSpaceID sets the "space_id" field.
func (u *FolderUpsertBulk) SetSpaceID(v string) *FolderUpsertBulk {
	return u.Update(func(s *FolderUpsert) {
		s.SetSpaceID(v)
	})
}

// UpdateSpaceID sets the "space_id" field to the value that was provided on create.
func (u *FolderUpsertBulk) UpdateSpaceID() *FolderUpsertBulk {
	return u.Update(func(s *FolderUpsert) {
		s.UpdateSpaceID()
	})
}

// SetName sets the "name" field.
func (u *FolderUpsertBulk) SetName(v string) *FolderUpsertBulk {
	return u.Update(func(s *FolderUpsert) {
		s.SetName(v)
	})
}

// UpdateName sets the "name" field to the value that was provided on create.
func (u *FolderUpsertBulk) UpdateName() *FolderUpsertBulk {
	return u.Update(func(s *FolderUpsert) {
		s.UpdateName()
	})
}

// SetUpdatedAt sets the "updated_at" field.
func (u *FolderUpsertBulk) SetUpdatedAt(v time.Time) *FolderUpsertBulk {
	return u.Update(func(s *FolderUpsert) {
		s.SetUpdatedAt(v)
	})
}

// UpdateUpdatedAt sets the "updated_at" field to the value that was provided on create.
func (u *FolderUpsertBulk) UpdateUpdatedAt() *FolderUpsertBulk {
	return u.Update(func(s *FolderUpsert) {
		s.UpdateUpdatedAt()
	})
}

// SetExpiresAt sets the "expires_at" field.
func (u *FolderUpsertBulk) SetExpiresAt(v time.Time) *FolderUpsertBulk {
	return u.Update(func(s *FolderUpsert) {
		s.SetExpiresAt(v)
	})
}

// UpdateExpiresAt sets the "expires_at" field to the value that was provided on create.
func (u *FolderUpsertBulk) UpdateExpiresAt() *FolderUpsertBulk {
	return u.Update(func(s *FolderUpsert) {
		s.UpdateExpiresAt()
	})
}

// ClearExpiresAt clears the value of the "expires_at" field.
func (u *FolderUpsertBulk) ClearExpiresAt() *FolderUpsertBulk {
	return u.Update(func(s *FolderUpsert) {
		s.ClearExpiresAt()
	})
}

// Exec executes the query.
func (u *FolderUpsertBulk) Exec(ctx context.Context) error {
	for i, b := range u.create.builders {
		if len(b.conflict) != 0 {
			return fmt.Errorf("ent: OnConflict was set for builder %d. Set it on the FolderCreateBulk instead", i)
		}
	}
	if len(u.create.conflict) == 0 {
		return errors.New("ent: missing options for FolderCreateBulk.OnConflict")
	}
	return u.create.Exec(ctx)
}

// ExecX is like Exec, but panics if an error occurs.
func (u *FolderUpsertBulk) ExecX(ctx context.Context) {
	if err := u.create.Exec(ctx); err != nil {
		panic(err)
	}
}
