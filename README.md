# gator

Local RSS aggregator that works in the command line.

# How it works

This tool supports multiple users/profiles where each user can add feeds to the database. Once a feed has been added other users can start following these feeds.
Content is automatically fetched by the program and the intended usage us to run `gator agg 10m` in a seperate window/background process so that feed content is automatically fetched and saved to the database.

Users can interact with the database using various command line arguments

## Installation

```bash
go install https://github.com/benkoben/gator
```

Once installed the program looks for a `.gatorconfig.json` file in the user's home directory. The file should contain the postgres connection string initially. The program will add and modify the configuration during runtime.

```
{
  "db_url": "postgres://kooijman:@localhost:5432/gator?sslmode=disable",
}
```

## Command line options

`./gator <option> [arg]`

OPTIONS:

* `users` - List all users/profiles registered. this is also used to see which current used is logged in.
* `login <NAME>` - Login as an existing user
* `register <NAME>` - register a new user/profile
* `feeds` - Lists all feeds that have been added to the database
* `addfeed <URL>` - Add a new feed to the database
* `follow <URL>` - Follow an existing feed
* `following` - List which feeds you are following
* `unfollow <URL>` - Unfollow a feed
* `browse <LIMIT>` - List the feeds content

# Dependencies

* Go version 1.23
* PostgresQL 15.12 (developed with)

The following third party packages have to be installed:

```bash
# Install goose to handle database migrations
go install github.com/pressly/goose/v3/cmd/goose@latest

# Install sqlc to interact with the database in a type safe manner
go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest

# postgresql driver
go get github.com/lib/pq
```

# Developer references

## Goose

```bash
goose postgres "postgres://kooijman:@localhost:5432/gator" up
```

```bash
goose postgres "postgres://kooijman:@localhost:5432/gator" down
```

# Future improvements

- Proper error handling
- middleware for verifying if an URL is valid
- enable SSL mode for postgres connections
- Add sorting and filtering options to the browse command
- Add pagination to the browse command
- Add concurrency to the agg command so that it can fetch more frequently
- Add a search command that allows for fuzzy searching of posts
- Add bookmarking or liking posts
- Add a TUI that allows you to select a post in the terminal and view it in a more readable format (either in the terminal or open in a browser)
- Add an HTTP API (and authentication/authorization) that allows other users to interact with the service remotely
- Write a service manager that keeps the agg command running in the background and restarts it if it crashes


# Additional docs

* [pq driver error codes](https://github.com/lib/pq/blob/master/error.go)
