# JKM Email Client

JKM is a simple TUI email client.

- It tends to use `less`-like keybindings for navigation.
- It is configurable via environment variables and a `.env` file.
- It works with IMAP/SMTP servers.

## Configuration

Uses `.env` files if you care to, or environment variables, or else use the configuration flow in the app without those set.

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

- Right now this refreshes the inbox, versus using IDLE.
- I ought to build an RPM/Deb/apk/brew package.
- The router in the app has no "memory" of where the user was in the inbox if they read or compose an email.
- The test coverage here is weak! I depriorized testing in this case because getting the IMAP/SMTP/TUI test coverage seemed like it would be pretty fragile when it came to adding new features and changing the flow, in this kind of short-lived project.

- Deploy it somewhere

## Usage

- Navigate using arrow keys or hjkl.
- Press Enter to read a selected email.
- Press c to compose a new email (in the mailbox view.)
- Press Ctrl+C, or 'q' to quit from the mailbox view or return to the mailbox from the compose/read views.

## Data Flow

```
             +-------------------+
             |                   |
             |  cmd/jkm/main.go  |
             |    (entrypoint)   |
             |                   |
             +--------+----------+
                      |
                      | initializes...
                      |
     +----------------+------------------+
     |                                   |
     |            router                 |
     |     (navigation controller)       |
     |                                   |
     +----------------+------------------+
                      |
                      | loads config through...
                      |
     +----------------+------------------+
     |                                   |
     |            configure              |
     |     (environment/settings)        |
     |                                   |
     +----------------+------------------+
                      |
                      | which goes back to...
                      |
     +----------------+------------------+
     |                                   |
     |            router                 |
     |     (navigation controller)       |
     |                                   |
     +-----------------------------------+
         |         |            |      |
    +--------+  +--------+ +---------+ |
    |        |  |        | |         | |
    | list   |  | read   | | compose | | ticks...
    |        |  |        | |         | |
    +----+---+  +---+----+ +---------+ |
         |          |           |      |
     +-----------------------------------+
     |                                   |
     |       caching email client        |
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
