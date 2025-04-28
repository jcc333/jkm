package io

import (
	"crypto/tls"
	"fmt"
	"io"
	"net/smtp"
	"strings"

	"github.com/emersion/go-imap"
	imapClient "github.com/emersion/go-imap/client"
	"github.com/jordan-wright/email"

	"github.com/jcc333/jkm/internal/configure"
	jkmemail "github.com/jcc333/jkm/internal/email"
	"github.com/jcc333/jkm/internal/log"
)

// IMAP and SMTP functionality.

// An IMAP-based Sender
type Client struct {
	// Application-wide configuration.
	cfg *configure.Config

	// Underlying IMAP client.
	// Nil if not connected.
	in *imapClient.Client

	// SMTP configuration
	smtpServer   string
	smtpPort     int
	smtpEmail    string
	smtpPassword string
	smtpAuth     smtp.Auth

	// UIDs of messages in the inbox.
	uids []uint32

	// Cached read emails.
	messageInfos map[uint32]*imap.Message
}

// Disconnect from the IMAP server.
func (i *Client) Disconnect() error {
	if i.in == nil {
		return nil
	}
	i.in.Logout()
	return nil
}

// Get a new IMAP-based Sender.
func New(cfg *configure.Config) (*Client, error) {
	c := &Client{
		cfg:          cfg,
		messageInfos: make(map[uint32]*imap.Message),
	}
	err := c.Connect()
	if err != nil {
		log.Errorf("Failed to connect to IMAP server: %v", err)
		return nil, err
	}
	err = c.FetchMessages()
	if err != nil {
		log.Errorf("Failed to fetch messages: %v", err)
		return nil, err
	}
	return c, nil
}

// Connect to the IMAP server.
func (c *Client) Connect() error {
	var err error
	defer func() {
		if err != nil {
			log.Errorf("Failed to connect to IMAP server: %v", err)
			c.in = nil
		}
	}()

	if c.in != nil && c.smtpAuth != nil {
		return nil
	}

	imapAddr := fmt.Sprintf("%s:%d", c.cfg.IMAPServer, c.cfg.IMAPPort)
	c.in, err = imapClient.DialTLS(imapAddr, &tls.Config{InsecureSkipVerify: true})
	if err != nil {
		return err
	}

	if err := c.in.Login(c.cfg.EmailAddress, c.cfg.IMAPPassword); err != nil {
		return err
	}

	err = c.FetchMessages()
	if err != nil {
		log.Errorf("Failed to fetch messages: %v", err)
		return err
	}

	// Setup SMTP configuration
	c.smtpServer = c.cfg.SMTPServer
	c.smtpPort = c.cfg.SMTPPort
	c.smtpEmail = c.cfg.EmailAddress
	c.smtpPassword = c.cfg.SMTPPassword

	// Setup authentication
	c.smtpAuth = smtp.PlainAuth("", c.smtpEmail, c.smtpPassword, c.smtpServer)

	return nil
}

// Refresh the cached inbox contents
func (c *Client) FetchMessages() error {
	var err error
	defer func() {
		if err != nil {
			log.Errorf("Failed to fetch messages: %v", err)
		}
	}()
	if c.in == nil {
		err = fmt.Errorf("Fetching before connected to IMAP server")
		return err
	}

	mailbox, err := c.in.Select("INBOX", false)
	if err != nil {
		return err
	}

	// Get all message sequence numbers (not UIDs) in descending order (newest first)
	// This will use message sequence numbers, which are sorted by arrival time
	if mailbox.Messages == 0 {
		c.uids = []uint32{}
	} else {
		// Create sequence set for all messages - we'll later sort them in newest-first order
		seqSet := new(imap.SeqSet)
		// From 1 to highest message sequence number
		seqSet.AddRange(1, mailbox.Messages)

		// Load all UIDs
		ids, err := c.in.UidSearch(&imap.SearchCriteria{
			WithoutFlags: []string{imap.DeletedFlag},
		})
		if err != nil {
			return err
		}

		// Sort the UIDs in descending order
		if len(ids) > 1 {
			for i, j := 0, len(ids)-1; i < j; i, j = i+1, j-1 {
				ids[i], ids[j] = ids[j], ids[i]
			}
		}
		c.uids = ids
	}

	return nil
}

// Count the number of messages in the IMAP server.
func (c *Client) CountMessages() (int, error) {
	err := c.Connect()
	if err != nil {
		log.Errorf("Failed to fetch message count: %v", err)
		return 0, err
	}
	mailbox, err := c.in.Select("INBOX", false)
	if err != nil {
		return 0, err
	}
	return int(mailbox.Messages), nil
}

// Get all messages from the IMAP server.
// TODO: Implement pagination for large mailboxes.
// TODO: This does double-duty as a cache-buster/refresh function.
// Ideally we'd use IDLE instead, but this is adequate for now.
func (c *Client) List() ([]jkmemail.MessageHeader, error) {
	log.Info("io list messages")
	err := c.Connect()
	defer func() {
		if err != nil {
			log.Errorf("Failed to list messages: %v", err)
			c.in = nil
		}
	}()
	if err != nil {
		return nil, err
	}

	// Refresh the list of UIDs
	count, err := c.CountMessages()
	if err != nil {
		return nil, err
	}
	if count != len(c.messageInfos) {
		err = c.FetchMessages()
		if err != nil {
			return nil, err
		}
	}

	if len(c.uids) == 0 {
		return []jkmemail.MessageHeader{}, nil
	}

	seqSet := new(imap.SeqSet)
	for _, uid := range c.uids {
		seqSet.AddNum(uid)
	}

	items := []imap.FetchItem{imap.FetchEnvelope, imap.FetchUid, imap.FetchFlags}
	messages := make(chan *imap.Message, 100)
	done := make(chan error, 1)

	go func() {
		done <- c.in.UidFetch(seqSet, items, messages)
	}()

	// Clear the cache before populating with fresh data
	c.messageInfos = make(map[uint32]*imap.Message)

	// Store all messages in our cache
	for msg := range messages {
		c.messageInfos[msg.Uid] = msg
	}

	if err := <-done; err != nil {
		return nil, err
	}

	headers := make([]jkmemail.MessageHeader, 0, len(c.uids))
	for _, uid := range c.uids {
		msg, ok := c.messageInfos[uid]
		if !ok || msg == nil || msg.Envelope == nil {
			continue
		}

		var from string
		if len(msg.Envelope.From) > 0 {
			addr := msg.Envelope.From[0]
			from = fmt.Sprintf("%s <%s@%s>", addr.PersonalName, addr.MailboxName, addr.HostName)
		}

		to := make([]string, 0, len(msg.Envelope.To))
		for _, addr := range msg.Envelope.To {
			to = append(to, fmt.Sprintf("%s <%s@%s>", addr.PersonalName, addr.MailboxName, addr.HostName))
		}

		headers = append(headers, jkmemail.MessageHeader{
			ID:      int(uid),
			Subject: msg.Envelope.Subject,
			Date:    msg.Envelope.Date,
			From:    from,
			To:      to,
		})
	}

	return headers, nil
}

// Get a single message from the IMAP server.
func (c *Client) Read(id int) (*jkmemail.Message, error) {
	var err error
	defer func() {
		if err != nil {
			log.Errorf("Failed to read email with id %d: '%v'", id, err)
		}
	}()
	err = c.Connect()
	if err != nil {
		return nil, err
	}

	uid := uint32(id)

	seqSet := new(imap.SeqSet)
	seqSet.AddNum(uid)

	section := &imap.BodySectionName{
		BodyPartName: imap.BodyPartName{
			Specifier: imap.TextSpecifier,
		},
		Peek: true,
	}

	items := []imap.FetchItem{imap.FetchEnvelope, imap.FetchBody, imap.FetchBodyStructure, section.FetchItem()}
	messages := make(chan *imap.Message, 1)
	done := make(chan error, 1)
	go func() {
		done <- c.in.UidFetch(seqSet, items, messages)
	}()

	msg := <-messages
	if msg == nil {
		return nil, fmt.Errorf("message %d not found", id)
	}

	if err := <-done; err != nil {
		return nil, err
	}

	var body string
	r := msg.GetBody(section)
	if r == nil {
		return nil, fmt.Errorf("unable to get message body")
	}
	buf := new(strings.Builder)
	if _, err := io.Copy(buf, r); err != nil {
		return nil, fmt.Errorf("failed to read message body: %w", err)
	}
	body = buf.String()

	var from string
	if len(msg.Envelope.From) > 0 {
		addr := msg.Envelope.From[0]
		from = fmt.Sprintf("%s <%s@%s>", addr.PersonalName, addr.MailboxName, addr.HostName)
	}

	to := make([]string, 0, len(msg.Envelope.To))
	for _, addr := range msg.Envelope.To {
		to = append(to, fmt.Sprintf("%s <%s@%s>", addr.PersonalName, addr.MailboxName, addr.HostName))
	}

	return &jkmemail.Message{
		MessageHeader: jkmemail.MessageHeader{
			ID:      id,
			From:    from,
			To:      to,
			Subject: msg.Envelope.Subject,
			Date:    msg.Envelope.Date,
		},
		Body: body,
	}, nil
}

// Send an email via SMTP.
func (c *Client) Send(msg jkmemail.Message) error {
	err := c.Connect()
	defer func() {
		if err != nil {
			log.Warnf("Failed to send email with subject %s: '%v'", msg.Subject, err)
		}
	}()

	if err != nil {
		return fmt.Errorf("connection error: %w", err)
	}

	m := email.NewEmail()
	m.From = c.smtpEmail
	m.To = msg.To
	m.Subject = msg.Subject
	m.Text = []byte(msg.Body)

	addr := fmt.Sprintf("%s:%d", c.smtpServer, c.smtpPort)
	fmt.Printf("Sending via SMTP: server=%s, port=%d, from=%s\n", c.smtpServer, c.smtpPort, c.smtpEmail)

	// For Gmail, Yahoo, and many other providers, we need to use their SMTP server domains directly
	// The auth hostname must match the SMTP server name for many providers
	c.smtpAuth = smtp.PlainAuth("", c.smtpEmail, c.smtpPassword, c.smtpServer)

	// Fixed TLS config with proper hostname verification
	tlsConfig := &tls.Config{
		ServerName:         c.smtpServer,
		InsecureSkipVerify: false,
	}

	// Try TLS first (for port 465)
	if c.smtpPort == 465 {
		err = m.SendWithTLS(addr, c.smtpAuth, tlsConfig)
		if err == nil {
			return nil
		}
		fmt.Printf("TLS error: %v\n", err)
	}

	// Try StartTLS (for port 587)
	err = m.SendWithStartTLS(addr, c.smtpAuth, tlsConfig)
	if err == nil {
		return nil
	}
	fmt.Printf("StartTLS error: %v\n", err)

	// Fallback to unencrypted as last resort
	return m.Send(addr, c.smtpAuth)
}
