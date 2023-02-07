// Code generated by ent, DO NOT EDIT.

package ent

import (
	"context"
	"errors"
	"fmt"
	"log"

	"github.com/ugent-library/deliver/ent/migrate"

	"github.com/ugent-library/deliver/ent/file"
	"github.com/ugent-library/deliver/ent/folder"
	"github.com/ugent-library/deliver/ent/space"

	"entgo.io/ent/dialect"
	"entgo.io/ent/dialect/sql"
	"entgo.io/ent/dialect/sql/sqlgraph"
)

// Client is the client that holds all ent builders.
type Client struct {
	config
	// Schema is the client for creating, migrating and dropping schema.
	Schema *migrate.Schema
	// File is the client for interacting with the File builders.
	File *FileClient
	// Folder is the client for interacting with the Folder builders.
	Folder *FolderClient
	// Space is the client for interacting with the Space builders.
	Space *SpaceClient
}

// NewClient creates a new client configured with the given options.
func NewClient(opts ...Option) *Client {
	cfg := config{log: log.Println, hooks: &hooks{}, inters: &inters{}}
	cfg.options(opts...)
	client := &Client{config: cfg}
	client.init()
	return client
}

func (c *Client) init() {
	c.Schema = migrate.NewSchema(c.driver)
	c.File = NewFileClient(c.config)
	c.Folder = NewFolderClient(c.config)
	c.Space = NewSpaceClient(c.config)
}

// Open opens a database/sql.DB specified by the driver name and
// the data source name, and returns a new client attached to it.
// Optional parameters can be added for configuring the client.
func Open(driverName, dataSourceName string, options ...Option) (*Client, error) {
	switch driverName {
	case dialect.MySQL, dialect.Postgres, dialect.SQLite:
		drv, err := sql.Open(driverName, dataSourceName)
		if err != nil {
			return nil, err
		}
		return NewClient(append(options, Driver(drv))...), nil
	default:
		return nil, fmt.Errorf("unsupported driver: %q", driverName)
	}
}

// Tx returns a new transactional client. The provided context
// is used until the transaction is committed or rolled back.
func (c *Client) Tx(ctx context.Context) (*Tx, error) {
	if _, ok := c.driver.(*txDriver); ok {
		return nil, errors.New("ent: cannot start a transaction within a transaction")
	}
	tx, err := newTx(ctx, c.driver)
	if err != nil {
		return nil, fmt.Errorf("ent: starting a transaction: %w", err)
	}
	cfg := c.config
	cfg.driver = tx
	return &Tx{
		ctx:    ctx,
		config: cfg,
		File:   NewFileClient(cfg),
		Folder: NewFolderClient(cfg),
		Space:  NewSpaceClient(cfg),
	}, nil
}

// BeginTx returns a transactional client with specified options.
func (c *Client) BeginTx(ctx context.Context, opts *sql.TxOptions) (*Tx, error) {
	if _, ok := c.driver.(*txDriver); ok {
		return nil, errors.New("ent: cannot start a transaction within a transaction")
	}
	tx, err := c.driver.(interface {
		BeginTx(context.Context, *sql.TxOptions) (dialect.Tx, error)
	}).BeginTx(ctx, opts)
	if err != nil {
		return nil, fmt.Errorf("ent: starting a transaction: %w", err)
	}
	cfg := c.config
	cfg.driver = &txDriver{tx: tx, drv: c.driver}
	return &Tx{
		ctx:    ctx,
		config: cfg,
		File:   NewFileClient(cfg),
		Folder: NewFolderClient(cfg),
		Space:  NewSpaceClient(cfg),
	}, nil
}

// Debug returns a new debug-client. It's used to get verbose logging on specific operations.
//
//	client.Debug().
//		File.
//		Query().
//		Count(ctx)
func (c *Client) Debug() *Client {
	if c.debug {
		return c
	}
	cfg := c.config
	cfg.driver = dialect.Debug(c.driver, c.log)
	client := &Client{config: cfg}
	client.init()
	return client
}

// Close closes the database connection and prevents new queries from starting.
func (c *Client) Close() error {
	return c.driver.Close()
}

// Use adds the mutation hooks to all the entity clients.
// In order to add hooks to a specific client, call: `client.Node.Use(...)`.
func (c *Client) Use(hooks ...Hook) {
	c.File.Use(hooks...)
	c.Folder.Use(hooks...)
	c.Space.Use(hooks...)
}

// Intercept adds the query interceptors to all the entity clients.
// In order to add interceptors to a specific client, call: `client.Node.Intercept(...)`.
func (c *Client) Intercept(interceptors ...Interceptor) {
	c.File.Intercept(interceptors...)
	c.Folder.Intercept(interceptors...)
	c.Space.Intercept(interceptors...)
}

// Mutate implements the ent.Mutator interface.
func (c *Client) Mutate(ctx context.Context, m Mutation) (Value, error) {
	switch m := m.(type) {
	case *FileMutation:
		return c.File.mutate(ctx, m)
	case *FolderMutation:
		return c.Folder.mutate(ctx, m)
	case *SpaceMutation:
		return c.Space.mutate(ctx, m)
	default:
		return nil, fmt.Errorf("ent: unknown mutation type %T", m)
	}
}

// FileClient is a client for the File schema.
type FileClient struct {
	config
}

// NewFileClient returns a client for the File from the given config.
func NewFileClient(c config) *FileClient {
	return &FileClient{config: c}
}

// Use adds a list of mutation hooks to the hooks stack.
// A call to `Use(f, g, h)` equals to `file.Hooks(f(g(h())))`.
func (c *FileClient) Use(hooks ...Hook) {
	c.hooks.File = append(c.hooks.File, hooks...)
}

// Use adds a list of query interceptors to the interceptors stack.
// A call to `Intercept(f, g, h)` equals to `file.Intercept(f(g(h())))`.
func (c *FileClient) Intercept(interceptors ...Interceptor) {
	c.inters.File = append(c.inters.File, interceptors...)
}

// Create returns a builder for creating a File entity.
func (c *FileClient) Create() *FileCreate {
	mutation := newFileMutation(c.config, OpCreate)
	return &FileCreate{config: c.config, hooks: c.Hooks(), mutation: mutation}
}

// CreateBulk returns a builder for creating a bulk of File entities.
func (c *FileClient) CreateBulk(builders ...*FileCreate) *FileCreateBulk {
	return &FileCreateBulk{config: c.config, builders: builders}
}

// Update returns an update builder for File.
func (c *FileClient) Update() *FileUpdate {
	mutation := newFileMutation(c.config, OpUpdate)
	return &FileUpdate{config: c.config, hooks: c.Hooks(), mutation: mutation}
}

// UpdateOne returns an update builder for the given entity.
func (c *FileClient) UpdateOne(f *File) *FileUpdateOne {
	mutation := newFileMutation(c.config, OpUpdateOne, withFile(f))
	return &FileUpdateOne{config: c.config, hooks: c.Hooks(), mutation: mutation}
}

// UpdateOneID returns an update builder for the given id.
func (c *FileClient) UpdateOneID(id string) *FileUpdateOne {
	mutation := newFileMutation(c.config, OpUpdateOne, withFileID(id))
	return &FileUpdateOne{config: c.config, hooks: c.Hooks(), mutation: mutation}
}

// Delete returns a delete builder for File.
func (c *FileClient) Delete() *FileDelete {
	mutation := newFileMutation(c.config, OpDelete)
	return &FileDelete{config: c.config, hooks: c.Hooks(), mutation: mutation}
}

// DeleteOne returns a builder for deleting the given entity.
func (c *FileClient) DeleteOne(f *File) *FileDeleteOne {
	return c.DeleteOneID(f.ID)
}

// DeleteOneID returns a builder for deleting the given entity by its id.
func (c *FileClient) DeleteOneID(id string) *FileDeleteOne {
	builder := c.Delete().Where(file.ID(id))
	builder.mutation.id = &id
	builder.mutation.op = OpDeleteOne
	return &FileDeleteOne{builder}
}

// Query returns a query builder for File.
func (c *FileClient) Query() *FileQuery {
	return &FileQuery{
		config: c.config,
		ctx:    &QueryContext{Type: TypeFile},
		inters: c.Interceptors(),
	}
}

// Get returns a File entity by its id.
func (c *FileClient) Get(ctx context.Context, id string) (*File, error) {
	return c.Query().Where(file.ID(id)).Only(ctx)
}

// GetX is like Get, but panics if an error occurs.
func (c *FileClient) GetX(ctx context.Context, id string) *File {
	obj, err := c.Get(ctx, id)
	if err != nil {
		panic(err)
	}
	return obj
}

// QueryFolder queries the folder edge of a File.
func (c *FileClient) QueryFolder(f *File) *FolderQuery {
	query := (&FolderClient{config: c.config}).Query()
	query.path = func(context.Context) (fromV *sql.Selector, _ error) {
		id := f.ID
		step := sqlgraph.NewStep(
			sqlgraph.From(file.Table, file.FieldID, id),
			sqlgraph.To(folder.Table, folder.FieldID),
			sqlgraph.Edge(sqlgraph.M2O, true, file.FolderTable, file.FolderColumn),
		)
		fromV = sqlgraph.Neighbors(f.driver.Dialect(), step)
		return fromV, nil
	}
	return query
}

// Hooks returns the client hooks.
func (c *FileClient) Hooks() []Hook {
	return c.hooks.File
}

// Interceptors returns the client interceptors.
func (c *FileClient) Interceptors() []Interceptor {
	return c.inters.File
}

func (c *FileClient) mutate(ctx context.Context, m *FileMutation) (Value, error) {
	switch m.Op() {
	case OpCreate:
		return (&FileCreate{config: c.config, hooks: c.Hooks(), mutation: m}).Save(ctx)
	case OpUpdate:
		return (&FileUpdate{config: c.config, hooks: c.Hooks(), mutation: m}).Save(ctx)
	case OpUpdateOne:
		return (&FileUpdateOne{config: c.config, hooks: c.Hooks(), mutation: m}).Save(ctx)
	case OpDelete, OpDeleteOne:
		return (&FileDelete{config: c.config, hooks: c.Hooks(), mutation: m}).Exec(ctx)
	default:
		return nil, fmt.Errorf("ent: unknown File mutation op: %q", m.Op())
	}
}

// FolderClient is a client for the Folder schema.
type FolderClient struct {
	config
}

// NewFolderClient returns a client for the Folder from the given config.
func NewFolderClient(c config) *FolderClient {
	return &FolderClient{config: c}
}

// Use adds a list of mutation hooks to the hooks stack.
// A call to `Use(f, g, h)` equals to `folder.Hooks(f(g(h())))`.
func (c *FolderClient) Use(hooks ...Hook) {
	c.hooks.Folder = append(c.hooks.Folder, hooks...)
}

// Use adds a list of query interceptors to the interceptors stack.
// A call to `Intercept(f, g, h)` equals to `folder.Intercept(f(g(h())))`.
func (c *FolderClient) Intercept(interceptors ...Interceptor) {
	c.inters.Folder = append(c.inters.Folder, interceptors...)
}

// Create returns a builder for creating a Folder entity.
func (c *FolderClient) Create() *FolderCreate {
	mutation := newFolderMutation(c.config, OpCreate)
	return &FolderCreate{config: c.config, hooks: c.Hooks(), mutation: mutation}
}

// CreateBulk returns a builder for creating a bulk of Folder entities.
func (c *FolderClient) CreateBulk(builders ...*FolderCreate) *FolderCreateBulk {
	return &FolderCreateBulk{config: c.config, builders: builders}
}

// Update returns an update builder for Folder.
func (c *FolderClient) Update() *FolderUpdate {
	mutation := newFolderMutation(c.config, OpUpdate)
	return &FolderUpdate{config: c.config, hooks: c.Hooks(), mutation: mutation}
}

// UpdateOne returns an update builder for the given entity.
func (c *FolderClient) UpdateOne(f *Folder) *FolderUpdateOne {
	mutation := newFolderMutation(c.config, OpUpdateOne, withFolder(f))
	return &FolderUpdateOne{config: c.config, hooks: c.Hooks(), mutation: mutation}
}

// UpdateOneID returns an update builder for the given id.
func (c *FolderClient) UpdateOneID(id string) *FolderUpdateOne {
	mutation := newFolderMutation(c.config, OpUpdateOne, withFolderID(id))
	return &FolderUpdateOne{config: c.config, hooks: c.Hooks(), mutation: mutation}
}

// Delete returns a delete builder for Folder.
func (c *FolderClient) Delete() *FolderDelete {
	mutation := newFolderMutation(c.config, OpDelete)
	return &FolderDelete{config: c.config, hooks: c.Hooks(), mutation: mutation}
}

// DeleteOne returns a builder for deleting the given entity.
func (c *FolderClient) DeleteOne(f *Folder) *FolderDeleteOne {
	return c.DeleteOneID(f.ID)
}

// DeleteOneID returns a builder for deleting the given entity by its id.
func (c *FolderClient) DeleteOneID(id string) *FolderDeleteOne {
	builder := c.Delete().Where(folder.ID(id))
	builder.mutation.id = &id
	builder.mutation.op = OpDeleteOne
	return &FolderDeleteOne{builder}
}

// Query returns a query builder for Folder.
func (c *FolderClient) Query() *FolderQuery {
	return &FolderQuery{
		config: c.config,
		ctx:    &QueryContext{Type: TypeFolder},
		inters: c.Interceptors(),
	}
}

// Get returns a Folder entity by its id.
func (c *FolderClient) Get(ctx context.Context, id string) (*Folder, error) {
	return c.Query().Where(folder.ID(id)).Only(ctx)
}

// GetX is like Get, but panics if an error occurs.
func (c *FolderClient) GetX(ctx context.Context, id string) *Folder {
	obj, err := c.Get(ctx, id)
	if err != nil {
		panic(err)
	}
	return obj
}

// QuerySpace queries the space edge of a Folder.
func (c *FolderClient) QuerySpace(f *Folder) *SpaceQuery {
	query := (&SpaceClient{config: c.config}).Query()
	query.path = func(context.Context) (fromV *sql.Selector, _ error) {
		id := f.ID
		step := sqlgraph.NewStep(
			sqlgraph.From(folder.Table, folder.FieldID, id),
			sqlgraph.To(space.Table, space.FieldID),
			sqlgraph.Edge(sqlgraph.M2O, true, folder.SpaceTable, folder.SpaceColumn),
		)
		fromV = sqlgraph.Neighbors(f.driver.Dialect(), step)
		return fromV, nil
	}
	return query
}

// QueryFiles queries the files edge of a Folder.
func (c *FolderClient) QueryFiles(f *Folder) *FileQuery {
	query := (&FileClient{config: c.config}).Query()
	query.path = func(context.Context) (fromV *sql.Selector, _ error) {
		id := f.ID
		step := sqlgraph.NewStep(
			sqlgraph.From(folder.Table, folder.FieldID, id),
			sqlgraph.To(file.Table, file.FieldID),
			sqlgraph.Edge(sqlgraph.O2M, false, folder.FilesTable, folder.FilesColumn),
		)
		fromV = sqlgraph.Neighbors(f.driver.Dialect(), step)
		return fromV, nil
	}
	return query
}

// Hooks returns the client hooks.
func (c *FolderClient) Hooks() []Hook {
	return c.hooks.Folder
}

// Interceptors returns the client interceptors.
func (c *FolderClient) Interceptors() []Interceptor {
	return c.inters.Folder
}

func (c *FolderClient) mutate(ctx context.Context, m *FolderMutation) (Value, error) {
	switch m.Op() {
	case OpCreate:
		return (&FolderCreate{config: c.config, hooks: c.Hooks(), mutation: m}).Save(ctx)
	case OpUpdate:
		return (&FolderUpdate{config: c.config, hooks: c.Hooks(), mutation: m}).Save(ctx)
	case OpUpdateOne:
		return (&FolderUpdateOne{config: c.config, hooks: c.Hooks(), mutation: m}).Save(ctx)
	case OpDelete, OpDeleteOne:
		return (&FolderDelete{config: c.config, hooks: c.Hooks(), mutation: m}).Exec(ctx)
	default:
		return nil, fmt.Errorf("ent: unknown Folder mutation op: %q", m.Op())
	}
}

// SpaceClient is a client for the Space schema.
type SpaceClient struct {
	config
}

// NewSpaceClient returns a client for the Space from the given config.
func NewSpaceClient(c config) *SpaceClient {
	return &SpaceClient{config: c}
}

// Use adds a list of mutation hooks to the hooks stack.
// A call to `Use(f, g, h)` equals to `space.Hooks(f(g(h())))`.
func (c *SpaceClient) Use(hooks ...Hook) {
	c.hooks.Space = append(c.hooks.Space, hooks...)
}

// Use adds a list of query interceptors to the interceptors stack.
// A call to `Intercept(f, g, h)` equals to `space.Intercept(f(g(h())))`.
func (c *SpaceClient) Intercept(interceptors ...Interceptor) {
	c.inters.Space = append(c.inters.Space, interceptors...)
}

// Create returns a builder for creating a Space entity.
func (c *SpaceClient) Create() *SpaceCreate {
	mutation := newSpaceMutation(c.config, OpCreate)
	return &SpaceCreate{config: c.config, hooks: c.Hooks(), mutation: mutation}
}

// CreateBulk returns a builder for creating a bulk of Space entities.
func (c *SpaceClient) CreateBulk(builders ...*SpaceCreate) *SpaceCreateBulk {
	return &SpaceCreateBulk{config: c.config, builders: builders}
}

// Update returns an update builder for Space.
func (c *SpaceClient) Update() *SpaceUpdate {
	mutation := newSpaceMutation(c.config, OpUpdate)
	return &SpaceUpdate{config: c.config, hooks: c.Hooks(), mutation: mutation}
}

// UpdateOne returns an update builder for the given entity.
func (c *SpaceClient) UpdateOne(s *Space) *SpaceUpdateOne {
	mutation := newSpaceMutation(c.config, OpUpdateOne, withSpace(s))
	return &SpaceUpdateOne{config: c.config, hooks: c.Hooks(), mutation: mutation}
}

// UpdateOneID returns an update builder for the given id.
func (c *SpaceClient) UpdateOneID(id string) *SpaceUpdateOne {
	mutation := newSpaceMutation(c.config, OpUpdateOne, withSpaceID(id))
	return &SpaceUpdateOne{config: c.config, hooks: c.Hooks(), mutation: mutation}
}

// Delete returns a delete builder for Space.
func (c *SpaceClient) Delete() *SpaceDelete {
	mutation := newSpaceMutation(c.config, OpDelete)
	return &SpaceDelete{config: c.config, hooks: c.Hooks(), mutation: mutation}
}

// DeleteOne returns a builder for deleting the given entity.
func (c *SpaceClient) DeleteOne(s *Space) *SpaceDeleteOne {
	return c.DeleteOneID(s.ID)
}

// DeleteOneID returns a builder for deleting the given entity by its id.
func (c *SpaceClient) DeleteOneID(id string) *SpaceDeleteOne {
	builder := c.Delete().Where(space.ID(id))
	builder.mutation.id = &id
	builder.mutation.op = OpDeleteOne
	return &SpaceDeleteOne{builder}
}

// Query returns a query builder for Space.
func (c *SpaceClient) Query() *SpaceQuery {
	return &SpaceQuery{
		config: c.config,
		ctx:    &QueryContext{Type: TypeSpace},
		inters: c.Interceptors(),
	}
}

// Get returns a Space entity by its id.
func (c *SpaceClient) Get(ctx context.Context, id string) (*Space, error) {
	return c.Query().Where(space.ID(id)).Only(ctx)
}

// GetX is like Get, but panics if an error occurs.
func (c *SpaceClient) GetX(ctx context.Context, id string) *Space {
	obj, err := c.Get(ctx, id)
	if err != nil {
		panic(err)
	}
	return obj
}

// QueryFolders queries the folders edge of a Space.
func (c *SpaceClient) QueryFolders(s *Space) *FolderQuery {
	query := (&FolderClient{config: c.config}).Query()
	query.path = func(context.Context) (fromV *sql.Selector, _ error) {
		id := s.ID
		step := sqlgraph.NewStep(
			sqlgraph.From(space.Table, space.FieldID, id),
			sqlgraph.To(folder.Table, folder.FieldID),
			sqlgraph.Edge(sqlgraph.O2M, false, space.FoldersTable, space.FoldersColumn),
		)
		fromV = sqlgraph.Neighbors(s.driver.Dialect(), step)
		return fromV, nil
	}
	return query
}

// Hooks returns the client hooks.
func (c *SpaceClient) Hooks() []Hook {
	return c.hooks.Space
}

// Interceptors returns the client interceptors.
func (c *SpaceClient) Interceptors() []Interceptor {
	return c.inters.Space
}

func (c *SpaceClient) mutate(ctx context.Context, m *SpaceMutation) (Value, error) {
	switch m.Op() {
	case OpCreate:
		return (&SpaceCreate{config: c.config, hooks: c.Hooks(), mutation: m}).Save(ctx)
	case OpUpdate:
		return (&SpaceUpdate{config: c.config, hooks: c.Hooks(), mutation: m}).Save(ctx)
	case OpUpdateOne:
		return (&SpaceUpdateOne{config: c.config, hooks: c.Hooks(), mutation: m}).Save(ctx)
	case OpDelete, OpDeleteOne:
		return (&SpaceDelete{config: c.config, hooks: c.Hooks(), mutation: m}).Exec(ctx)
	default:
		return nil, fmt.Errorf("ent: unknown Space mutation op: %q", m.Op())
	}
}
