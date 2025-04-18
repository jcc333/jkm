package io

import (
	"testing"
	"os"
	"time"
	
	"github.com/jcc333/jkm/internal/configure"
	jkmemail "github.com/jcc333/jkm/internal/email"
)

func TestEmailClientWithTestServer(t *testing.T) {
	// This test will be skipped unless we explicitly run in integration test mode
	if testing.Short() {
		t.Skip("Skipping integration test")
	}
	
	// Check if we're running in test environment
	if os.Getenv("JKM_IMAP_SERVER") != "localhost" && os.Getenv("JKM_IMAP_SERVER") != "imap" {
		t.Skip("Skipping test: not running in Docker test environment")
	}
	
	// Load configuration
	config, err := configure.Load()
	if err != nil {
		t.Fatal(err)
	}
	if config.IMAPServer == "" {
		t.Fatal("IMAPServer is empty")
	}
	
	t.Run("Can connect to test server", func(t *testing.T) {
		// Create a new email client
		client, err := New(config)
		if err != nil {
			t.Fatal(err)
		}
		if client == nil {
			t.Fatal("client is nil")
		}
		
		// Clean up
		defer client.Disconnect()
		
		// Test count messages
		count, err := client.CountMessages()
		if err != nil {
			t.Fatal(err)
		}
		if count < 0 {
			t.Fatalf("Expected count >= 0, got %d", count)
		}
	})
	
	t.Run("Can send and receive email", func(t *testing.T) {
		// Create a new email client
		client, err := New(config)
		if err != nil {
			t.Fatal(err)
		}
		if client == nil {
			t.Fatal("client is nil")
		}
		
		// Clean up
		defer client.Disconnect()
		
		// Send a test email
		testSubject := "Test Email"
		testBody := "This is a test email from the integration test."
		err = client.Send(jkmemail.Message{
			MessageHeader: jkmemail.MessageHeader{
				Subject: testSubject,
				To:      []string{config.EmailAddress},
			},
			Body: testBody,
		})
		if err != nil {
			t.Fatal(err)
		}
		
		// Wait a moment for the email to be delivered
		time.Sleep(2 * time.Second)
		
		// List messages
		msgs, err := client.List(10, 0)
		if err != nil {
			t.Fatal(err)
		}
		
		// We should at least get one message
		if len(msgs) == 0 {
			t.Fatal("No messages found")
		}
		
		// Try to read the most recent message
		if len(msgs) > 0 {
			msg, err := client.Read(msgs[0].ID)
			if err != nil {
				t.Fatal(err)
			}
			if msg == nil {
				t.Fatal("message is nil")
			}
		}
	})
}