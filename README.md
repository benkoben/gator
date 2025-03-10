# gator

Local RSS aggregator that works in the command line. Using a database it saves state between usage.

# dependencies

```bash
# Install goose to handle database migrations
go install github.com/pressly/goose/v3/cmd/goose@latest

# Install sqlc to interact with the database in a type safe manner
go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest

# postgresql driver
go get github.com/lib/pq
```
