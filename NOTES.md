---
 Connection string: postgres://postgres:postgres@localhost:5432/gator

Start postgres server
- sudo service postgresql start

Enter postgres sql shell with
- Linux: `sudo -u postgres psql

Connect Database
- \c {database}

| command | description |
| --- | --- |
| `sudo service postgresql start` | start postgres server in the background |
| `sudo -u postgres psql` |enter postgres shell|
| \c {database} | connect to specific database|

- goose cmd
goose -dir ./sql/schema postgres "postgres://postgres:postgres@localhost:5432/gator" up


UPDATE feeds
SET updated_at = '2026-03-9 22:18:01.737593', last_fetched_at = '2026-03-9 22:18:01.737593'
WHERE url='fb.com';

