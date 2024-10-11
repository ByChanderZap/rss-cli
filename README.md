# Overview
rss-cli is a CLI application for you to retrieve posts via RSS and have all of them on a single place. YOUR TERMINAL!

# Requirements
- You need a postgres database (there is where we are going to store our posts)
- Create a json config file on your home path, it should be: ~/.gatorconfig.json
- The config file needs a db_url (from the database where we will store our posts)

### Example json config file:
```json
{
   "db_url":"postgres://<dbuser>:<password>@localhost:5432/<db_name>?sslmode=disable",
}
```

# Commands list:
- Register a user:
```bash
gator register <name>
```
- Add a feed:
```bash
gator addfeed <url>
```
- Start aggregation:
```bash
gator agg <time> (example: 1m to execute every minute)
```
- Browse posts:
```bash
gator browse <limit>
```
- Login on a different user:
```bash
gator login <name>
```
- List feeds:
```bash
gator feeds
```
- List users:
```bash
gator users
```
- Follow feed:
```bash
gator follow <url>
```
- Unfopllow feed:
```bash
gator unfollow <url>
```


#### Notes:
Use goose for database migrations
