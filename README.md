# CLI RSS feed aggregator in Go

This project uses PostgreSQL as its database and is written in Go.

## Requirements

To run it, you’ll need:

- PostgreSQL installed and running (with a database you can connect to)
  > https://www.postgresql.org/download/
- Go (version 1.25 or newer) installed and available on your PATH
  > https://go.dev/doc/install

## Installation

To install the _gator_ CLI, run:

```bash
go install github.com/juandrzej/postgreSQL-blog-aggregator/cmd/gator@latest
```

> This will build the binary and place gator in your $GOPATH/bin (or $GOBIN), which should be on your PATH.

## Configuration:

Create a config file at:

```bash
~/.gatorconfig.json`
```

with contents similar to:

```json
{
  "db_url": "postgres://user:password@localhost:5432/dbname?sslmode=disable"
}
```

**Make sure this URL matches your local Postgres setup.**

## Usage

After installing and configuring _gator_, you can run commands like:

```bash
gator register <username> # Create a new user
gator addfeed <url> # Follow an RSS feed
gator agg # Fetch new posts from followed feeds
gator browse # View recent posts
```

Run _gator help_ or _gator help <command>_ for more options.
