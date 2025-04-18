# JKM Testing

This directory contains testing utilities for JKM.

## Integration Tests

The integration tests use Docker Compose to spin up a test environment with:

- GreenMail IMAP/SMTP server
- Pre-loaded test emails

### Running Integration Tests

To run the integration tests, execute:

```sh
./scripts/run_integration_tests.sh
```

### Requirements

- Docker Compose
- Go 1.23+

### Test Environment

The test environment uses the following configuration:

- IMAP server on port 3143
- SMTP server on port 3025
- Test account: `test@example.com` with password `password`

### Adding Test Emails

To add test emails to be pre-loaded in the test environment, add `.eml` files to the `testdata/emails` directory. These will be mounted into the GreenMail container.

### Manual Testing

To manually start the test environment:

```sh
docker compose up -d
```

You can then connect to the IMAP/SMTP servers directly for testing.