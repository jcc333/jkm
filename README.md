# JKM Email Client

JKM is a simple TUI email client.

pageSize?
refreshing emails automatically?
email logic cleanup :(

// TODO: when we exit another view and return to the listing model, we should
// restore the previous state of the model, including the current offset and
// page size.

## TODO
- Cleanup pass in the code.
- List view memory in the router (return to previous page of emails)
- Overview configurable size limit (10K default? 100K default?)
- Set up a fastmail account for jkm to show it off.
- Configuration view before the mailer gets initialized.
- Stand this up on my k8s cluster.
- Expose it to the open internet.
- Write some unit tests for it to get coverage.
- Add some logging to a file.
- Publish a RPM/Deb/apk/brew package.
- Add a GH build.
- GH Pages about it.

## Features

- List emails from inbox
- Read email messages
- Compose new emails using your system's default editor ($EDITOR)
- Send emails
- Configured via environment variables and .env file

## Configuration

JKM is configured using environment variables or a `.env` file in either your home directory or the current directory. The following variables are required:

- `JKM_SMTP_SERVER`: SMTP server for sending emails
- `JKM_SMTP_PORT`: SMTP port (defaults to 587)
- `JKM_IMAP_SERVER`: IMAP server for receiving emails
- `JKM_IMAP_PORT`: IMAP port (defaults to 993)
- `JKM_EMAIL`: Your email address
- `JKM_PASSWORD`: Your email password

Additionally, you need to have an editor set up:

- `EDITOR`: Your preferred text editor (e.g., vim, nano, emacs)

## Building

```bash
go build -o jkm ./cmd/jkm
```

## Running

```bash
./jkm
```

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
