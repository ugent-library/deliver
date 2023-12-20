create table users (
    id text primary key,
    username text unique not null,
    name text not null,
    email text not null,
    remember_token text unique not null,
    created_at timestamptz not null,
    updated_at timestamptz not null
);

create table spaces (
    id text primary key,
    name text not null,
    admins jsonb,
    created_at timestamptz not null,
    updated_at timestamptz not null
);

create table folders (
    id text primary key,
    name text not null,
    created_at timestamptz not null,
    updated_at timestamptz not null,
    expires_at timestamptz,
    space_id text not null references spaces(id) on delete cascade
);

create table files (
    id text primary key,
    md5 text not null,
    name text not null,
    size bigint not null,
    content_type text not null,
    downloads bigint not null default 0,
    created_at timestamptz not null,
    updated_at timestamptz not null,
    folder_id text not null references folders(id) on delete cascade
);

---- create above / drop below ----

drop table files;
drop table folders;
drop table spaces;
drop table users;
