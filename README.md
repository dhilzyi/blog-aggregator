## Description

This is a blog aggregator CLI project built by following a course from Boot.dev. The application fetches RSS feeds and stores posts in a local PostgreSQL database.

Course link:  
[Build a Blog Aggregator – Boot.dev](https://www.boot.dev/lessons/14b7179b-ced3-4141-9fa5-e67dbc3e5242)

---

## Requirements

Make sure you have the following installed:

- PostgreSQL 16+
- Go 1.22+

---

## Installation

Install the CLI tool using:

```bash
go install github.com/dhilzyi/blog-aggregator@latest
```

Then run:

```bash
blog-aggregator
```

---

## Configuration

### 1. Create the database

```bash
createdb gator
```

---

### 2. Create config file

Create the file:

```bash
~/.gatorconfig.json
```

Add the following:

```json
{
  "db_url": "postgres://username:password@localhost:5432/gator?sslmode=disable",
  "current_user_name": "admin"
}
```

Update the database URL with your own credentials.

---

## Usage

Commands are passed as arguments to the CLI:

```bash
blog-aggregator <command> [arguments...]
```

---

## Commands

| Command   | Description                                      | Arguments                     | Example |
|----------|--------------------------------------------------|------------------------------|--------|
| login    | Switch to an existing user                       | `<name>`                     | `blog-aggregator login ruru` |
| register | Register a new user                              | `<name>`                     | `blog-aggregator register ruru` |
| reset    | Delete all data in all tables ⚠️                 | none                         | `blog-aggregator reset` |
| addfeed  | Add a new RSS feed                               | `<name> <url>`               | `blog-aggregator addfeed "Hacker News" "https://news.ycombinator.com/rss"` |
| feeds    | List all feeds for the current user              | none                         | `blog-aggregator feeds` |
| browse   | Show posts from followed feeds                   | none                         | `blog-aggregator browse` |
| agg      | Start feed aggregation loop                      | `<interval>` (e.g. `1m`)     | `blog-aggregator agg 1m` |
