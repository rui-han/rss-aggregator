# Introduction

This is a [RSS](https://en.wikipedia.org/wiki/RSS) feed aggregator, which is a web server that allow clients to:

- Add RSS feeds to be collected
- Follow and unfollow RSS feeds that other users have added
- Fetch all of the latest posts from the RSS feeds they follow

You can use this project to keep up with your favorite blogs, news sites, podcasts, and more!

# PostgreSQL

## 1. Install:

```zsh
brew install postgresql
```

## 2. Check version:

```zsh
psql --version
```

## 3. Start Postgres server in the background

```zsh
brew services start postgresql
```

## 4. Connect to the server

Download and install [pgAdmin](https://www.pgadmin.org/). Then create a new server connection:

- Host: `localhost`
- Port: `5432`
- Username: postgres
- Password: your password

## 5. Create and query a database

```sql
SELECT version();
```

# sqlc and goose

## 1. Installation

```zsh
brew install sqlc
```

Then run `sqlc version` to make sure it's installed correctly.

```zsh
brew install goose
```

Then run `goose -version` to make sure it's installed correctly.

## 2. Migration using goose

I recommend creating an sql directory in the root of your project, and in there creating a schema directory.

A "migration" is a SQL file that describes a change to your database schema. For now, we need our first migration to create a users table. The simplest format for these files is:

```sql
-- +goose Up
CREATE TABLE ...

-- +goose Down
DROP TABLE users;
```

The `-- +goose Up` and `-- +goose Down` comments are required. They tell Goose how to run the migration. An "up" migration moves your database from its old state to a new state. A "down" migration moves your database from its new state back to its old state.

By running all of the "up" migrations on a blank database, you should end up with a database in a ready-to-use state. "Down" migrations are only used when you need to roll back a migration, or if you need to reset a local testing database to a known state.

### Run the migration:

```zsh
goose postgres protocol://username:password@host:port/database up
```

## 3. Configure sqlc

Create a file called `sqlc.yaml` in the root of the project.

```yaml
version: "2"
sql:
  - schema: "sql/schema"
    queries: "sql/queries"
    engine: "postgresql"
    gen:
      go:
        out: "internal/database"
```

We're telling SQLC to look in the `sql/schema` directory for our schema structure (which is the same set of files that Goose uses, but sqlc automatically ignores "down" migrations), and in the `sql/queries` directory for queries. We're also telling it to generate Go code in the `internal/database` directory.
