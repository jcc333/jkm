# JKM Email Client

JKM is a simple TUI email client.

- It tends to use `less`-like keybindings for navigation.
- It is configurable via environment variables and a `.env` file.
- It works with IMAP/SMTP servers.

## Configuration

```
JKM_EMAIL=flast@example.com
JKM_IMAP_SERVER=imap.fastmail.com
JKM_IMAP_PORT=993
JKM_IMAP_PASSWORD=somepassword
JKM_SMTP_SERVER=smtp.fastmail.com
JKM_SMTP_PORT=465
JKM_SMTP_PASSWORD=otherpassword
JKM_LOGGING=true #logging to jkm.logs.jsonl
```

## Still to be Done

- Right now this doesn't IDLE or refresh the fetched inbox.
- Currently we return to the inbox at the first page of results.
- I ought to build an RPM/Deb/apk/brew package.
- The router in the app has no "memory" of where the user was in the inbox if they read or compose an email.

// TODO: when we exit another view and return to the listing model, we should
// restore the previous state of the model, including the current offset and
// page size.

- Set up a fastmail account for jkm to show it off.
- Deploy it somewhere
- Cleanup pass in the code.
- List view memory in the router (keep track of the last UID scrolled over?)
- List view refreshes periodically or IDLEs
- Expose it to the open internet.
- Finish clean-up pass.
- Write some unit tests for it to get coverage.

## Usage

- Navigate using arrow keys or j/k
- Press Enter to read a selected email
- Press c to compose a new email
- Press r to refresh the inbox
- Press ? to toggle help
- Press Ctrl+C to quit


## Dependencies

- BubbleTea - Terminal UI framework
- go-imap - IMAP client library
- go-smtp - SMTP client library
- godotenv - .env file support

## Data Flow
```
             +-------------------+
             |                   |
             |  cmd/jkm/main.go  |
             |    (entrypoint)   |
             |                   |
             +--------+----------+
                      |
                      | loads config
                      v
     +----------------+------------------+
     |                                   |
     |     configure/config.go           |
     |     (environment/settings)        |
     |                                   |
     +----------------+------------------+
                      |
                      | initializes
                      v
     +----------------+------------------+
     |                                   |
     |     router/model.go               |
     |     (navigation controller)       |
     |                                   |
     +---+---------+--------+--------+---+
         |         |        |        |
         |         |        |        |
    +----v---+  +--v-----+  |   +----v----+
    |        |  |        |  |   |         |
    | list   |  | read   |  |   | compose |
    | model  |  | model  |  |   | model   |
    |        |  |        |  |   |         |
    +----+---+  +---+----+  |   +---------+
         ^          ^       |        |
         |          |       |        |
         |          |       |        |
         |          |       |        |
     +---v----------v-------v--------v---+
     |                                   |
     |          email/client.go          |
     |                                   |
     +---+-------------------------+-----+
         |                         |
         | IMAP                    | SMTP
         |                         |
  +------v--------+         +------v--------+
  |               |         |               |
  | External IMAP |         | External SMTP |
  |   Server      |         |    Server     |
  |               |         |               |
  +---------------+         +---------------+
```

1. main.go loads configuration from configure/config.go
2. main.go initializes the router from router/model.go
3. The router dispatches to different views:
- list/model.go for browsing emails
- read/model.go for reading a specific email
- compose/model.go for writing new emails
4. All views interact with email/client.go which:
- Connects to IMAP server for retrieving messages
- Connects to SMTP server for sending messages
5. Message passing between components is handled by messages/messages.go
6. The router processes events from sub-models and handles view transitions
