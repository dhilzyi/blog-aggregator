## Description
It's a blog aggregator project by following guide course from boot.dev

## Requirements
You need to install some program in your computer to able run this program.
- Postgres 16.13+
- Go 1.22.2+

## Installation
After you have golang in your system. Run this command in terminal to install.  
`go install github.com/dhilzyi/blog-aggregator@latest`
  
Type `blog-aggregator` to run it.

## Configuration
You might want to create the config file for yourself  
[1]. First create the file following the path `~/.gatorconfig.json`.  
[2]. Copy and paste below code into the file  
```
{
 "db_url": "postgres://postgres:postgres@localhost:5432/gator?sslmode=disable",
 "current_user_name": "admin"
}
```  
[3]. Save the file

You can change the db url connection according to your own user and password.

## Commands

You need pass command name as a first argument then after that the other arguments will be passed for the command itself.

| Command | Explantation | Arguments |Example |
| --- | --- | --- | --- |
| login | Change account to specific user |( name )| blog-aggregator login ruru |
| register | Register new account to table users |( name )| .. register ruru |
| reset | WARNING: this will delete all rows in all tables | | .. reset|
| addfeed |Add new feed to the table feeds| ( name ) ( link )|.. "Hacker News" "https://news.ycombinator.com/rss"|
| feeds |Show all feeds for current user| |.. feeds|
| browse |Show all posts for current user| |.. browse|
| agg |Will fetch every interval time| ( interval time ) |.. 1m|
