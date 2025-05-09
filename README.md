# Gator CLI

Gator is a command-line tool for managing RSS feeds and users, built with Go and PostgreSQL.

## Prerequisites

- **Go** (1.18+)
- **PostgreSQL** (running and accessible)

## Installation

Install the Gator CLI with:

```sh
go install github.com/lucashthiele/gator@latest
```

This will place the `gator` binary in your `$GOPATH/bin` (make sure this is in your `PATH`).

## Configuration

Create a `.gatorconfig.json` file in your home directory:

```json
{
  "db_url": "postgres://username:password@localhost:5432/yourdb?sslmode=disable",
  "current_user_name": ""
}
```

Replace `username`, `password`, and `yourdb` with your PostgreSQL credentials and database name.

## Database Setup

This project uses [goose](https://github.com/pressly/goose) for database migrations. Install goose if you haven't already:

```sh
go install github.com/pressly/goose/v3/cmd/goose@latest
```

Run the migrations in the `sql/schema/` directory:

```sh
goose -dir sql/schema postgres "postgres://username:password@localhost:5432/yourdb?sslmode=disable" up
```

Replace `username`, `password`, and `yourdb` with your PostgreSQL credentials and database name.

## Usage

After configuring your database and installing the CLI, you can run commands like:

- `gator register <username>` — Register a new user
- `gator login <username>` — Set the current user
- `gator addfeed <name> <url>` — Add a new RSS feed and follow it
- `gator feeds` — List all available feeds
- `gator follow <feed-url>` — Follow a feed by its URL
- `gator unfollow <feed-url>` — Unfollow a feed
- `gator following` — List feeds you are following
- `gator browse [limit]` — Browse recent posts from followed feeds
- `gator agg <duration>` — Start the feed aggregator (e.g., `1m` for every minute)
- `gator users` — List all users
- `gator reset` — Delete all users (for development/testing)

## Notes

- The CLI stores the current user in the config file.
- Make sure your PostgreSQL server is running and accessible via the `db_url` in your config.
- For development, you may want to reset the database using the `reset` command.
